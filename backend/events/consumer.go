package events

import (
	"context"
	"encoding/json"
	"fmt"
	"net/smtp"

	"github.com/geoo115/property-manager/config"
	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/logger"
	"github.com/geoo115/property-manager/models"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// StartKafkaConsumer runs the Kafka consumer to process maintenance events
func StartKafkaConsumer(cfg *config.Config) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.Kafka.Broker},
		Topic:    "maintenance-requests",
		GroupID:  "maintenance-group",
		MaxBytes: 10e6,
	})

	defer func() {
		if err := reader.Close(); err != nil {
			logger.LogError(err, "Failed to close Kafka reader", nil)
		}
	}()

	logger.LogInfo("Kafka Consumer started - Listening for maintenance requests", nil)

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			logger.LogError(err, "Error reading Kafka message", nil)
			continue
		}

		var maintenance models.Maintenance
		if err := json.Unmarshal(msg.Value, &maintenance); err != nil {
			logger.LogError(err, "Failed to parse maintenance request", logrus.Fields{
				"message": string(msg.Value),
			})
			continue
		}

		if err := processMaintenanceEvent(maintenance, cfg); err != nil {
			logger.LogError(err, "Failed to process maintenance event", logrus.Fields{
				"maintenance_id": maintenance.ID,
			})
		} else {
			logger.LogInfo("Successfully processed maintenance event", logrus.Fields{
				"maintenance_id": maintenance.ID,
			})
		}
	}
}

// processMaintenanceEvent processes a maintenance request event
func processMaintenanceEvent(maintenance models.Maintenance, cfg *config.Config) error {
	// Log the event in audit logs
	if err := logAuditEvent(maintenance); err != nil {
		logger.LogError(err, "Failed to log audit event", logrus.Fields{
			"maintenance_id": maintenance.ID,
		})
	}

	// Send notification to maintenance team
	if err := notifyMaintenanceTeam(maintenance, cfg); err != nil {
		logger.LogError(err, "Failed to notify maintenance team", logrus.Fields{
			"maintenance_id": maintenance.ID,
		})
	}

	return nil
}

// logAuditEvent logs the maintenance event to audit logs
func logAuditEvent(maintenance models.Maintenance) error {
	auditLog := models.AuditLog{
		UserID:      maintenance.RequestedByID,
		Action:      "CREATE",
		EntityType:  "maintenance",
		EntityID:    maintenance.ID,
		NewData:     fmt.Sprintf("Maintenance request created: %s", maintenance.Description),
		Description: fmt.Sprintf("New maintenance request created for property %d", maintenance.PropertyID),
	}

	if err := db.DB.Create(&auditLog).Error; err != nil {
		return fmt.Errorf("failed to log to audit: %v", err)
	}

	return nil
}

// notifyMaintenanceTeam sends an email notification to the maintenance team
func notifyMaintenanceTeam(maintenance models.Maintenance, cfg *config.Config) error {
	// Check if email configuration is available
	if cfg.Email.SMTPHost == "" || cfg.Email.SMTPUser == "" ||
		cfg.Email.SMTPPass == "" || cfg.Email.MaintenanceTeamEmail == "" {
		return fmt.Errorf("missing SMTP configuration")
	}

	auth := smtp.PlainAuth("", cfg.Email.SMTPUser, cfg.Email.SMTPPass, cfg.Email.SMTPHost)
	subject := fmt.Sprintf("New Maintenance Request #%d", maintenance.ID)
	body := fmt.Sprintf(
		"A new maintenance request has been created:\n\n"+
			"ID: %d\n"+
			"Property ID: %d\n"+
			"Title: %s\n"+
			"Description: %s\n"+
			"Priority: %s\n"+
			"Category: %s\n"+
			"Requested At: %s\n"+
			"Status: %s\n\n"+
			"Please review and take appropriate action.",
		maintenance.ID, maintenance.PropertyID, maintenance.Title, maintenance.Description,
		maintenance.Priority, maintenance.Category,
		maintenance.RequestedAt.Format("2006-01-02 15:04:05"), maintenance.Status,
	)

	message := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", cfg.Email.MaintenanceTeamEmail, subject, body)
	addr := fmt.Sprintf("%s:%d", cfg.Email.SMTPHost, cfg.Email.SMTPPort)

	err := smtp.SendMail(addr, auth, cfg.Email.SMTPUser, []string{cfg.Email.MaintenanceTeamEmail}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	logger.LogInfo("Email notification sent to maintenance team", logrus.Fields{
		"maintenance_id": maintenance.ID,
		"to_email":       cfg.Email.MaintenanceTeamEmail,
	})

	return nil
}
