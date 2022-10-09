package api

import (
	"net/http"

	"github.com/broswen/vex/internal/project"
	"github.com/broswen/vex/internal/stats"
	"github.com/rs/zerolog/log"
)

func (api *API) CreateProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId, err := accountId(r)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		p := &project.Project{}
		err = readJSON(w, r, p)
		if err != nil {
			writeErr(w, nil, ErrBadRequest.WithError(err))
			return
		}
		defer r.Body.Close()
		p.AccountID = accountId
		newProject, err := api.Project.Insert(r.Context(), p)

		if err != nil {
			writeErr(w, nil, err)
			return
		}
		stats.ProjectCreated.Inc()
		err = writeOK(w, http.StatusOK, newProject)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
	}
}

func (api *API) UpdateProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId, err := accountId(r)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		projectId, err := projectId(r)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		p := &project.Project{}
		err = readJSON(w, r, p)
		if err != nil {
			writeErr(w, nil, ErrBadRequest.WithError(err))
			return
		}
		defer r.Body.Close()

		p.ID = projectId
		p.AccountID = accountId

		updatedProject, err := api.Project.Update(r.Context(), p)

		if err != nil {
			writeErr(w, nil, err)
			return
		}
		err = writeOK(w, http.StatusOK, updatedProject)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
	}
}

func (api *API) ListProjects() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountId, err := accountId(r)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		p := pagination(r)
		projects, err := api.Project.List(r.Context(), accountId, p.Limit, p.Offset)
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

func (api *API) GetProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectId, err := projectId(r)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		p, err := api.Project.Get(r.Context(), projectId)
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

func (api *API) DeleteProject() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectId, err := projectId(r)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		err = api.Project.Delete(r.Context(), projectId)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		err = api.Provisioner.DeprovisionProject(r.Context(), &project.Project{ID: projectId})
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
