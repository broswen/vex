package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/broswen/vex/internal/account"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateAccountHandler(t *testing.T) {
	a1 := &account.Account{
		Name:        "test",
		Description: "test account",
	}
	reqBody, err := json.Marshal(a1)
	assert.Nil(t, err)
	req, err := http.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(reqBody))
	assert.Nil(t, err)
	req.WithContext(context.Background())
	rr := httptest.NewRecorder()
	store := account.NewMockStore()
	store.On("Insert").Return(nil)
	app := &API{
		Account: store,
	}
	r := chi.NewRouter()
	r.Post("/accounts", app.CreateAccount())
	r.ServeHTTP(rr, req)
	assert.Equalf(t, rr.Code, http.StatusOK, "should create new account")
	store.AssertExpectations(t)
}

func TestGetAccountHandler(t *testing.T) {
	a1 := &account.Account{
		Name:        "test",
		Description: "test account",
	}
	store := account.NewMockStore()
	//create test account
	store.On("Insert").Return(nil)
	err := store.Insert(context.Background(), a1)
	assert.Nil(t, err)
	assert.Equal(t, "0", a1.ID)

	//get test account through chi router and handler
	rr := httptest.NewRecorder()

	store.On("Get").Return(nil)

	app := &API{
		Account: store,
	}
	r := chi.NewRouter()
	r.Get("/accounts/{accountId}", app.GetAccount())

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%s", a1.ID), nil)
	assert.Nil(t, err)
	r.ServeHTTP(rr, req)
	assert.Equalf(t, http.StatusOK, rr.Code, "should get existing account")

	//non existant account should return 404
	rr = httptest.NewRecorder()
	store.On("Get").Return(nil)

	req, err = http.NewRequest(http.MethodGet, "/accounts/1", nil)
	assert.Nil(t, err)
	r.ServeHTTP(rr, req)
	assert.Equalf(t, http.StatusNotFound, rr.Code, "should not find nonexistant account")
	store.AssertExpectations(t)
}
