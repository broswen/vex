package provisioner

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/broswen/vex/internal/flag"
	"github.com/broswen/vex/internal/project"
	"github.com/broswen/vex/internal/stats"
	"github.com/broswen/vex/internal/token"
	"github.com/rs/zerolog/log"
	"time"
)

type KafkaProvisioner struct {
	provisionProjectTopic   string
	deprovisionProjectTopic string
	provisionFlagTopic      string
	deprovisionFlagTopic    string
	provisionTokenTopic     string
	deprovisionTokenTopic   string
	broker                  string
	producer                sarama.SyncProducer
}

func NewKafkaProvisioner(provisionProjectTopic, deprovisionProjectTopic, provisionTokenTopic, deprovisionTokenTopic string, broker string) (*KafkaProvisioner, error) {
	config := sarama.NewConfig()
	config.ClientID = "vex-config"
	version, err := sarama.ParseKafkaVersion("3.1.0")
	if err != nil {
		log.Fatal().Err(err)
	}
	config.Version = version
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	producer, err := sarama.NewSyncProducer([]string{broker}, config)
	if err != nil {
		return nil, err
	}
	return &KafkaProvisioner{
		provisionProjectTopic:   provisionProjectTopic,
		deprovisionProjectTopic: deprovisionProjectTopic,
		provisionTokenTopic:     provisionTokenTopic,
		deprovisionTokenTopic:   deprovisionTokenTopic,
		broker:                  broker,
		producer:                producer,
	}, nil
}

func (p *KafkaProvisioner) ProvisionFlag(ctx context.Context, f *flag.Flag) error {
	b, err := json.Marshal(f)
	if err != nil {
		return err
	}
	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic:     p.provisionFlagTopic,
		Key:       sarama.StringEncoder(f.ID),
		Value:     sarama.StringEncoder(b),
		Timestamp: time.Now(),
	})
	if err != nil {
		stats.ProvisionError.Inc()
	}
	return err
}

func (p *KafkaProvisioner) DeprovisionFlag(ctx context.Context, f *flag.Flag) error {
	b, err := json.Marshal(f)
	if err != nil {
		return err
	}
	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic:     p.deprovisionProjectTopic,
		Key:       sarama.StringEncoder(b),
		Value:     sarama.StringEncoder(b),
		Timestamp: time.Now(),
	})
	if err != nil {
		stats.DeprovisionError.Inc()
	}
	return err
}

func (p *KafkaProvisioner) ProvisionProject(ctx context.Context, pr *project.Project) error {
	_, _, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic:     p.provisionProjectTopic,
		Key:       sarama.StringEncoder(pr.ID),
		Value:     sarama.StringEncoder(pr.ID),
		Timestamp: time.Now(),
	})
	if err != nil {
		stats.ProvisionError.Inc()
	}
	return err
}

func (p *KafkaProvisioner) DeprovisionProject(ctx context.Context, pr *project.Project) error {
	_, _, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic:     p.deprovisionProjectTopic,
		Key:       sarama.StringEncoder(pr.ID),
		Value:     sarama.StringEncoder(pr.ID),
		Timestamp: time.Now(),
	})
	if err != nil {
		stats.DeprovisionError.Inc()
	}
	return err
}

func (p *KafkaProvisioner) ProvisionToken(ctx context.Context, t *token.Token) error {
	_, _, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic:     p.provisionTokenTopic,
		Key:       sarama.StringEncoder(t.ID),
		Value:     sarama.StringEncoder(t.ID),
		Timestamp: time.Now(),
	})
	if err != nil {
		stats.ProvisionError.Inc()
	}
	return err
}

func (p *KafkaProvisioner) DeprovisionToken(ctx context.Context, t *token.Token) error {
	_, _, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic:     p.deprovisionTokenTopic,
		Key:       sarama.StringEncoder(t.Token),
		Value:     sarama.ByteEncoder(t.TokenHash),
		Timestamp: time.Now(),
	})
	if err != nil {
		stats.DeprovisionError.Inc()
	}
	return err
}
