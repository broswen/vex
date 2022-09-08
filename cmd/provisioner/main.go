package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/broswen/vex/internal/db"
	flag2 "github.com/broswen/vex/internal/flag"
	"github.com/broswen/vex/internal/project"
	"github.com/broswen/vex/internal/provisioner"
	"github.com/broswen/vex/internal/stats"
	"github.com/broswen/vex/internal/token"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

func main() {
	// cloudflare api token to manage KV
	cloudflareToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	// cloudflare api token to manage KV
	cloudflareAccountId := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	// KV namespace id
	projectKVNamespaceID := os.Getenv("PROJECT_KV_NAMESPACE_ID")
	tokenKVNamespaceID := os.Getenv("TOKEN_KV_NAMESPACE_ID")

	skipProvision := os.Getenv("SKIP_PROVISION")
	if skipProvision == "true" {
		log.Printf("SKIP_PROVISION=%s", skipProvision)
	}

	// postgres connection string
	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatalf("postgres DSN is empty")
	}

	//	initialize database connection
	database, err := db.InitDB(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	flagStore, err := flag2.NewPostgresStore(database)
	if err != nil {
		log.Fatal(err)
	}
	tokenStore, err := token.NewPostgresStore(database)
	if err != nil {
		log.Fatal(err)
	}

	cloudflareProvisioner, err := provisioner.NewCloudflareProvisioner(cloudflareToken, cloudflareAccountId, projectKVNamespaceID, tokenKVNamespaceID, flagStore, tokenStore)

	// port for prometheus
	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "8081"
	}

	metricsPath := os.Getenv("METRICS_PATH")
	if metricsPath == "" {
		metricsPath = "/metrics"
	}

	group := os.Getenv("GROUP")
	if group == "" {
		group = "vex-cloudflareProvisioner"
	}
	topics := os.Getenv("TOPICS")
	if topics == "" {
		topics = "vex-provision,vex-deprovision,vex-provision-token,vex-deprovision-token"
	}
	brokers := os.Getenv("BROKERS")
	if brokers == "" {
		brokers = "kafka-clusterip.kafka.svc.cluster.local:9092"
	}
	flag.Parse()

	config := sarama.NewConfig()
	config.ClientID = "vex-cloudflareProvisioner"
	version, err := sarama.ParseKafkaVersion("3.1.0")
	if err != nil {
		log.Fatal(err)
	}
	config.Version = version

	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

	consumer := NewConsumer(skipProvision == "true")
	consumer.HandleFunc("vex-provision", func(message *sarama.ConsumerMessage) error {
		log.Printf("provisioning %s", string(message.Value))
		stats.ProjectProvisioned.Inc()
		return cloudflareProvisioner.ProvisionProject(context.Background(), &project.Project{ID: string(message.Value)})
	})
	consumer.HandleFunc("vex-deprovision", func(message *sarama.ConsumerMessage) error {
		log.Printf("deprovisioning %s", string(message.Value))
		stats.ProjectDeprovisioned.Inc()
		return cloudflareProvisioner.DeprovisionProject(context.Background(), &project.Project{ID: string(message.Value)})
	})

	consumer.HandleFunc("vex-provision-token", func(message *sarama.ConsumerMessage) error {
		log.Printf("provisioning %s", string(message.Value))
		stats.ProjectProvisioned.Inc()
		return cloudflareProvisioner.ProvisionToken(context.Background(), &token.Token{ID: string(message.Value)})
	})
	consumer.HandleFunc("vex-deprovision-token", func(message *sarama.ConsumerMessage) error {
		log.Printf("deprovisioning %s", string(message.Value))
		stats.ProjectDeprovisioned.Inc()
		return cloudflareProvisioner.DeprovisionToken(context.Background(), &token.Token{Token: string(message.Value)})
	})

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), group, config)
	if err != nil {
		log.Panic(err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	log.Println("starting consume loop")
	go func() {
		defer wg.Done()

		for {
			if err := client.Consume(ctx, strings.Split(topics, ","), consumer); err != nil {
				log.Panic(err.Error())
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-consumer.ready

	// start promhttp listener on metrics port
	m := chi.NewRouter()
	m.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%s", metricsPort), m); err != nil {
			log.Fatal(err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	select {
	case s := <-sigs:
		log.Println("received signal:", s)
	}
	cancel()
	wg.Wait()
	if err := client.Close(); err != nil {
		log.Println(err.Error())
	}
}

func NewConsumer(skip bool) *Consumer {
	return &Consumer{
		skip:     skip,
		ready:    make(chan bool),
		handlers: make(map[string]MessageHandler),
	}
}

type MessageHandler func(message *sarama.ConsumerMessage) error

type Consumer struct {
	// whether in dev mode and shouldn't actually provision anything
	skip     bool
	ready    chan bool
	handlers map[string]MessageHandler
}

func (c *Consumer) Handle(message *sarama.ConsumerMessage) error {
	for k, f := range c.handlers {
		if message.Topic == k {
			if c.skip {
				log.Printf("skipping message %s %s", message.Topic, message.Value)
				return nil
			}
			return f(message)
		}
	}
	return errors.New("no matching handler: " + message.Topic)
}

func (c *Consumer) HandleFunc(topic string, f MessageHandler) {
	c.handlers[topic] = f
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			err := c.Handle(message)
			if err != nil {
				log.Println(err)
				continue
			}
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
