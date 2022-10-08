package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/broswen/vex/internal/account"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var now = time.Now()
var accountID = "da863681-2f59-432d-848d-a64fbfbeab34"

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
	store.On("Insert", mock.Anything, a1).Return(&account.Account{
		ID:          accountID,
		Name:        "test",
		Description: "test account",
		CreatedOn:   now,
		ModifiedOn:  now,
	}, nil)
	app := &API{
		Account: store,
	}
	r := chi.NewRouter()
	r.Post("/accounts", app.CreateAccount())
	r.ServeHTTP(rr, req)
	assert.Equalf(t, http.StatusOK, rr.Code, "should return ok")
	store.AssertExpectations(t)
}

func TestUpdateAccountHandler(t *testing.T) {
	a1 := &account.Account{
		Name:        "test",
		Description: "test account",
	}
	reqBody, err := json.Marshal(a1)
	assert.Nil(t, err)
	req, err := http.NewRequest(http.MethodPut, "/accounts/"+accountID, bytes.NewReader(reqBody))
	assert.Nil(t, err)
	req.WithContext(context.Background())
	rr := httptest.NewRecorder()
	store := account.NewMockStore()
	store.On("Update", mock.Anything, &account.Account{
		ID:          accountID,
		Name:        "test",
		Description: "test account",
	}).Return(&account.Account{
		ID:          accountID,
		Name:        "test",
		Description: "test account",
		CreatedOn:   now,
		ModifiedOn:  now,
	}, nil)
	app := &API{
		Account: store,
	}
	r := chi.NewRouter()
	r.Put("/accounts/{accountId}", app.UpdateAccount())
	r.ServeHTTP(rr, req)
	assert.Equalf(t, http.StatusOK, rr.Code, "should return ok")
	store.AssertExpectations(t)
}

func TestListAccountsHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/accounts", nil)
	assert.Nil(t, err)
	req.WithContext(context.Background())
	rr := httptest.NewRecorder()
	store := account.NewMockStore()
	store.On("List", mock.Anything, int64(100), int64(0)).Return([]*account.Account{
		{
			ID:          accountID,
			Name:        "test",
			Description: "test account",
			CreatedOn:   now,
			ModifiedOn:  now,
		},
	}, nil)
	app := &API{
		Account: store,
	}
	r := chi.NewRouter()
	r.Get("/accounts", app.ListAccounts())
	r.ServeHTTP(rr, req)
	assert.Equalf(t, rr.Code, http.StatusOK, "should return ok")
	store.AssertExpectations(t)
}

func TestGetAccountHandler_GetAccount(t *testing.T) {
	store := account.NewMockStore()
	//get test account through chi router and handler
	rr := httptest.NewRecorder()

	store.On("Get", mock.Anything, accountID).Return(&account.Account{
		ID:          accountID,
		Name:        "test",
		Description: "test account",
		CreatedOn:   now,
		ModifiedOn:  now,
	}, nil)

	app := &API{
		Account: store,
	}
	r := chi.NewRouter()
	r.Get("/accounts/{accountId}", app.GetAccount())

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%s", accountID), nil)
	assert.Nil(t, err)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	store.AssertExpectations(t)
}

func TestGetAccountHandler_BadId(t *testing.T) {
	store := account.NewMockStore()
	//get test account through chi router and handler
	rr := httptest.NewRecorder()

	app := &API{
		Account: store,
	}
	r := chi.NewRouter()
	r.Get("/accounts/{accountId}", app.GetAccount())

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%s", "123"), nil)
	assert.Nil(t, err)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetAccountHandler_NotFound(t *testing.T) {
	store := account.NewMockStore()
	//get test account through chi router and handler
	rr := httptest.NewRecorder()

	store.On("Get", mock.Anything, accountID).Return(&account.Account{}, account.ErrAccountNotFound{})

	app := &API{
		Account: store,
	}
	r := chi.NewRouter()
	r.Get("/accounts/{accountId}", app.GetAccount())

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%s", accountID), nil)
	assert.Nil(t, err)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
	store.AssertExpectations(t)
}
