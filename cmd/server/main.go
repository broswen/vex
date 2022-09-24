package main

import (
	"context"
	"fmt"
	"github.com/broswen/vex/internal/account"
	"github.com/broswen/vex/internal/api"
	"github.com/broswen/vex/internal/db"
	"github.com/broswen/vex/internal/flag"
	"github.com/broswen/vex/internal/project"
	"github.com/broswen/vex/internal/provisioner"
	"github.com/broswen/vex/internal/token"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// postgres connection string
	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatal().Msg("postgres DSN is empty")
	}
	// port for api
	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		apiPort = "8080"
	}
	// port for admin api
	adminPort := os.Getenv("ADMIN_PORT")
	if apiPort == "" {
		apiPort = "8082"
	}
	// port for prometheus
	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "8081"
	}

	metricsPath := os.Getenv("METRICS_PATH")
	if metricsPath == "" {
		metricsPath = "/metrics"
	}

	provisionTopic := os.Getenv("PROVISION_TOPIC")
	if provisionTopic == "" {
		provisionTopic = "vex-provision"
	}
	deprovisionTopic := os.Getenv("DEPROVISION_TOPIC")
	if deprovisionTopic == "" {
		deprovisionTopic = "vex-deprovision"
	}
	tokenProvisionTopic := os.Getenv("TOKEN_PROVISION_TOPIC")
	if tokenProvisionTopic == "" {
		tokenProvisionTopic = "vex-provision-token"
	}
	tokenDeprovisionTopic := os.Getenv("TOKEN_DEPROVISION_TOPIC")
	if tokenDeprovisionTopic == "" {
		tokenDeprovisionTopic = "vex-deprovision-token"
	}

	brokers := os.Getenv("BROKERS")
	if brokers == "" {
		brokers = "kafka-clusterip.kafka.svc.cluster.local:9092"
	}

	//	initialize database connection
	database, err := db.InitDB(context.Background(), dsn)
	if err != nil {
		log.Fatal().Err(err)
	}
	projectStore, err := project.NewPostgresStore(database)
	if err != nil {
		log.Fatal().Err(err)
	}
	flagStore, err := flag.NewPostgresStore(database)
	if err != nil {
		log.Fatal().Err(err)
	}
	accountStore, err := account.NewPostgresStore(database)
	if err != nil {
		log.Fatal().Err(err)
	}
	tokenStore, err := token.NewPostgresStore(database)
	if err != nil {
		log.Fatal().Err(err)
	}

	provisioner, err := provisioner.NewKafkaProvisioner(provisionTopic, deprovisionTopic, tokenProvisionTopic, tokenDeprovisionTopic, brokers)
	if err != nil {
		log.Fatal().Err(err)
	}

	// start promhttp listener on metrics port
	m := chi.NewRouter()
	m.Handle(metricsPath, promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%s", metricsPort), m); err != nil {
			log.Fatal().Err(err)
		}
	}()

	eg := errgroup.Group{}

	adminRouter := api.AdminRouter(accountStore, tokenStore, provisioner)
	// start admin api listener on api port
	eg.Go(func() error {
		log.Debug().Msgf("admin api listening on :%s", adminPort)
		return http.ListenAndServe(fmt.Sprintf(":%s", adminPort), adminRouter)
	})

	// start api listener on api port
	publicRouter := api.Router(projectStore, flagStore, accountStore, tokenStore, provisioner)
	eg.Go(func() error {
		log.Debug().Msgf("public api listening on :%s", apiPort)
		return http.ListenAndServe(fmt.Sprintf(":%s", apiPort), publicRouter)
	})

	if err := eg.Wait(); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
