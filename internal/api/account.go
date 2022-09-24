package api

import (
	"github.com/broswen/vex/internal/account"
	"net/http"
)

func CreateAccount(accountStore account.AccountStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &account.Account{}
		err := readJSON(w, r, p)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		defer r.Body.Close()
		err = accountStore.Insert(r.Context(), p)

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

func UpdateAccount(accountStore account.AccountStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId, err := accountId(r)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		p := &account.Account{}
		err = readJSON(w, r, p)
		if err != nil {
			writeErr(w, nil, ErrBadRequest.WithError(err))
			return
		}
		defer r.Body.Close()

		p.ID = accountId

		err = accountStore.Update(r.Context(), p)

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

func ListAccounts(accountStore account.AccountStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := accountStore.List(r.Context())
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

func GetAccount(accountStore account.AccountStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId, err := accountId(r)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		p, err := accountStore.Get(r.Context(), accountId)
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

func DeleteAccount(accountStore account.AccountStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId, err := accountId(r)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		err = accountStore.Delete(r.Context(), accountId)
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
