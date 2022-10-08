package api

import (
	"github.com/broswen/vex/internal/stats"
	"github.com/broswen/vex/internal/token"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (api *API) ListTokens() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId, err := accountId(r)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		defer r.Body.Close()
		p := pagination(r)
		tokens, err := api.Token.List(r.Context(), accountId, p.Limit, p.Offset)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		err = writeOK(w, http.StatusOK, tokens)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
	}
}

func (api *API) GenerateToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId, err := accountId(r)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		readOnly := r.URL.Query().Get("readOnly")
		defer r.Body.Close()
		t, err := api.Token.Generate(r.Context(), accountId, readOnly == "true")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		stats.TokenCreated.Inc()

		err = api.Provisioner.ProvisionToken(r.Context(), t)
		if err != nil {
			log.Warn().Str("id", t.ID).Err(err).Msg("could not provision token")
		}
		err = writeOK(w, http.StatusOK, t)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
	}
}

func (api *API) RerollToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//accountId := chi.URLParam(r, "accountId")
		tokenId := chi.URLParam(r, "tokenId")
		defer r.Body.Close()
		t, err := api.Token.Get(r.Context(), tokenId)
		if err != nil {
			writeErr(w, nil, err)
			return
		}

		updatedToken, err := api.Token.Reroll(r.Context(), t.ID)
		if err != nil {
			log.Error().Err(err).Msg("")
			writeErr(w, nil, err)
			return
		}
		stats.TokenRolled.Inc()

		err = api.Provisioner.ProvisionToken(r.Context(), updatedToken)
		if err != nil {
			log.Warn().Str("id", updatedToken.ID).Err(err).Msg("could not provision new token")
		}

		//deprovision old token value
		err = api.Provisioner.DeprovisionToken(r.Context(), t)
		if err != nil {
			log.Warn().Str("id", t.ID).Err(err).Msg("could not deprovision old token")
		}

		err = writeOK(w, http.StatusOK, updatedToken)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
	}
}

func (api *API) DeleteToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenId := chi.URLParam(r, "tokenId")
		t, err := api.Token.Get(r.Context(), tokenId)
		if err != nil {
			writeErr(w, nil, err)
			return
		}

		err = api.Token.Delete(r.Context(), tokenId)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		stats.TokenDeleted.Inc()

		err = api.Provisioner.DeprovisionToken(r.Context(), &token.Token{Token: t.Token})
		if err != nil {
			log.Warn().Str("id", t.ID).Err(err).Msg("could not deprovision token")
		}

		err = writeOK(w, http.StatusOK, t)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
	}
}
