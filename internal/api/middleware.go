package api

import (
	"github.com/broswen/vex/internal/token"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

func AccountAuthorizer(next http.Handler, tokenStore token.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		accountId := chi.URLParam(r, "accountId")
		if authHeader == "" {
			RenderError(w, ErrUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			RenderError(w, ErrUnauthorized)
			return
		}
		token := parts[1]
		log.Debug().Str("token", token).Msg("found authorization token")
		t, err := tokenStore.GetByToken(r.Context(), token)
		if err != nil {
			log.Debug().Err(err).Msg("couldn't get by token")
			RenderError(w, ErrUnauthorized)
			return
		}
		//token must match account in path
		if t.AccountID != accountId {
			RenderError(w, ErrUnauthorized)
			return
		}
		//readonly tokens can only GET
		if t.ReadOnly && r.Method != http.MethodGet {
			RenderError(w, ErrUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
