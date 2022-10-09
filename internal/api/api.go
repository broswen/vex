package api

import (
	"github.com/broswen/vex/internal/account"
	"github.com/broswen/vex/internal/flag"
	"github.com/broswen/vex/internal/project"
	"github.com/broswen/vex/internal/provisioner"
	"github.com/broswen/vex/internal/token"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type API struct {
	Account     account.Store
	Project     project.Store
	Flag        flag.Store
	Token       token.Store
	Provisioner provisioner.Provisioner
}

func (api *API) AdminRouter(teamDomain, policyAUD string) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Heartbeat("/healthz"))
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Use(CloudflareAccessVerifier(teamDomain, policyAUD))

	r.Get("/admin/accounts", api.ListAccounts())
	r.Post("/admin/accounts", api.CreateAccount())

	r.Post("/admin/accounts/{accountId}/tokens", api.GenerateToken())
	r.Put("/admin/accounts/{accountId}/tokens/{tokenId}", api.RerollToken())

	return r
}

func (api *API) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Heartbeat("/healthz"))
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		writeErr(w, nil, ErrNotFound)
	})

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		writeOK(w, http.StatusOK, "OK")
	})

	//disable creating accounts through api for now
	//r.Post("/accounts", CreateAccount(accountStore))
	//r.Get("/accounts/", http.NotFound)
	r.Route("/accounts/{accountId}", func(r chi.Router) {
		r.Use(AccountAuthorizer(api.Token))
		r.Get("/", api.GetAccount())
		r.Put("/", api.UpdateAccount())
		r.Delete("/", api.DeleteAccount())

		r.Get("/tokens", api.ListTokens())
		r.Post("/tokens", api.GenerateToken())
		r.Put("/tokens/{tokenId}", api.RerollToken())
		r.Delete("/tokens/{tokenId}", api.DeleteToken())

		r.Post("/projects", api.CreateProject())
		r.Get("/projects", api.ListProjects())
		r.Put("/projects/{projectId}", api.UpdateProject())
		r.Get("/projects/{projectId}", api.GetProject())
		r.Delete("/projects/{projectId}", api.DeleteProject())

		r.Post("/projects/{projectId}/flags", api.CreateFlag())
		r.Get("/projects/{projectId}/flags", api.ListFlags())
		r.Put("/projects/{projectId}/flags/{flagId}", api.UpdateFlag())
		r.Get("/projects/{projectId}/flags/{flagId}", api.GetFlag())
		r.Delete("/projects/{projectId}/flags/{flagId}", api.DeleteFlag())
	})

	return r
}
