package models

import "time"

type AuditLog struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	EventType   string    `json:"event_type"`
	EntityID    uint      `json:"entity_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
