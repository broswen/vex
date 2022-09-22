package api

import (
	"github.com/broswen/vex/internal/token"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

func AccountAuthorizer(next http.Handler, tokenStore token.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		log.Debug().Str("token", token).Msg("found authorization token")
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
}
