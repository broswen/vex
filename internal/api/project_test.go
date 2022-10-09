package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/broswen/vex/internal/project"
	provisioner2 "github.com/broswen/vex/internal/provisioner"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var projectID = "da863681-2f59-432d-848d-a64fbfbeab38"

func TestCreateProjectHandler(t *testing.T) {
	p1 := &project.Project{
		AccountID:   accountID,
		Name:        "test",
		Description: "test project",
	}
	reqBody, err := json.Marshal(p1)
	assert.Nil(t, err)
	req, err := http.NewRequest(http.MethodPost, "/accounts/"+accountID+"/projects", bytes.NewReader(reqBody))
	assert.Nil(t, err)
	req.WithContext(context.Background())
	rr := httptest.NewRecorder()
	store := project.NewMockStore()
	store.On("Insert", mock.Anything, p1).Return(&project.Project{
		ID:          projectID,
		AccountID:   accountID,
		Name:        "test",
		Description: "test project",
		CreatedOn:   now,
		ModifiedOn:  now,
	}, nil)
	app := &API{
		Project: store,
	}
	r := chi.NewRouter()
	r.Post("/accounts/{accountId}/projects", app.CreateProject())
	r.ServeHTTP(rr, req)
	assert.Equalf(t, http.StatusOK, rr.Code, "should return ok")
	store.AssertExpectations(t)
}

func TestGetProjectHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/accounts/"+accountID+"/projects/"+projectID, nil)
	assert.Nil(t, err)
	req.WithContext(context.Background())
	rr := httptest.NewRecorder()
	store := project.NewMockStore()
	store.On("Get", mock.Anything, projectID).Return(&project.Project{
		ID:          projectID,
		AccountID:   accountID,
		Name:        "test",
		Description: "test project",
		CreatedOn:   now,
		ModifiedOn:  now,
	}, nil)
	app := &API{
		Project: store,
	}
	r := chi.NewRouter()
	r.Get("/accounts/{accountId}/projects/{projectId}", app.GetProject())
	r.ServeHTTP(rr, req)
	assert.Equalf(t, http.StatusOK, rr.Code, "should return ok")
	store.AssertExpectations(t)
}

func TestListProjectsHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/accounts/"+accountID+"/projects", nil)
	assert.Nil(t, err)
	req.WithContext(context.Background())
	rr := httptest.NewRecorder()
	store := project.NewMockStore()
	store.On("List", mock.Anything, accountID, int64(100), int64(0)).Return([]*project.Project{
		{
			ID:          projectID,
			AccountID:   accountID,
			Name:        "test",
			Description: "test project",
			CreatedOn:   now,
			ModifiedOn:  now,
		},
	}, nil)
	app := &API{
		Project: store,
	}
	r := chi.NewRouter()
	r.Get("/accounts/{accountId}/projects", app.ListProjects())
	r.ServeHTTP(rr, req)
	assert.Equalf(t, http.StatusOK, rr.Code, "should return ok")
	store.AssertExpectations(t)
}

func TestUpdateProjectHandler(t *testing.T) {
	p1 := &project.Project{
		Name:        "test",
		Description: "test project",
	}
	reqBody, err := json.Marshal(p1)
	assert.Nil(t, err)
	req, err := http.NewRequest(http.MethodPut, "/accounts/"+accountID+"/projects/"+projectID, bytes.NewReader(reqBody))
	assert.Nil(t, err)
	req.WithContext(context.Background())
	rr := httptest.NewRecorder()
	store := project.NewMockStore()
	store.On("Update", mock.Anything, &project.Project{
		ID:          projectID,
		AccountID:   accountID,
		Name:        "test",
		Description: "test project",
	}).Return(&project.Project{
		ID:          projectID,
		AccountID:   accountID,
		Name:        "test",
		Description: "test project",
		CreatedOn:   now,
		ModifiedOn:  now,
	}, nil)
	app := &API{
		Project: store,
	}
	r := chi.NewRouter()
	r.Put("/accounts/{accountId}/projects/{projectId}", app.UpdateProject())
	r.ServeHTTP(rr, req)
	assert.Equalf(t, http.StatusOK, rr.Code, "should return ok")
	store.AssertExpectations(t)
}

func TestDeleteProjectHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "/accounts/"+accountID+"/projects/"+projectID, nil)
	assert.Nil(t, err)
	req.WithContext(context.Background())
	rr := httptest.NewRecorder()
	store := project.NewMockStore()
	store.On("Delete", mock.Anything, projectID).Return(nil)
	provisioner := provisioner2.NewMockProvisioner()
	provisioner.On("DeprovisionProject", mock.Anything, &project.Project{ID: projectID}).Return(nil)
	app := &API{
		Project:     store,
		Provisioner: provisioner,
	}
	r := chi.NewRouter()
	r.Delete("/accounts/{accountId}/projects/{projectId}", app.DeleteProject())
	r.ServeHTTP(rr, req)
	assert.Equalf(t, http.StatusOK, rr.Code, "should return ok")
	store.AssertExpectations(t)
}
