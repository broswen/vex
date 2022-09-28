package api

import (
	"context"
	"fmt"
	"github.com/broswen/vex/internal/token"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

func AccountAuthorizer(tokenStore token.TokenStore) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			accountId, err := accountId(r)
			if err != nil {
				writeErr(w, nil, err)
				return
			}
			if authHeader == "" {
				writeErr(w, nil, ErrUnauthorized)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 {
				writeErr(w, nil, ErrUnauthorized)
				return
			}
			token := parts[1]
			log.Debug().Msg("found authorization token")
			t, err := tokenStore.GetByToken(r.Context(), token)
			if err != nil {
				log.Debug().Err(err).Msg("couldn't get by token")
				writeErr(w, nil, ErrUnauthorized)
				return
			}
			//token must match account in path
			if t.AccountID != accountId {
				writeErr(w, nil, ErrUnauthorized)
				return
			}
			//readonly tokens can only GET
			if t.ReadOnly && r.Method != http.MethodGet {
				writeErr(w, nil, ErrUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func CloudflareAccessVerifier(teamDomain, policyAUD string) func(next http.Handler) http.Handler {

	if teamDomain == "" || policyAUD == "" {
		log.Warn().Str("teamDomain", teamDomain).Str("policyAUD", policyAUD).Msg("Cloudflare Access verification disabled")
		//Skip JWT verification
		return func(next http.Handler) http.Handler {
			fn := func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			}
			return http.HandlerFunc(fn)
		}
	} else {
		log.Debug().Str("teamDomain", teamDomain).Str("policyAUD", policyAUD).Msg("Cloudflare Access verification enabled")
	}
	var certsURL = fmt.Sprintf("%s/cdn-cgi/access/certs", teamDomain)

	var config = &oidc.Config{
		ClientID: policyAUD,
	}
	var keySet = oidc.NewRemoteKeySet(context.Background(), certsURL)
	var verifier = oidc.NewVerifier(teamDomain, keySet, config)

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			headers := r.Header

			// Make sure that the incoming request has our token header
			//  Could also look in the cookies for CF_AUTHORIZATION
			accessJWT := headers.Get("Cf-Access-Jwt-Assertion")
			if accessJWT == "" {
				log.Debug().Msg("couldn't get authorization token")
				log.Debug().Str("accessJWT", accessJWT).Msg("")
				writeErr(w, nil, ErrUnauthorized)
				return
			}

			// Verify the access token
			ctx := r.Context()
			_, err := verifier.Verify(ctx, accessJWT)
			if err != nil {
				log.Debug().Err(err).Msg("invalid token")
				writeErr(w, nil, ErrUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
