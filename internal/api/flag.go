package api

import (
	"github.com/broswen/vex/internal/flag"
	"github.com/broswen/vex/internal/project"
	"github.com/broswen/vex/internal/provisioner"
	"github.com/broswen/vex/internal/stats"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"net/http"
)

func CreateFlag(flagStore flag.FlagStore, projectStore project.ProjectStore, provisioner provisioner.Provisioner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectId := chi.URLParam(r, "projectId")
		p, err := projectStore.Get(r.Context(), projectId)
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

		err = flagStore.Insert(r.Context(), f)

		if err != nil {
			writeErr(w, nil, err)
			return
		}

		err = provisioner.ProvisionProject(r.Context(), p)
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

func UpdateFlag(flagStore flag.FlagStore, projectStore project.ProjectStore, provisioner provisioner.Provisioner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flagId := chi.URLParam(r, "flagId")
		projectId := chi.URLParam(r, "projectId")
		p, err := projectStore.Get(r.Context(), projectId)
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

		err = flagStore.Update(r.Context(), f)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		err = provisioner.ProvisionProject(r.Context(), p)
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

func ListFlags(flagStore flag.FlagStore, projectStore project.ProjectStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectId := chi.URLParam(r, "projectId")
		project, err := projectStore.Get(r.Context(), projectId)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		flags, err := flagStore.List(r.Context(), project.ID, 100, 0)
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

func GetFlag(flagStore flag.FlagStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//projectId := chi.URLParam(r, "projectId")
		flagId := chi.URLParam(r, "flagId")
		f, err := flagStore.Get(r.Context(), flagId)
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

func DeleteFlag(flagStore flag.FlagStore, projectStore project.ProjectStore, provisioner provisioner.Provisioner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectId := chi.URLParam(r, "projectId")
		p, err := projectStore.Get(r.Context(), projectId)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		flagId := chi.URLParam(r, "flagId")
		err = flagStore.Delete(r.Context(), flagId)
		if err != nil {
			writeErr(w, nil, err)
			return
		}
		err = provisioner.ProvisionProject(r.Context(), p)
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
