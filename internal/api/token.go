package api

import (
	"encoding/json"
	"github.com/broswen/vex/internal/db"
	"github.com/broswen/vex/internal/token"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GenerateToken(tokenStore token.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId := chi.URLParam(r, "accountId")
		readOnly := r.URL.Query().Get("readOnly")
		defer r.Body.Close()
		t, err := tokenStore.Generate(r.Context(), accountId, readOnly == "true")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func RerollToken(tokenStore token.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId := chi.URLParam(r, "accountId")
		tokenId := chi.URLParam(r, "tokenId")
		t := &token.Token{}
		t.AccountID = accountId
		t.ID = tokenId
		defer r.Body.Close()

		err := tokenStore.Reroll(r.Context(), t)
		if err != nil {
			if err == db.ErrNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		err = json.NewEncoder(w).Encode(t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func DeleteToken(tokenStore token.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenId := chi.URLParam(r, "tokenId")
		err := tokenStore.Delete(r.Context(), tokenId)
		if err != nil {
			if err == db.ErrNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		err = json.NewEncoder(w).Encode(&struct{ id string }{id: tokenId})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
