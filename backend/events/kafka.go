package events

import (
	"context"
	"log"

	"github.com/geoo115/property-manager/config"
	"github.com/segmentio/kafka-go"
)

var KafkaWriter *kafka.Writer

func InitKafka(cfg *config.Config) error {
	KafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP(cfg.Kafka.Broker),
		Topic:    cfg.Kafka.Topic,
		Balancer: &kafka.LeastBytes{},
	}
	return nil
}

func PublishEvent(message string) {
	if KafkaWriter == nil {
		log.Println("Kafka writer not initialized")
		return
	}

	err := KafkaWriter.WriteMessages(context.Background(),
		kafka.Message{Value: []byte(message)},
	)
	if err != nil {
		log.Printf("Failed to publish Kafka message: %v", err)
	}
}

func CloseKafka() {
	if KafkaWriter != nil {
		KafkaWriter.Close()
	}
}
