package api

import (
	"context"
	provisioner2 "github.com/broswen/vex/internal/provisioner"
	"github.com/broswen/vex/internal/token"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var tokenID = "da863681-2f59-432d-848d-a64fbfbeab51"

func TestGenerateTokenHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/accounts/"+accountID+"/tokens", nil)
	assert.Nil(t, err)
	req.WithContext(context.Background())
	rr := httptest.NewRecorder()
	store := token.NewMockStore()
	store.On("Generate", mock.Anything, accountID, false).Return(&token.Token{
		ID:         tokenID,
		AccountID:  accountID,
		Token:      "abc123",
		TokenHash:  []byte("abc123"),
		ReadOnly:   false,
		CreatedOn:  now,
		ModifiedOn: now,
	}, nil)

	provisioner := provisioner2.NewMockProvisioner()
	provisioner.On("ProvisionToken", mock.Anything, &token.Token{
		ID:         tokenID,
		AccountID:  accountID,
		Token:      "abc123",
		TokenHash:  []byte("abc123"),
		ReadOnly:   false,
		CreatedOn:  now,
		ModifiedOn: now,
	}).Return(nil)

	app := &API{
		Token:       store,
		Provisioner: provisioner,
	}
	r := chi.NewRouter()
	r.Post("/accounts/{accountId}/tokens", app.GenerateToken())
	r.ServeHTTP(rr, req)
	assert.Equalf(t, http.StatusOK, rr.Code, "should return ok")
	store.AssertExpectations(t)
}

func TestRerollTokenHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "/accounts/"+accountID+"/tokens/"+tokenID, nil)
	assert.Nil(t, err)
	req.WithContext(context.Background())
	rr := httptest.NewRecorder()
	store := token.NewMockStore()

	store.On("Get", mock.Anything, tokenID).Return(&token.Token{
		ID:         tokenID,
		AccountID:  accountID,
		Token:      "abc123",
		TokenHash:  []byte("abc123"),
		ReadOnly:   false,
		CreatedOn:  now,
		ModifiedOn: now,
	}, nil)
	store.On("Reroll", mock.Anything, tokenID).Return(&token.Token{
		ID:         tokenID,
		AccountID:  accountID,
		Token:      "abc456",
		TokenHash:  []byte("abc456"),
		ReadOnly:   false,
		CreatedOn:  now,
		ModifiedOn: now,
	}, nil)

	provisioner := provisioner2.NewMockProvisioner()
	provisioner.On("ProvisionToken", mock.Anything, &token.Token{
		ID:         tokenID,
		AccountID:  accountID,
		Token:      "abc456",
		TokenHash:  []byte("abc456"),
		ReadOnly:   false,
		CreatedOn:  now,
		ModifiedOn: now,
	}).Return(nil)

	provisioner.On("DeprovisionToken", mock.Anything, &token.Token{
		ID:         tokenID,
		AccountID:  accountID,
		Token:      "abc123",
		TokenHash:  []byte("abc123"),
		ReadOnly:   false,
		CreatedOn:  now,
		ModifiedOn: now,
	}).Return(nil)

	app := &API{
		Token:       store,
		Provisioner: provisioner,
	}
	r := chi.NewRouter()
	r.Put("/accounts/{accountId}/tokens/{tokenId}", app.RerollToken())
	r.ServeHTTP(rr, req)
	assert.Equalf(t, http.StatusOK, rr.Code, "should return ok")
	store.AssertExpectations(t)
}

func TestListTokensHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/accounts/"+accountID+"/tokens", nil)
	assert.Nil(t, err)
	req.WithContext(context.Background())
	rr := httptest.NewRecorder()
	store := token.NewMockStore()
	store.On("List", mock.Anything, accountID, int64(100), int64(0)).Return([]*token.Token{
		{
			ID:         tokenID,
			AccountID:  accountID,
			Token:      "abc123",
			TokenHash:  []byte("abc123"),
			ReadOnly:   false,
			CreatedOn:  now,
			ModifiedOn: now,
		},
	}, nil)

	app := &API{
		Token: store,
	}
	r := chi.NewRouter()
	r.Get("/accounts/{accountId}/tokens", app.ListTokens())
	r.ServeHTTP(rr, req)
	assert.Equalf(t, http.StatusOK, rr.Code, "should return ok")
	store.AssertExpectations(t)
}

func TestDeleteTokenHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "/accounts/"+accountID+"/tokens/"+tokenID, nil)
	assert.Nil(t, err)
	req.WithContext(context.Background())
	rr := httptest.NewRecorder()
	store := token.NewMockStore()
	store.On("Delete", mock.Anything, tokenID).Return(nil)

	store.On("Get", mock.Anything, tokenID).Return(&token.Token{
		ID:         tokenID,
		AccountID:  accountID,
		Token:      "abc123",
		TokenHash:  []byte("abc123"),
		ReadOnly:   false,
		CreatedOn:  now,
		ModifiedOn: now,
	}, nil)
	provisioner := provisioner2.NewMockProvisioner()

	provisioner.On("DeprovisionToken", mock.Anything, &token.Token{
		ID:         "",
		AccountID:  "",
		Token:      "",
		TokenHash:  []byte("abc123"),
		ReadOnly:   false,
		CreatedOn:  time.Time{},
		ModifiedOn: time.Time{},
	}).Return(nil)

	app := &API{
		Token:       store,
		Provisioner: provisioner,
	}
	r := chi.NewRouter()
	r.Delete("/accounts/{accountId}/tokens/{tokenId}", app.DeleteToken())
	r.ServeHTTP(rr, req)
	assert.Equalf(t, http.StatusOK, rr.Code, "should return ok")
	store.AssertExpectations(t)
}
