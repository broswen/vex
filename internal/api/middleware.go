package api

import (
	"github.com/broswen/vex/internal/token"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strings"
)

func AccountAuthorizer(next http.Handler, tokenStore token.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		accountId := chi.URLParam(r, "accountId")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := parts[1]
		log.Printf("found authorization token: %s", token)
		t, err := tokenStore.Get(r.Context(), token)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		//token must match account in path
		if t.AccountID != accountId {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		//readonly tokens can only GET
		if t.ReadOnly && r.Method != http.MethodGet {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
