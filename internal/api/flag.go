package api

import (
	"encoding/json"
	"github.com/broswen/vex/internal/flag"
	"github.com/broswen/vex/internal/project"
	"github.com/broswen/vex/internal/provisioner"
	"github.com/broswen/vex/internal/stats"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func CreateFlag(flagStore flag.FlagStore, provisioner provisioner.Provisioner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectId := chi.URLParam(r, "projectId")
		accountId := chi.URLParam(r, "accountId")
		f := &flag.Flag{}
		err := json.NewDecoder(r.Body).Decode(f)
		if err != nil {
			RenderError(w, ErrBadRequest.WithError(err))
			return
		}
		defer r.Body.Close()
		f.ProjectID = projectId
		f.AccountID = accountId

		if err = flag.Validate(*f); err != nil {
			RenderError(w, ErrBadRequest.WithError(err))
			return
		}

		err = flagStore.Insert(r.Context(), f)

		log.Printf("1: %#v\n", f)
		if err != nil {
			RenderError(w, err)
			return
		}

		err = provisioner.ProvisionProject(r.Context(), &project.Project{ID: projectId})
		if err != nil {
			log.Printf("provision %s: %v", projectId, err)
		}

		stats.FlagCreated.Inc()

		err = json.NewEncoder(w).Encode(f)
		if err != nil {
			RenderError(w, err)
			return
		}
	}
}

func UpdateFlag(flagStore flag.FlagStore, provisioner provisioner.Provisioner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flagId := chi.URLParam(r, "flagId")
		projectId := chi.URLParam(r, "projectId")
		accountId := chi.URLParam(r, "accountId")
		f := &flag.Flag{}
		err := json.NewDecoder(r.Body).Decode(f)
		if err != nil {
			RenderError(w, ErrBadRequest.WithError(err))
			return
		}
		defer r.Body.Close()

		f.ID = flagId
		f.ProjectID = projectId
		f.AccountID = accountId

		if err = flag.Validate(*f); err != nil {
			RenderError(w, err)
			return
		}

		err = flagStore.Update(r.Context(), f)
		if err != nil {
			RenderError(w, err)
			return
		}
		err = provisioner.ProvisionProject(r.Context(), &project.Project{ID: projectId})
		if err != nil {
			log.Printf("provision %s: %v", projectId, err)
		}

		stats.FlagUpdated.Inc()

		err = json.NewEncoder(w).Encode(f)
		if err != nil {
			RenderError(w, err)
			return
		}
	}
}

func ListFlags(flagStore flag.FlagStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectId := chi.URLParam(r, "projectId")
		flags, err := flagStore.List(r.Context(), projectId, 100, 0)
		if err != nil {
			RenderError(w, err)
			return
		}
		err = json.NewEncoder(w).Encode(flags)
		if err != nil {
			RenderError(w, err)
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
			RenderError(w, err)
			return
		}
		err = json.NewEncoder(w).Encode(f)
		if err != nil {
			RenderError(w, err)
			return
		}
	}
}

func DeleteFlag(flagStore flag.FlagStore, provisioner provisioner.Provisioner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectId := chi.URLParam(r, "projectId")
		flagId := chi.URLParam(r, "flagId")
		err := flagStore.Delete(r.Context(), flagId)
		if err != nil {
			RenderError(w, err)
			return
		}
		err = provisioner.ProvisionProject(r.Context(), &project.Project{ID: projectId})
		if err != nil {
			log.Printf("provision %s: %v", projectId, err)
		}

		stats.FlagDeleted.Inc()

		err = json.NewEncoder(w).Encode(&struct{ id string }{id: flagId})
		if err != nil {
			RenderError(w, err)
			return
		}
	}
}
