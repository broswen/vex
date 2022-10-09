package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/broswen/vex/internal/flag"
	"github.com/broswen/vex/internal/project"
	provisioner2 "github.com/broswen/vex/internal/provisioner"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var flagID = "ke863682-2f60-431d-846d-a66fbfbeab39"

func TestCreateFlagHandler(t *testing.T) {
	f1 := &flag.Flag{
		Key:   "flag1",
		Type:  "STRING",
		Value: "test",
	}
	reqBody, err := json.Marshal(f1)
	assert.Nil(t, err)
	req, err := http.NewRequest(http.MethodPost, "/accounts/"+accountID+"/projects/"+projectID+"/flags", bytes.NewReader(reqBody))
	assert.Nil(t, err)
	req.WithContext(context.Background())
	rr := httptest.NewRecorder()
	p1 := &project.Project{
		ID:          projectID,
		AccountID:   accountID,
		Name:        "test",
		Description: "test",
		CreatedOn:   time.Time{},
		ModifiedOn:  time.Time{},
	}
	projectStore := project.NewMockStore()
	projectStore.On("Get", mock.Anything, projectID).Return(p1, nil)
	store := flag.NewMockStore()
	store.On("Insert", mock.Anything, &flag.Flag{
		ProjectID: projectID,
		AccountID: accountID,
		Key:       "flag1",
		Type:      flag.STRING,
		Value:     "test",
	}).Return(&flag.Flag{
		ID:         flagID,
		ProjectID:  projectID,
		AccountID:  accountID,
		Key:        "flag1",
		Type:       flag.STRING,
		Value:      "test",
		CreatedOn:  now,
		ModifiedOn: now,
	}, nil)
	provisioner := provisioner2.NewMockProvisioner()
	provisioner.On("ProvisionProject", mock.Anything, p1).Return(nil)
	app := &API{
		Flag:        store,
		Project:     projectStore,
		Provisioner: provisioner,
	}
	r := chi.NewRouter()
	r.Post("/accounts/{accountId}/projects/{projectId}/flags", app.CreateFlag())
	r.ServeHTTP(rr, req)
	assert.Equalf(t, http.StatusOK, rr.Code, "should return ok")
	store.AssertExpectations(t)
}

func TestGetFlagHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/accounts/"+accountID+"/projects/"+projectID+"/flags/"+flagID, nil)
	assert.Nil(t, err)
	req.WithContext(context.Background())
	rr := httptest.NewRecorder()
	p1 := &project.Project{
		ID:          projectID,
		AccountID:   accountID,
		Name:        "test",
		Description: "test",
		CreatedOn:   time.Time{},
		ModifiedOn:  time.Time{},
	}
	projectStore := project.NewMockStore()
	projectStore.On("Get", mock.Anything, projectID).Return(p1, nil)
	store := flag.NewMockStore()
	store.On("Get", mock.Anything, flagID).Return(&flag.Flag{
		ID:         flagID,
		ProjectID:  projectID,
		AccountID:  accountID,
		Key:        "flag1",
		Type:       flag.STRING,
		Value:      "test",
		CreatedOn:  now,
		ModifiedOn: now,
	}, nil)
	app := &API{
		Flag:    store,
		Project: projectStore,
	}
	r := chi.NewRouter()
	r.Get("/accounts/{accountId}/projects/{projectId}/flags/{flagId}", app.GetFlag())
	r.ServeHTTP(rr, req)
	assert.Equalf(t, http.StatusOK, rr.Code, "should return ok")
	store.AssertExpectations(t)
}

func TestListFlagHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/accounts/"+accountID+"/projects/"+projectID+"/flags", nil)
	assert.Nil(t, err)
	req.WithContext(context.Background())
	rr := httptest.NewRecorder()
	p1 := &project.Project{
		ID:          projectID,
		AccountID:   accountID,
		Name:        "test",
		Description: "test",
		CreatedOn:   time.Time{},
		ModifiedOn:  time.Time{},
	}
	projectStore := project.NewMockStore()
	projectStore.On("Get", mock.Anything, projectID).Return(p1, nil)
	store := flag.NewMockStore()
	store.On("List", mock.Anything, projectID, int64(100), int64(0)).Return([]*flag.Flag{
		{
			ID:         flagID,
			ProjectID:  projectID,
			AccountID:  accountID,
			Key:        "flag1",
			Type:       flag.STRING,
			Value:      "test",
			CreatedOn:  now,
			ModifiedOn: now,
		},
	}, nil)
	app := &API{
		Flag:    store,
		Project: projectStore,
	}
	r := chi.NewRouter()
	r.Get("/accounts/{accountId}/projects/{projectId}/flags", app.ListFlags())
	r.ServeHTTP(rr, req)
	assert.Equalf(t, http.StatusOK, rr.Code, "should return ok")
	store.AssertExpectations(t)
}

func TestUpdateFlagHandler(t *testing.T) {
	f1 := &flag.Flag{
		Key:   "flag1",
		Type:  "STRING",
		Value: "test",
	}
	reqBody, err := json.Marshal(f1)
	assert.Nil(t, err)
	req, err := http.NewRequest(http.MethodPut, "/accounts/"+accountID+"/projects/"+projectID+"/flags/"+flagID, bytes.NewReader(reqBody))
	assert.Nil(t, err)
	req.WithContext(context.Background())
	rr := httptest.NewRecorder()
	p1 := &project.Project{
		ID:          projectID,
		AccountID:   accountID,
		Name:        "test",
		Description: "test",
		CreatedOn:   time.Time{},
		ModifiedOn:  time.Time{},
	}
	projectStore := project.NewMockStore()
	projectStore.On("Get", mock.Anything, projectID).Return(p1, nil)
	store := flag.NewMockStore()
	store.On("Update", mock.Anything, &flag.Flag{
		ID:        flagID,
		ProjectID: projectID,
		AccountID: accountID,
		Key:       "flag1",
		Type:      flag.STRING,
		Value:     "test",
	}).Return(&flag.Flag{
		ID:         flagID,
		ProjectID:  projectID,
		AccountID:  accountID,
		Key:        "flag1",
		Type:       flag.STRING,
		Value:      "test",
		CreatedOn:  now,
		ModifiedOn: now,
	}, nil)
	provisioner := provisioner2.NewMockProvisioner()
	provisioner.On("ProvisionProject", mock.Anything, p1).Return(nil)
	app := &API{
		Flag:        store,
		Project:     projectStore,
		Provisioner: provisioner,
	}
	r := chi.NewRouter()
	r.Put("/accounts/{accountId}/projects/{projectId}/flags/{flagId}", app.UpdateFlag())
	r.ServeHTTP(rr, req)
	assert.Equalf(t, http.StatusOK, rr.Code, "should return ok")
	store.AssertExpectations(t)
}

func TestDeleteFlagHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "/accounts/"+accountID+"/projects/"+projectID+"/flags/"+flagID, nil)
	assert.Nil(t, err)
	req.WithContext(context.Background())
	rr := httptest.NewRecorder()
	p1 := &project.Project{
		ID:          projectID,
		AccountID:   accountID,
		Name:        "test",
		Description: "test",
		CreatedOn:   time.Time{},
		ModifiedOn:  time.Time{},
	}
	projectStore := project.NewMockStore()
	projectStore.On("Get", mock.Anything, projectID).Return(p1, nil)
	store := flag.NewMockStore()
	store.On("Delete", mock.Anything, flagID).Return(nil)
	provisioner := provisioner2.NewMockProvisioner()
	provisioner.On("ProvisionProject", mock.Anything, p1).Return(nil)
	app := &API{
		Flag:        store,
		Project:     projectStore,
		Provisioner: provisioner,
	}
	r := chi.NewRouter()
	r.Delete("/accounts/{accountId}/projects/{projectId}/flags/{flagId}", app.DeleteFlag())
	r.ServeHTTP(rr, req)
	assert.Equalf(t, http.StatusOK, rr.Code, "should return ok")
	store.AssertExpectations(t)
}
