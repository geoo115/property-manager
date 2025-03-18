package events

import (
	"context"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

var KafkaWriter *kafka.Writer

func InitKafka() {
	KafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP(os.Getenv("KAFKA_BROKER")),
		Topic:    "property-events",
		Balancer: &kafka.LeastBytes{},
	}
}

func PublishEvent(message string) {
	err := KafkaWriter.WriteMessages(context.Background(),
		kafka.Message{Value: []byte(message)},
	)
	if err != nil {
		log.Printf("Failed to publish Kafka message: %v", err)
	}
}
