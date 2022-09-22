package api

import (
	"github.com/broswen/vex/internal/project"
	"github.com/broswen/vex/internal/provisioner"
	"github.com/broswen/vex/internal/stats"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"net/http"
)

func CreateProject(projectStore project.ProjectStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId := chi.URLParam(r, "accountId")
		p := &project.Project{}
		err := readJSON(w, r, p)
		if err != nil {
			writeErr(w, nil, ErrBadRequest.WithError(err))
			return
		}
		defer r.Body.Close()
		p.AccountID = accountId
		err = projectStore.Insert(r.Context(), p)

		if err != nil {
			writeErr(w, nil, err)
			return
		}
		stats.ProjectCreated.Inc()
		err = writeOK(w, http.StatusOK, p)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
	}
}

func UpdateProject(projectStore project.ProjectStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId := chi.URLParam(r, "accountId")
		projectId := chi.URLParam(r, "projectId")
		p := &project.Project{}
		err := readJSON(w, r, p)
		if err != nil {
			writeErr(w, nil, ErrBadRequest.WithError(err))
			return
		}
		defer r.Body.Close()

		p.ID = projectId
		p.AccountID = accountId

		err = projectStore.Update(r.Context(), p)

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

func ListProjects(projectStore project.ProjectStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId := chi.URLParam(r, "accountId")
		projects, err := projectStore.List(r.Context(), accountId, 100, 0)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		err = writeOK(w, http.StatusOK, projects)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
	}
}

func GetProject(projectStore project.ProjectStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//accountId := chi.URLParam(r, "accountId")
		projectId := chi.URLParam(r, "projectId")
		p, err := projectStore.Get(r.Context(), projectId)
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

func DeleteProject(projectStore project.ProjectStore, provisioner provisioner.Provisioner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectId := chi.URLParam(r, "projectId")
		//accountId := chi.URLParam(r, "accountId")
		err := projectStore.Delete(r.Context(), projectId)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		err = provisioner.DeprovisionProject(r.Context(), &project.Project{ID: projectId})
		if err != nil {
			log.Warn().Str("id", projectId).Err(err).Msg("could not deprovision project")
		}
		stats.ProjectDeleted.Inc()
		err = writeOK(w, http.StatusOK, &struct{ id string }{id: projectId})
		if err != nil {
			writeErr(w, nil, err)
			return
		}
	}
}
