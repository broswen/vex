package stats

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ProvisionError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "provision_error",
	})

	ProjectProvisioned = promauto.NewCounter(prometheus.CounterOpts{
		Name: "project_provisioned",
	})

	DeprovisionError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "deprovision_error",
	})

	ProjectDeprovisioned = promauto.NewCounter(prometheus.CounterOpts{
		Name: "project_deprovisioned",
	})

	ProjectCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "project_created",
	})

	ProjectDeleted = promauto.NewCounter(prometheus.CounterOpts{
		Name: "project_deleted",
	})

	FlagCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "flag_created",
	})

	FlagUpdated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "flag_updated",
	})

	FlagDeleted = promauto.NewCounter(prometheus.CounterOpts{
		Name: "flag_deleted",
	})
)
