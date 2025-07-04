package models

import (
	"time"
)

type AuditLog struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	Action      string    `json:"action" gorm:"not null;index"`
	EntityType  string    `json:"entity_type" gorm:"not null;index"`
	EntityID    uint      `json:"entity_id" gorm:"not null;index"`
	OldData     string    `json:"old_data" gorm:"type:text"`
	NewData     string    `json:"new_data" gorm:"type:text"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

// AuditLogResponse represents audit log response
type AuditLogResponse struct {
	ID          uint         `json:"id"`
	Action      string       `json:"action"`
	EntityType  string       `json:"entity_type"`
	EntityID    uint         `json:"entity_id"`
	OldData     string       `json:"old_data"`
	NewData     string       `json:"new_data"`
	IPAddress   string       `json:"ip_address"`
	UserAgent   string       `json:"user_agent"`
	Description string       `json:"description"`
	User        UserResponse `json:"user"`
	CreatedAt   time.Time    `json:"created_at"`
}

// ToResponse converts AuditLog to AuditLogResponse
func (a *AuditLog) ToResponse() AuditLogResponse {
	return AuditLogResponse{
		ID:          a.ID,
		Action:      a.Action,
		EntityType:  a.EntityType,
		EntityID:    a.EntityID,
		OldData:     a.OldData,
		NewData:     a.NewData,
		IPAddress:   a.IPAddress,
		UserAgent:   a.UserAgent,
		Description: a.Description,
		User:        a.User.ToResponse(),
		CreatedAt:   a.CreatedAt,
	}
}

// TableName returns the table name for AuditLog model
func (AuditLog) TableName() string {
	return "audit_logs"
}
