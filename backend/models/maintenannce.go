package models

import (
	"time"

	"github.com/geoo115/property-manager/validator"
)

type Maintenance struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	RequestedByID uint       `json:"requested_by_id" gorm:"not null;index"` // User who requested maintenance
	PropertyID    uint       `json:"property_id" gorm:"not null;index"`
	LeaseID       *uint      `json:"lease_id" gorm:"index"`       // Optional lease reference
	AssignedToID  *uint      `json:"assigned_to_id" gorm:"index"` // Maintenance team member
	Title         string     `json:"title" gorm:"not null"`
	Description   string     `json:"description" gorm:"type:text;not null"`
	Status        string     `json:"status" gorm:"default:'pending';check:status IN ('pending','in_progress','completed','cancelled')"`
	Priority      string     `json:"priority" gorm:"default:'medium';check:priority IN ('low','medium','high','urgent')"`
	Category      string     `json:"category" gorm:"default:'general';check:category IN ('plumbing','electrical','heating','appliances','general','emergency')"`
	EstimatedCost float64    `json:"estimated_cost" gorm:"default:0"`
	ActualCost    float64    `json:"actual_cost" gorm:"default:0"`
	RequestedAt   time.Time  `json:"requested_at"`
	ScheduledAt   *time.Time `json:"scheduled_at"`
	CompletedAt   *time.Time `json:"completed_at"`
	Notes         string     `json:"notes" gorm:"type:text"`
	Images        []string   `json:"images" gorm:"type:text"` // URLs to uploaded images
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at" gorm:"index"`

	// Relationships
	RequestedBy User     `json:"requested_by" gorm:"foreignKey:RequestedByID;constraint:OnDelete:CASCADE;"`
	Property    Property `json:"property" gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE;"`
	Lease       *Lease   `json:"lease,omitempty" gorm:"foreignKey:LeaseID;constraint:OnDelete:SET NULL;"`
	AssignedTo  *User    `json:"assigned_to,omitempty" gorm:"foreignKey:AssignedToID;constraint:OnDelete:SET NULL;"`
}

// MaintenanceCreateRequest represents maintenance creation request
type MaintenanceCreateRequest struct {
	PropertyID    uint       `json:"property_id" binding:"required"`
	LeaseID       *uint      `json:"lease_id"`
	Title         string     `json:"title" binding:"required"`
	Description   string     `json:"description" binding:"required"`
	Priority      string     `json:"priority"`
	Category      string     `json:"category"`
	EstimatedCost float64    `json:"estimated_cost"`
	ScheduledAt   *time.Time `json:"scheduled_at"`
	Images        []string   `json:"images"`
}

// MaintenanceUpdateRequest represents maintenance update request
type MaintenanceUpdateRequest struct {
	Title         *string    `json:"title"`
	Description   *string    `json:"description"`
	Status        *string    `json:"status"`
	Priority      *string    `json:"priority"`
	Category      *string    `json:"category"`
	EstimatedCost *float64   `json:"estimated_cost"`
	ActualCost    *float64   `json:"actual_cost"`
	ScheduledAt   *time.Time `json:"scheduled_at"`
	CompletedAt   *time.Time `json:"completed_at"`
	Notes         *string    `json:"notes"`
	Images        []string   `json:"images"`
	AssignedToID  *uint      `json:"assigned_to_id"`
}

// MaintenanceResponse represents maintenance response
type MaintenanceResponse struct {
	ID            uint             `json:"id"`
	Title         string           `json:"title"`
	Description   string           `json:"description"`
	Status        string           `json:"status"`
	Priority      string           `json:"priority"`
	Category      string           `json:"category"`
	EstimatedCost float64          `json:"estimated_cost"`
	ActualCost    float64          `json:"actual_cost"`
	RequestedAt   time.Time        `json:"requested_at"`
	ScheduledAt   *time.Time       `json:"scheduled_at"`
	CompletedAt   *time.Time       `json:"completed_at"`
	Notes         string           `json:"notes"`
	Images        []string         `json:"images"`
	RequestedBy   UserResponse     `json:"requested_by"`
	Property      PropertyResponse `json:"property"`
	Lease         *LeaseResponse   `json:"lease"`
	AssignedTo    *UserResponse    `json:"assigned_to"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

// ToResponse converts Maintenance to MaintenanceResponse
func (m *Maintenance) ToResponse() MaintenanceResponse {
	response := MaintenanceResponse{
		ID:            m.ID,
		Title:         m.Title,
		Description:   m.Description,
		Status:        m.Status,
		Priority:      m.Priority,
		Category:      m.Category,
		EstimatedCost: m.EstimatedCost,
		ActualCost:    m.ActualCost,
		RequestedAt:   m.RequestedAt,
		ScheduledAt:   m.ScheduledAt,
		CompletedAt:   m.CompletedAt,
		Notes:         m.Notes,
		Images:        m.Images,
		RequestedBy:   m.RequestedBy.ToResponse(),
		Property:      m.Property.ToResponse(),
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}

	if m.Lease != nil {
		leaseResponse := m.Lease.ToResponse()
		response.Lease = &leaseResponse
	}

	if m.AssignedTo != nil {
		assignedToResponse := m.AssignedTo.ToResponse()
		response.AssignedTo = &assignedToResponse
	}

	return response
}

// Validate validates maintenance creation request
func (req *MaintenanceCreateRequest) Validate() error {
	errors := validator.CollectValidationErrors(
		validator.ValidateRequired(req.Title, "title"),
		validator.ValidateMaxLength(req.Title, 200, "title"),
		validator.ValidateRequired(req.Description, "description"),
		validator.ValidateMaxLength(req.Description, 2000, "description"),
		validator.ValidateNonNegativeFloat(req.EstimatedCost, "estimated_cost"),
	)

	// Validate priority
	if req.Priority != "" {
		validPriorities := []string{"low", "medium", "high", "urgent"}
		valid := false
		for _, p := range validPriorities {
			if req.Priority == p {
				valid = true
				break
			}
		}
		if !valid {
			errors = append(errors, validator.ValidationError{
				Field:   "priority",
				Message: "must be one of: low, medium, high, urgent",
				Value:   req.Priority,
			})
		}
	}

	// Validate category
	if req.Category != "" {
		validCategories := []string{"plumbing", "electrical", "heating", "appliances", "general", "emergency"}
		valid := false
		for _, c := range validCategories {
			if req.Category == c {
				valid = true
				break
			}
		}
		if !valid {
			errors = append(errors, validator.ValidationError{
				Field:   "category",
				Message: "must be one of: plumbing, electrical, heating, appliances, general, emergency",
				Value:   req.Category,
			})
		}
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// IsCompleted checks if maintenance is completed
func (m *Maintenance) IsCompleted() bool {
	return m.Status == "completed"
}

// IsOverdue checks if maintenance is overdue
func (m *Maintenance) IsOverdue() bool {
	if m.ScheduledAt == nil {
		return false
	}
	return m.ScheduledAt.Before(time.Now()) && !m.IsCompleted()
}

// IsUrgent checks if maintenance is urgent priority
func (m *Maintenance) IsUrgent() bool {
	return m.Priority == "urgent"
}

// DaysOverdue returns the number of days overdue
func (m *Maintenance) DaysOverdue() int {
	if !m.IsOverdue() {
		return 0
	}
	duration := time.Since(*m.ScheduledAt)
	return int(duration.Hours() / 24)
}

// TableName returns the table name for Maintenance model
func (Maintenance) TableName() string {
	return "maintenance_requests"
}
