package api

import (
	"encoding/json"
	"github.com/broswen/vex/internal/project"
	"github.com/broswen/vex/internal/provisioner"
	"github.com/broswen/vex/internal/stats"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func CreateProject(projectStore project.ProjectStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId := chi.URLParam(r, "accountId")
		p := &project.Project{}
		err := json.NewDecoder(r.Body).Decode(p)
		if err != nil {
			RenderError(w, ErrBadRequest.WithError(err))
			return
		}
		defer r.Body.Close()
		p.AccountID = accountId
		err = projectStore.Insert(r.Context(), p)

		if err != nil {
			RenderError(w, err)
			return
		}
		stats.ProjectCreated.Inc()
		err = json.NewEncoder(w).Encode(p)
		if err != nil {
			RenderError(w, err)
			return
		}
	}
}

func UpdateProject(projectStore project.ProjectStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId := chi.URLParam(r, "accountId")
		projectId := chi.URLParam(r, "projectId")
		p := &project.Project{}
		err := json.NewDecoder(r.Body).Decode(p)
		if err != nil {
			RenderError(w, ErrBadRequest.WithError(err))
			return
		}
		defer r.Body.Close()

		p.ID = projectId
		p.AccountID = accountId

		err = projectStore.Update(r.Context(), p)

		if err != nil {
			RenderError(w, err)
			return
		}

		err = json.NewEncoder(w).Encode(p)
		if err != nil {
			RenderError(w, err)
			return
		}
	}
}

func ListProjects(projectStore project.ProjectStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId := chi.URLParam(r, "accountId")
		projects, err := projectStore.List(r.Context(), accountId, 100, 0)
		if err != nil {
			RenderError(w, err)
			return
		}
		err = json.NewEncoder(w).Encode(projects)
		if err != nil {
			RenderError(w, err)
			return
		}
	}
}

func GetProject(projectStore project.ProjectStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId := chi.URLParam(r, "accountId")
		projectId := chi.URLParam(r, "projectId")
		p, err := projectStore.Get(r.Context(), projectId, accountId)
		if err != nil {
			RenderError(w, err)
			return
		}
		err = json.NewEncoder(w).Encode(p)
		if err != nil {
			RenderError(w, err)
			return
		}
	}
}

func DeleteProject(projectStore project.ProjectStore, provisioner provisioner.Provisioner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectId := chi.URLParam(r, "projectId")
		accountId := chi.URLParam(r, "accountId")
		err := projectStore.Delete(r.Context(), projectId, accountId)
		if err != nil {
			RenderError(w, err)
			return
		}
		err = provisioner.DeprovisionProject(r.Context(), &project.Project{ID: projectId})
		if err != nil {
			log.Printf("deprovision %s: %v", projectId, err)
		}
		stats.ProjectDeleted.Inc()
		err = json.NewEncoder(w).Encode(&struct{ id string }{id: projectId})
		if err != nil {
			RenderError(w, err)
			return
		}
	}
}
