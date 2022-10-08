package api

import (
	"github.com/broswen/vex/internal/account"
	"net/http"
)

func (api *API) CreateAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a := &account.Account{}
		err := readJSON(w, r, a)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		defer r.Body.Close()
		newAccount, err := api.Account.Insert(r.Context(), a)

		if err != nil {
			writeErr(w, nil, err)
			return
		}

		err = writeOK(w, http.StatusOK, newAccount)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
	}
}

func (api *API) UpdateAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId, err := accountId(r)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		a := &account.Account{}
		err = readJSON(w, r, a)
		if err != nil {
			writeErr(w, nil, ErrBadRequest.WithError(err))
			return
		}
		defer r.Body.Close()

		a.ID = accountId

		updatedAccount, err := api.Account.Update(r.Context(), a)

		if err != nil {
			writeErr(w, nil, err)
			return
		}

		err = writeOK(w, http.StatusOK, updatedAccount)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
	}
}

func (api *API) ListAccounts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := pagination(r)
		a, err := api.Account.List(r.Context(), p.Limit, p.Offset)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		err = writeOK(w, http.StatusOK, a)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
	}
}

func (api *API) GetAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId, err := accountId(r)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		p, err := api.Account.Get(r.Context(), accountId)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		err = writeOK(w, http.StatusOK, p)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
	}
}

func (api *API) DeleteAccount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId, err := accountId(r)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		err = api.Account.Delete(r.Context(), accountId)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		err = writeOK(w, http.StatusOK, &struct{ id string }{id: accountId})
		if err != nil {
			writeErr(w, nil, err)
			return
		}
	}
}
