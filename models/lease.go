package models

import "time"

type Lease struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	TenantID        uint      `json:"tenant_id"`
	PropertyID      uint      `json:"property_id"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	MonthlyRent     float64   `json:"monthly_rent"`
	SecurityDeposit float64   `json:"security_deposit"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	Tenant   User     `json:"tenant" gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE;"`
	Property Property `json:"property" gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE;"`
}
