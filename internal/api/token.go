package api

import (
	"encoding/json"
	"github.com/broswen/vex/internal/db"
	"github.com/broswen/vex/internal/provisioner"
	"github.com/broswen/vex/internal/token"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func ListTokens(tokenStore token.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId := chi.URLParam(r, "accountId")
		defer r.Body.Close()
		tokens, err := tokenStore.List(r.Context(), accountId, 100, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(tokens)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GenerateToken(tokenStore token.TokenStore, provisioner provisioner.Provisioner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId := chi.URLParam(r, "accountId")
		readOnly := r.URL.Query().Get("readOnly")
		defer r.Body.Close()
		t, err := tokenStore.Generate(r.Context(), accountId, readOnly == "true")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = provisioner.ProvisionToken(r.Context(), t)
		if err != nil {
			log.Printf("provision %s: %v", t.ID, err)
		}
		err = json.NewEncoder(w).Encode(t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func RerollToken(tokenStore token.TokenStore, provisioner provisioner.Provisioner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//accountId := chi.URLParam(r, "accountId")
		tokenId := chi.URLParam(r, "tokenId")
		defer r.Body.Close()
		t, err := tokenStore.Get(r.Context(), tokenId)
		if err != nil {
			if err == db.ErrNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		oldToken := t.Token

		err = tokenStore.Reroll(r.Context(), t)
		if err != nil {
			if err == db.ErrNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		err = provisioner.ProvisionToken(r.Context(), t)
		if err != nil {
			log.Printf("provision %s: %v", t.ID, err)
		}

		//deprovision old token value
		err = provisioner.DeprovisionToken(r.Context(), &token.Token{Token: oldToken})
		if err != nil {
			log.Printf("deprovision %s: %v", tokenId, err)
		}

		err = json.NewEncoder(w).Encode(t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func DeleteToken(tokenStore token.TokenStore, provisioner provisioner.Provisioner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenId := chi.URLParam(r, "tokenId")
		t, err := tokenStore.Get(r.Context(), tokenId)
		if err != nil {
			if err == db.ErrNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		err = tokenStore.Delete(r.Context(), tokenId)
		if err != nil {
			if err == db.ErrNotFound {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		err = provisioner.DeprovisionToken(r.Context(), &token.Token{Token: t.Token})
		if err != nil {
			log.Printf("deprovision %s: %v", tokenId, err)
		}
		err = json.NewEncoder(w).Encode(&struct{ id string }{id: tokenId})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
