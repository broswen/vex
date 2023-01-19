package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/broswen/vex/internal/db"
	"github.com/broswen/vex/internal/flag"
	"github.com/broswen/vex/internal/project"
	"github.com/broswen/vex/internal/provisioner"
	"github.com/broswen/vex/internal/stats"
	"github.com/broswen/vex/internal/token"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	log2 "log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// cloudflare api token to manage KV
	cloudflareToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	// cloudflare api token to manage KV
	cloudflareAccountId := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	// KV namespace id
	projectKVNamespaceID := os.Getenv("PROJECT_KV_NAMESPACE_ID")
	flagKVNamespaceID := os.Getenv("FLAG_KV_NAMESPACE_ID")
	tokenKVNamespaceID := os.Getenv("TOKEN_KV_NAMESPACE_ID")

	skipProvision := os.Getenv("SKIP_PROVISION")
	if skipProvision == "true" {
		log.Debug().Msgf("SKIP_PROVISION=%s", skipProvision)
	}

	// postgres connection string
	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatal().Msgf("postgres DSN is empty")
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
	flagStore, err := flag2.NewPostgresStore(database)
	if err != nil {
		log.Fatal().Err(err)
	}
	tokenStore, err := token.NewPostgresStore(database)
	if err != nil {
		log.Fatal().Err(err)
	}

	cloudflareProvisioner, err := provisioner.NewCloudflareProvisioner(cloudflareToken, cloudflareAccountId, projectKVNamespaceID, tokenKVNamespaceID, flagKVNamespaceID, projectStore, flagStore, tokenStore)

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
		topics = "vex-provision,vex-deprovision,vex-provision-token,vex-deprovision-token,vex-provision-flag,vex-deprovision-flag"
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
		log.Fatal().Err(err)
	}
	config.Version = version

	sarama.Logger = log2.New(os.Stdout, "[sarama] ", log2.LstdFlags)

	consumer := NewConsumer(skipProvision == "true")
	consumer.HandleFunc("vex-provision", func(message *sarama.ConsumerMessage) error {
		log.Debug().Str("id", string(message.Value)).Msg("provisioning project")
		stats.ProjectProvisioned.Inc()
		return cloudflareProvisioner.ProvisionProject(context.Background(), &project.Project{ID: string(message.Value)})
	})
	consumer.HandleFunc("vex-deprovision", func(message *sarama.ConsumerMessage) error {
		log.Debug().Str("id", string(message.Value)).Msg("deprovisioning project")
		stats.ProjectDeprovisioned.Inc()
		return cloudflareProvisioner.DeprovisionProject(context.Background(), &project.Project{ID: string(message.Value)})
	})
	consumer.HandleFunc("vex-provision-flag", func(message *sarama.ConsumerMessage) error {
		f := &flag.Flag{}
		err := json.Unmarshal(message.Value, f)
		if err != nil {
			return err
		}
		log.Debug().Str("flag_id", f.ID).Msg("provisioning flag")
		stats.FlagProvisioned.Inc()
		return cloudflareProvisioner.ProvisionFlag(context.Background(), f)
	})
	consumer.HandleFunc("vex-deprovision-flag", func(message *sarama.ConsumerMessage) error {
		f := &flag.Flag{}
		err := json.Unmarshal(message.Value, f)
		if err != nil {
			return err
		}
		log.Debug().Str("flag_id", f.ID).Msg("deprovisioning flag")
		stats.FlagDeprovisioned.Inc()
		return cloudflareProvisioner.DeprovisionFlag(context.Background(), f)
	})

	consumer.HandleFunc("vex-provision-token", func(message *sarama.ConsumerMessage) error {
		log.Debug().Str("id", string(message.Value)).Msg("provisioning token")
		stats.ProjectProvisioned.Inc()
		return cloudflareProvisioner.ProvisionToken(context.Background(), &token.Token{ID: string(message.Value)})
	})
	consumer.HandleFunc("vex-deprovision-token", func(message *sarama.ConsumerMessage) error {
		log.Debug().Str("token_hash", hex.EncodeToString(message.Value)).Msg("deprovisioning token")
		stats.ProjectDeprovisioned.Inc()
		return cloudflareProvisioner.DeprovisionToken(context.Background(), &token.Token{TokenHash: message.Value})
	})

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), group, config)
	if err != nil {
		log.Panic().Err(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	log.Debug().Msg("starting consume loop")
	go func() {
		defer wg.Done()

		for {
			if err := client.Consume(ctx, strings.Split(topics, ","), consumer); err != nil {
				log.Panic().Err(err)
			}
			if ctx.Err() != nil {
				log.Error().Err(err).Msg("")
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready

	// start promhttp listener on metrics port
	m := chi.NewRouter()
	m.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%s", metricsPort), m); err != nil {
			log.Fatal().Err(err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	select {
	case s := <-sigs:
		log.Debug().Str("signal", s.String()).Msg("received signal")
	}
	cancel()
	wg.Wait()
	if err := client.Close(); err != nil {
		log.Debug().Err(err)
	}
}

func NewConsumer(skip bool) *Consumer {
	return &Consumer{
		ready:    make(chan bool),
		skip:     skip,
		handlers: make(map[string]MessageHandler),
	}
}

type MessageHandler func(message *sarama.ConsumerMessage) error

type Consumer struct {
	ready chan bool
	// whether in dev mode and shouldn't actually provision anything
	skip     bool
	handlers map[string]MessageHandler
}

func (c *Consumer) Handle(message *sarama.ConsumerMessage) error {
	for k, f := range c.handlers {
		if message.Topic == k {
			if c.skip {
				log.Debug().Str("topic", message.Topic).Str("value", string(message.Value)).Msg("provisioning skipped")
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
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}
			err := c.Handle(message)
			if err != nil {
				log.Debug().Err(err)
				continue
			}
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
