package models

import "time"

type Maintenance struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	TenantID    uint      `json:"tenant_id"`
	PropertyID  uint      `json:"property_id"`
	Description string    `json:"description"`
	RequestedAt time.Time `json:"requested_at"`
	Status      string    `json:"status"` // e.g., "pending", "completed"
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Tenant   User     `json:"tenant" gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE;"`
	Property Property `json:"property" gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE;"`
}
