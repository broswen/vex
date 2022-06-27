package provisioner

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/broswen/vex/internal/project"
	"github.com/broswen/vex/internal/stats"
	"log"
	"time"
)

type KafkaProvisioner struct {
	provisionTopic   string
	deprovisionTopic string
	broker           string
	producer         sarama.SyncProducer
}

func NewKafkaProvisioner(provisionTopic, deprovisionTopic, broker string) (*KafkaProvisioner, error) {
	config := sarama.NewConfig()
	config.ClientID = "vex-config"
	version, err := sarama.ParseKafkaVersion("3.1.0")
	if err != nil {
		log.Fatal(err)
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
		provisionTopic:   provisionTopic,
		deprovisionTopic: deprovisionTopic,
		broker:           broker,
		producer:         producer,
	}, nil
}

func (p *KafkaProvisioner) ProvisionProject(ctx context.Context, pr *project.Project) error {
	_, _, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic:     p.provisionTopic,
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
		Topic:     p.deprovisionTopic,
		Value:     sarama.StringEncoder(pr.ID),
		Timestamp: time.Now(),
	})
	if err != nil {
		stats.DeprovisionError.Inc()
	}
	return err
}
