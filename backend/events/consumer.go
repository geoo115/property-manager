package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/models"
	"github.com/segmentio/kafka-go"
)

// StartKafkaConsumer runs the Kafka consumer to process maintenance events
func StartKafkaConsumer() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_BROKER")},
		Topic:    "maintenance-requests",
		GroupID:  "maintenance-group",
		MaxBytes: 10e6,
	})

	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("❌ Failed to close Kafka reader: %v", err)
		}
	}()

	fmt.Println("🚀 Kafka Consumer started... Listening for maintenance requests.")
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("❌ Error reading Kafka message: %v", err)
			continue
		}

		var maintenance models.Maintenance
		if err := json.Unmarshal(msg.Value, &maintenance); err != nil {
			log.Printf("❌ Failed to parse maintenance request: %v", err)
			continue
		}

		if err := processMaintenanceEvent(maintenance); err != nil {
			log.Printf("❌ Failed to process maintenance event for ID %d: %v", maintenance.ID, err)
			// Optional: Add retry logic or dead-letter queue here
		} else {
			fmt.Printf("✅ Successfully processed maintenance request ID %d\n", maintenance.ID)
		}
	}
}

// processMaintenanceEvent handles the processing of a maintenance request event
func processMaintenanceEvent(maintenance models.Maintenance) error {
	// 1. Notify the maintenance team
	if err := notifyMaintenanceTeam(maintenance); err != nil {
		return fmt.Errorf("failed to notify maintenance team: %v", err)
	}

	// 2. Update status in the database
	if err := updateMaintenanceStatus(maintenance); err != nil {
		return fmt.Errorf("failed to update maintenance status: %v", err)
	}

	// 3. Log the event to an audit table
	if err := logToAudit(maintenance); err != nil {
		return fmt.Errorf("failed to log to audit: %v", err)
	}

	return nil
}

// notifyMaintenanceTeam sends an email notification to the maintenance team
func notifyMaintenanceTeam(maintenance models.Maintenance) error {
	// Email configuration from environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	toEmail := os.Getenv("MAINTENANCE_TEAM_EMAIL")

	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" || toEmail == "" {
		return fmt.Errorf("missing SMTP configuration")
	}

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	subject := fmt.Sprintf("New Maintenance Request #%d for Property #%d", maintenance.ID, maintenance.PropertyID)
	body := fmt.Sprintf(
		"A new maintenance request has been created:\n\n"+
			"ID: %d\n"+
			"Property ID: %d\n"+
			"Description: %s\n"+
			"Reported By: %d\n"+
			"Requested At: %s\n"+
			"Status: %s\n\n"+
			"Please review and take appropriate action.",
		maintenance.ID, maintenance.PropertyID, maintenance.Description, maintenance.ReporterID,
		maintenance.RequestedAt.Format(time.RFC3339), maintenance.Status,
	)

	msg := []byte(fmt.Sprintf(
		"To: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n"+
			"%s\r\n",
		toEmail, subject, body,
	))

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	err := smtp.SendMail(addr, auth, smtpUser, []string{toEmail}, msg)
	if err != nil {
		log.Printf("❌ Failed to send email notification: %v", err)
		return err
	}

	fmt.Printf("📧 Email notification sent to %s for maintenance request #%d\n", toEmail, maintenance.ID)
	return nil
}

// updateMaintenanceStatus updates the maintenance request status in the database
func updateMaintenanceStatus(maintenance models.Maintenance) error {
	// Example logic: Mark as "received" if new, escalate if urgent
	var newStatus string
	if maintenance.Status == "pending" {
		// Simple urgency check based on description (customize as needed)
		if containsUrgentKeywords(maintenance.Description) {
			newStatus = "urgent"
		} else {
			newStatus = "received"
		}
	} else {
		// No change if status isn’t "pending"
		return nil
	}

	result := db.DB.Model(&maintenance).Where("id = ?", maintenance.ID).Update("status", newStatus)
	if result.Error != nil {
		log.Printf("❌ Failed to update maintenance status: %v", result.Error)
		return result.Error
	}

	fmt.Printf("📝 Updated maintenance request #%d status to '%s'\n", maintenance.ID, newStatus)
	return nil
}

// containsUrgentKeywords checks for urgent keywords in the description
func containsUrgentKeywords(description string) bool {
	keywords := []string{"urgent", "emergency", "leak", "fire", "broken"}
	for _, kw := range keywords {
		if strings.Contains(strings.ToLower(description), kw) {
			return true
		}
	}
	return false
}

// logToAudit logs the maintenance event to an audit table
func logToAudit(maintenance models.Maintenance) error {
	audit := models.AuditLog{
		EventType:   "maintenance_request_created",
		EntityID:    maintenance.ID,
		Description: fmt.Sprintf("Maintenance request created for property %d: %s", maintenance.PropertyID, maintenance.Description),
		CreatedAt:   time.Now(),
	}

	if err := db.DB.Create(&audit).Error; err != nil {
		log.Printf("❌ Failed to log audit event: %v", err)
		return err
	}

	fmt.Printf("📝 Audit log created for maintenance request #%d\n", maintenance.ID)
	return nil
}
