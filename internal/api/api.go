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

func Router(projectStore project.ProjectStore, flagStore flag.FlagStore, accountStore account.AccountStore, tokenStore token.TokenStore, provisioner provisioner.Provisioner) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Heartbeat("/healthz"))
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)

	//disable creating accounts through api for now
	//r.Post("/accounts", CreateAccount(accountStore))
	r.Get("/accounts/", http.NotFound)
	r.Get("/accounts/{accountId}", AccountAuthorizer(GetAccount(accountStore), tokenStore))
	r.Put("/accounts/{accountId}", AccountAuthorizer(UpdateAccount(accountStore), tokenStore))
	r.Delete("/accounts/{accountId}", AccountAuthorizer(DeleteAccount(accountStore), tokenStore))

	r.Post("/accounts/{accountId}/tokens", AccountAuthorizer(GenerateToken(tokenStore), tokenStore))
	r.Put("/accounts/{accountId}/tokens/{tokenId}", AccountAuthorizer(RerollToken(tokenStore), tokenStore))
	r.Delete("/accounts/{accountId}/tokens/{tokenId}", AccountAuthorizer(DeleteToken(tokenStore), tokenStore))

	r.Post("/accounts/{accountId}/projects", AccountAuthorizer(CreateProject(projectStore), tokenStore))
	r.Get("/accounts/{accountId}/projects", AccountAuthorizer(ListProjects(projectStore), tokenStore))
	r.Put("/accounts/{accountId}/projects/{projectId}", AccountAuthorizer(UpdateProject(projectStore), tokenStore))
	r.Get("/accounts/{accountId}/projects/{projectId}", AccountAuthorizer(GetProject(projectStore), tokenStore))
	r.Delete("/accounts/{accountId}/projects/{projectId}", AccountAuthorizer(DeleteProject(projectStore, provisioner), tokenStore))

	r.Post("/accounts/{accountId}/projects/{projectId}/flags", AccountAuthorizer(CreateFlag(flagStore, provisioner), tokenStore))
	r.Get("/accounts/{accountId}/projects/{projectId}/flags", AccountAuthorizer(ListFlags(flagStore), tokenStore))
	r.Put("/accounts/{accountId}/projects/{projectId}/flags/{flagId}", AccountAuthorizer(UpdateFlag(flagStore, provisioner), tokenStore))
	r.Get("/accounts/{accountId}/projects/{projectId}/flags/{flagId}", AccountAuthorizer(GetFlag(flagStore), tokenStore))
	r.Delete("/accounts/{accountId}/projects/{projectId}/flags/{flagId}", AccountAuthorizer(DeleteFlag(flagStore, provisioner), tokenStore))
	return r
}
