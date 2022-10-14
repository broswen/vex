package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/broswen/vex/internal/token"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"time"
)

var (
	AccessCookieName   = "CF_Authorization"
	AccessIdentityPath = "/cdn-cgi/access/get-identity"
)

type AccessIdentity struct {
	Email     string `json:"email"`
	UserUUID  string `json:"user_uuid"`
	AccountID string `json:"account_id"`
}

func AccountAuthorizer(tokenStore token.Store) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			accountId, err := accountId(r)
			if err != nil {
				writeErr(w, nil, err)
				return
			}
			if authHeader == "" {
				log.Warn().Msg("didn't find authorization header")
				writeErr(w, nil, ErrUnauthorized)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 {
				log.Warn().Msg("didn't find bearer token")
				writeErr(w, nil, ErrUnauthorized)
				return
			}
			token := parts[1]
			t, err := tokenStore.GetByHash(r.Context(), token)
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

func CloudflareAccessIdentityLogger() func(next http.Handler) http.Handler {
	client := http.Client{
		Timeout: time.Second * 3,
	}
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			cfAuthorization, err := r.Cookie(AccessCookieName)
			if err != nil {
				log.Warn().Err(err).Msg("no CF_Authorization cookie found")
				next.ServeHTTP(w, r)
				return
			}
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://%s%s", r.Host, AccessIdentityPath), nil)
			if err != nil {
				log.Error().Err(err).Msg("creating access identity request")
				next.ServeHTTP(w, r)
				return
			}
			req.AddCookie(cfAuthorization)
			res, err := client.Do(req)
			if err != nil {
				log.Error().Err(err).Msg("sending access identity request")
				next.ServeHTTP(w, r)
				return
			}
			if res.StatusCode >= http.StatusBadRequest {
				log.Error().Str("status", res.Status).Int("code", res.StatusCode).Msg("received access identity request")
				next.ServeHTTP(w, r)
				return
			}
			identity := &AccessIdentity{}
			err = json.NewDecoder(res.Body).Decode(identity)
			if err != nil {
				log.Error().Err(err).Msg("decoding access identity")
				next.ServeHTTP(w, r)
				return
			}
			log.Debug().Str("email", identity.Email).Str("user_uuid", identity.UserUUID).Str("account_id", identity.AccountID).Str("method", r.Method).Str("path", r.URL.Path).Msg("access identity")
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
