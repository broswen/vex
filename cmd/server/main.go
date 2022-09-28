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
	"os/signal"
	"syscall"
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
	if adminPort == "" {
		adminPort = "8082"
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

	//Cloudflare Access Application policy AUD
	policyAUD := os.Getenv("POLICY_AUD")
	//Cloudflare Access team domain <team>.cloudflareaccess.com
	teamDomain := os.Getenv("TEAM_DOMAIN")

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

	eg := errgroup.Group{}

	// start promhttp listener on metrics port
	m := chi.NewRouter()
	m.Handle(metricsPath, promhttp.Handler())
	eg.Go(func() error {
		return http.ListenAndServe(fmt.Sprintf(":%s", metricsPort), m)
	})

	app := &api.API{
		Account:     accountStore,
		Project:     projectStore,
		Flag:        flagStore,
		Token:       tokenStore,
		Provisioner: provisioner,
	}

	adminRouter := app.AdminRouter(teamDomain, policyAUD)
	adminServer := http.Server{
		Addr:    fmt.Sprintf(":%s", adminPort),
		Handler: adminRouter,
	}
	// start admin api listener on api port
	eg.Go(func() error {
		log.Debug().Msgf("admin api listening on :%s", adminPort)
		if err := adminServer.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				return err
			}
		}
		return nil
	})

	publicRouter := app.Router()
	publicServer := http.Server{
		Addr:    fmt.Sprintf(":%s", apiPort),
		Handler: publicRouter,
	}
	// start public api listener on api port
	eg.Go(func() error {
		log.Debug().Msgf("public api listening on :%s", apiPort)
		if err := publicServer.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				return err
			}
		}
		return nil
	})

	eg.Go(func() error {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		sig := <-sigint
		log.Debug().Str("signal", sig.String()).Msg("received signal")

		var err error
		log.Debug().Msg("shutting down admin server")
		if er := adminServer.Shutdown(context.Background()); er != nil {
			log.Error().Err(er).Msg("error shutting down admin server")
			err = er
		}
		log.Debug().Msg("shutting down public server")
		if er := publicServer.Shutdown(context.Background()); er != nil {
			log.Error().Err(er).Msg("error shutting down public server")
			err = er
		}
		return err
	})

	if err := eg.Wait(); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
