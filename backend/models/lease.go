package models

import (
	"time"

	"github.com/geoo115/property-manager/validator"
)

type Lease struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	TenantID        uint       `json:"tenant_id" gorm:"not null;index"`
	PropertyID      uint       `json:"property_id" gorm:"not null;index"`
	StartDate       time.Time  `json:"start_date" gorm:"type:timestamp"`
	EndDate         time.Time  `json:"end_date" gorm:"type:timestamp"`
	MonthlyRent     float64    `json:"monthly_rent" gorm:"not null"`
	SecurityDeposit float64    `json:"security_deposit" gorm:"not null"`
	Status          string     `json:"status" gorm:"default:'active';check:status IN ('active','expired','terminated','pending')"`
	LeaseType       string     `json:"lease_type" gorm:"default:'fixed';check:lease_type IN ('fixed','periodic','short_term')"`
	RenewalTerms    string     `json:"renewal_terms" gorm:"type:text"`
	SpecialTerms    string     `json:"special_terms" gorm:"type:text"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at" gorm:"index"`

	// Relationships
	Tenant              User          `json:"tenant" gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE;"`
	Property            Property      `json:"property" gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE;"`
	MaintenanceRequests []Maintenance `json:"maintenance_requests,omitempty" gorm:"foreignKey:LeaseID;constraint:OnDelete:CASCADE;"`
	Invoices            []Invoice     `json:"invoices,omitempty" gorm:"foreignKey:LeaseID;constraint:OnDelete:CASCADE;"`
}

// LeaseCreateRequest represents lease creation request
type LeaseCreateRequest struct {
	TenantID        uint      `json:"tenant_id" binding:"required"`
	PropertyID      uint      `json:"property_id" binding:"required"`
	StartDate       time.Time `json:"start_date" binding:"required"`
	EndDate         time.Time `json:"end_date" binding:"required"`
	MonthlyRent     float64   `json:"monthly_rent" binding:"required"`
	SecurityDeposit float64   `json:"security_deposit" binding:"required"`
	LeaseType       string    `json:"lease_type"`
	RenewalTerms    string    `json:"renewal_terms"`
	SpecialTerms    string    `json:"special_terms"`
}

// LeaseUpdateRequest represents lease update request
type LeaseUpdateRequest struct {
	StartDate       *time.Time `json:"start_date"`
	EndDate         *time.Time `json:"end_date"`
	MonthlyRent     *float64   `json:"monthly_rent"`
	SecurityDeposit *float64   `json:"security_deposit"`
	Status          *string    `json:"status"`
	LeaseType       *string    `json:"lease_type"`
	RenewalTerms    *string    `json:"renewal_terms"`
	SpecialTerms    *string    `json:"special_terms"`
}

// LeaseResponse represents lease response
type LeaseResponse struct {
	ID              uint             `json:"id"`
	StartDate       time.Time        `json:"start_date"`
	EndDate         time.Time        `json:"end_date"`
	MonthlyRent     float64          `json:"monthly_rent"`
	SecurityDeposit float64          `json:"security_deposit"`
	Status          string           `json:"status"`
	LeaseType       string           `json:"lease_type"`
	RenewalTerms    string           `json:"renewal_terms"`
	SpecialTerms    string           `json:"special_terms"`
	Tenant          UserResponse     `json:"tenant"`
	Property        PropertyResponse `json:"property"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

// ToResponse converts Lease to LeaseResponse
func (l *Lease) ToResponse() LeaseResponse {
	return LeaseResponse{
		ID:              l.ID,
		StartDate:       l.StartDate,
		EndDate:         l.EndDate,
		MonthlyRent:     l.MonthlyRent,
		SecurityDeposit: l.SecurityDeposit,
		Status:          l.Status,
		LeaseType:       l.LeaseType,
		RenewalTerms:    l.RenewalTerms,
		SpecialTerms:    l.SpecialTerms,
		Tenant:          l.Tenant.ToResponse(),
		Property:        l.Property.ToResponse(),
		CreatedAt:       l.CreatedAt,
		UpdatedAt:       l.UpdatedAt,
	}
}

// Validate validates lease creation request
func (req *LeaseCreateRequest) Validate() error {
	errors := validator.CollectValidationErrors(
		validator.ValidatePositiveFloat(req.MonthlyRent, "monthly_rent"),
		validator.ValidateNonNegativeFloat(req.SecurityDeposit, "security_deposit"),
	)

	// Check if end date is after start date
	if !req.EndDate.After(req.StartDate) {
		errors = append(errors, validator.ValidationError{
			Field:   "end_date",
			Message: "must be after start date",
			Value:   req.EndDate,
		})
	}

	// Check if start date is not in the past
	if req.StartDate.Before(time.Now().Truncate(24 * time.Hour)) {
		errors = append(errors, validator.ValidationError{
			Field:   "start_date",
			Message: "cannot be in the past",
			Value:   req.StartDate,
		})
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// IsActive checks if lease is currently active
func (l *Lease) IsActive() bool {
	now := time.Now()
	return l.Status == "active" && l.StartDate.Before(now) && l.EndDate.After(now)
}

// IsExpired checks if lease has expired
func (l *Lease) IsExpired() bool {
	return l.EndDate.Before(time.Now()) || l.Status == "expired"
}

// DaysRemaining returns the number of days remaining in the lease
func (l *Lease) DaysRemaining() int {
	if l.IsExpired() {
		return 0
	}
	duration := l.EndDate.Sub(time.Now())
	return int(duration.Hours() / 24)
}

// TotalRent calculates total rent for the lease period
func (l *Lease) TotalRent() float64 {
	months := l.EndDate.Sub(l.StartDate).Hours() / (24 * 30)
	return l.MonthlyRent * months
}

// TableName returns the table name for Lease model
func (Lease) TableName() string {
	return "leases"
}
