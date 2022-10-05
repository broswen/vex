package api

import (
	"github.com/broswen/vex/internal/flag"
	"github.com/broswen/vex/internal/stats"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (api *API) CreateFlag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectId := chi.URLParam(r, "projectId")
		p, err := api.Project.Get(r.Context(), projectId)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		accountId := chi.URLParam(r, "accountId")
		f := &flag.Flag{}
		err = readJSON(w, r, f)
		if err != nil {
			writeErr(w, nil, ErrBadRequest.WithError(err))
			return
		}
		defer r.Body.Close()
		f.ProjectID = p.ID
		f.AccountID = accountId

		if err = flag.Validate(*f); err != nil {
			writeErr(w, nil, ErrBadRequest.WithError(err))
			return
		}

		err = api.Flag.Insert(r.Context(), f)

		if err != nil {
			writeErr(w, nil, err)
			return
		}

		err = api.Provisioner.ProvisionProject(r.Context(), p)
		if err != nil {
			log.Warn().Str("id", projectId).Err(err).Msg("could not provision project")
		}

		stats.FlagCreated.Inc()

		err = writeOK(w, http.StatusOK, f)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
	}
}

func (api *API) UpdateFlag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flagId, err := flagId(r)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		projectId := chi.URLParam(r, "projectId")
		p, err := api.Project.Get(r.Context(), projectId)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		accountId := chi.URLParam(r, "accountId")
		f := &flag.Flag{}
		err = readJSON(w, r, f)
		if err != nil {
			writeErr(w, nil, ErrBadRequest.WithError(err))
			return
		}
		defer r.Body.Close()

		f.ID = flagId
		f.ProjectID = p.ID
		f.AccountID = accountId

		if err = flag.Validate(*f); err != nil {
			writeErr(w, nil, err)
			return
		}

		err = api.Flag.Update(r.Context(), f)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		err = api.Provisioner.ProvisionProject(r.Context(), p)
		if err != nil {
			log.Warn().Str("id", projectId).Err(err).Msg("could not provision project")
		}

		stats.FlagUpdated.Inc()

		err = writeOK(w, http.StatusOK, f)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
	}
}

func (api *API) ListFlags() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectId := chi.URLParam(r, "projectId")
		project, err := api.Project.Get(r.Context(), projectId)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		p := pagination(r)
		flags, err := api.Flag.List(r.Context(), project.ID, p.Limit, p.Offset)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		err = writeOK(w, http.StatusOK, flags)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
	}
}

func (api *API) GetFlag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//projectId := chi.URLParam(r, "projectId")
		flagId, err := flagId(r)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		f, err := api.Flag.Get(r.Context(), flagId)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		err = writeOK(w, http.StatusOK, f)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
	}
}

func (api *API) DeleteFlag() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectId := chi.URLParam(r, "projectId")
		p, err := api.Project.Get(r.Context(), projectId)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		flagId, err := flagId(r)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		err = api.Flag.Delete(r.Context(), flagId)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		err = api.Provisioner.ProvisionProject(r.Context(), p)
		if err != nil {
			log.Warn().Str("id", projectId).Err(err).Msg("could not provision project")
		}

		stats.FlagDeleted.Inc()

		err = writeOK(w, http.StatusOK, &struct{ id string }{id: flagId})
		if err != nil {
			writeErr(w, nil, err)
			return
		}
	}
}
