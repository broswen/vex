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
	"log"
	"net/http"
	"os"
)

func main() {
	// postgres connection string
	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatalf("postgres DSN is empty")
	}
	// port for api
	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		apiPort = "8080"
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
		log.Fatal(err)
	}
	projectStore, err := project.NewPostgresStore(database)
	if err != nil {
		log.Fatal(err)
	}
	flagStore, err := flag.NewPostgresStore(database)
	if err != nil {
		log.Fatal(err)
	}
	accountStore, err := account.NewPostgresStore(database)
	if err != nil {
		log.Fatal(err)
	}
	tokenStore, err := token.NewPostgresStore(database)
	if err != nil {
		log.Fatal(err)
	}

	provisioner, err := provisioner.NewKafkaProvisioner(provisionTopic, deprovisionTopic, tokenProvisionTopic, tokenDeprovisionTopic, brokers)
	if err != nil {
		log.Fatal(err)
	}

	// start promhttp listener on metrics port
	m := chi.NewRouter()
	m.Handle(metricsPath, promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%s", metricsPort), m); err != nil {
			log.Fatal(err)
		}
	}()

	// start api listener on api port
	r := api.Router(projectStore, flagStore, accountStore, tokenStore, provisioner)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", apiPort), r); err != nil {
		log.Fatal(err)
	}

	//	listen for os signals?

}
