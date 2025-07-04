package models

import (
	"time"

	"github.com/geoo115/property-manager/validator"
)

type User struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Username  string     `json:"username" gorm:"unique;not null;index"`
	FirstName string     `json:"first_name" gorm:"not null"`
	LastName  string     `json:"last_name" gorm:"not null"`
	Password  string     `json:"-" gorm:"not null"` // Never expose password in JSON
	Email     string     `json:"email" gorm:"unique;not null;index"`
	Role      string     `json:"role" gorm:"not null;index;check:role IN ('admin','tenant','landlord','maintenanceTeam')"`
	Phone     string     `json:"phone" gorm:"unique;not null"`
	Avatar    string     `json:"avatar"`
	IsActive  bool       `json:"is_active" gorm:"default:true"`
	LastLogin *time.Time `json:"last_login"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`

	// Relationships
	OwnedProperties     []Property    `json:"owned_properties,omitempty" gorm:"foreignKey:OwnerID"`
	Leases              []Lease       `json:"leases,omitempty" gorm:"foreignKey:TenantID"`
	MaintenanceRequests []Maintenance `json:"maintenance_requests,omitempty" gorm:"foreignKey:RequestedByID"`
	Invoices            []Invoice     `json:"invoices,omitempty" gorm:"foreignKey:TenantID"`
	CreatedInvoices     []Invoice     `json:"created_invoices,omitempty" gorm:"foreignKey:CreatedByID"`
	Expenses            []Expense     `json:"expenses,omitempty" gorm:"foreignKey:CreatedByID"`
	AuditLogs           []AuditLog    `json:"audit_logs,omitempty" gorm:"foreignKey:UserID"`
}

// UserCreateRequest represents user creation request
type UserCreateRequest struct {
	Username  string `json:"username" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Role      string `json:"role" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Avatar    string `json:"avatar"`
}

// UserUpdateRequest represents user update request
type UserUpdateRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	Avatar    *string `json:"avatar"`
	IsActive  *bool   `json:"is_active"`
}

// UserResponse represents user response (without sensitive data)
type UserResponse struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	Phone     string     `json:"phone"`
	Avatar    string     `json:"avatar"`
	IsActive  bool       `json:"is_active"`
	LastLogin *time.Time `json:"last_login"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Role:      u.Role,
		Phone:     u.Phone,
		Avatar:    u.Avatar,
		IsActive:  u.IsActive,
		LastLogin: u.LastLogin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// Validate validates user creation request
func (req *UserCreateRequest) Validate() error {
	errors := validator.CollectValidationErrors(
		validator.ValidateRequired(req.Username, "username"),
		validator.ValidateMinLength(req.Username, 3, "username"),
		validator.ValidateMaxLength(req.Username, 50, "username"),
		validator.ValidateRequired(req.FirstName, "first_name"),
		validator.ValidateMaxLength(req.FirstName, 100, "first_name"),
		validator.ValidateRequired(req.LastName, "last_name"),
		validator.ValidateMaxLength(req.LastName, 100, "last_name"),
		validator.ValidatePassword(req.Password, "password"),
		validator.ValidateEmail(req.Email, "email"),
		validator.ValidatePhone(req.Phone, "phone"),
		validator.ValidateRole(req.Role, "role"),
	)

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// GetFullName returns the full name of the user
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// IsAdmin checks if user is admin
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// IsLandlord checks if user is landlord
func (u *User) IsLandlord() bool {
	return u.Role == "landlord"
}

// IsTenant checks if user is tenant
func (u *User) IsTenant() bool {
	return u.Role == "tenant"
}

// IsMaintenanceTeam checks if user is maintenance team member
func (u *User) IsMaintenanceTeam() bool {
	return u.Role == "maintenanceTeam"
}
