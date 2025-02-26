// event/producer.go
package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/geoo115/property-manager/models"
	"github.com/segmentio/kafka-go"
)

// ProduceMaintenanceRequest sends a maintenance event to Kafka
func ProduceMaintenanceRequest(maintenance models.Maintenance) error {
	// Create Kafka writer
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKER")},
		Topic:    "maintenance-requests",
		Balancer: &kafka.LeastBytes{},
	})

	// Convert maintenance request to JSON
	data, err := json.Marshal(maintenance)
	if err != nil {
		log.Printf("❌ Error encoding maintenance request: %v", err)
		return err
	}

	// ✅ Log event before producing to Kafka
	fmt.Printf("📤 Sending Maintenance Request to Kafka: %s\n", string(data))

	// Send message to Kafka
	err = writer.WriteMessages(context.Background(), kafka.Message{
		Value: data,
	})
	if err != nil {
		log.Printf("❌ Failed to send Kafka message: %v", err)
		return err
	}

	fmt.Println("✅ Maintenance Request Sent Successfully to Kafka")
	return nil
}
