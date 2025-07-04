// models/accounting.go
package models

import (
	"time"

	"github.com/geoo115/property-manager/validator"
)

// Invoice represents an invoice for a tenant (could be for rent, utilities, etc.)
type Invoice struct {
	ID                uint       `json:"id" gorm:"primaryKey"`
	TenantID          uint       `json:"tenant_id" gorm:"not null;index"`
	PropertyID        uint       `json:"property_id" gorm:"not null;index"`
	LeaseID           *uint      `json:"lease_id" gorm:"index"`
	CreatedByID       uint       `json:"created_by_id" gorm:"not null;index"`
	InvoiceNumber     string     `json:"invoice_number" gorm:"unique;not null"`
	Amount            float64    `json:"amount" gorm:"not null"`
	PaidAmount        float64    `json:"paid_amount" gorm:"default:0"`
	InvoiceDate       time.Time  `json:"invoice_date" gorm:"not null"`
	Category          string     `json:"category" gorm:"not null;check:category IN ('rent','utilities','late_fee','deposit','maintenance','other')"`
	DueDate           time.Time  `json:"due_date" gorm:"not null"`
	PaymentStatus     string     `json:"payment_status" gorm:"default:'pending';check:payment_status IN ('paid','pending','overdue','cancelled')"`
	RefundedAmount    float64    `json:"refunded_amount" gorm:"default:0"`
	RecurringInterval string     `json:"recurring_interval" gorm:"check:recurring_interval IN ('','monthly','quarterly','yearly')"`
	Recurring         bool       `json:"recurring" gorm:"default:false"`
	PaymentMethod     string     `json:"payment_method" gorm:"check:payment_method IN ('','cash','bank_transfer','card','cheque')"`
	Notes             string     `json:"notes" gorm:"type:text"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at" gorm:"index"`

	// Relationships
	Tenant    User     `json:"tenant" gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE;"`
	Property  Property `json:"property" gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE;"`
	Lease     *Lease   `json:"lease,omitempty" gorm:"foreignKey:LeaseID;constraint:OnDelete:SET NULL;"`
	CreatedBy User     `json:"created_by" gorm:"foreignKey:CreatedByID;constraint:OnDelete:CASCADE;"`
}

// Expense represents an expense related to a property or business
type Expense struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	PropertyID    uint       `json:"property_id" gorm:"not null;index"`
	CreatedByID   uint       `json:"created_by_id" gorm:"not null;index"`
	ExpenseNumber string     `json:"expense_number" gorm:"unique;not null"`
	Description   string     `json:"description" gorm:"not null"`
	Category      string     `json:"category" gorm:"not null;check:category IN ('maintenance','utilities','taxes','insurance','repairs','supplies','other')"`
	Amount        float64    `json:"amount" gorm:"not null"`
	ExpenseDate   time.Time  `json:"expense_date" gorm:"not null"`
	VendorName    string     `json:"vendor_name"`
	VendorEmail   string     `json:"vendor_email"`
	VendorPhone   string     `json:"vendor_phone"`
	PaymentMethod string     `json:"payment_method" gorm:"check:payment_method IN ('','cash','bank_transfer','card','cheque')"`
	ReceiptURL    string     `json:"receipt_url"`
	Notes         string     `json:"notes" gorm:"type:text"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at" gorm:"index"`

	// Relationships
	Property  Property `json:"property" gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE;"`
	CreatedBy User     `json:"created_by" gorm:"foreignKey:CreatedByID;constraint:OnDelete:CASCADE;"`
}

// InvoiceCreateRequest represents invoice creation request
type InvoiceCreateRequest struct {
	TenantID          uint      `json:"tenant_id" binding:"required"`
	PropertyID        uint      `json:"property_id" binding:"required"`
	LeaseID           *uint     `json:"lease_id"`
	Amount            float64   `json:"amount" binding:"required"`
	InvoiceDate       time.Time `json:"invoice_date" binding:"required"`
	Category          string    `json:"category" binding:"required"`
	DueDate           time.Time `json:"due_date" binding:"required"`
	RecurringInterval string    `json:"recurring_interval"`
	Recurring         bool      `json:"recurring"`
	PaymentMethod     string    `json:"payment_method"`
	Notes             string    `json:"notes"`
}

// InvoiceUpdateRequest represents invoice update request
type InvoiceUpdateRequest struct {
	Amount            *float64   `json:"amount"`
	PaidAmount        *float64   `json:"paid_amount"`
	InvoiceDate       *time.Time `json:"invoice_date"`
	Category          *string    `json:"category"`
	DueDate           *time.Time `json:"due_date"`
	PaymentStatus     *string    `json:"payment_status"`
	RefundedAmount    *float64   `json:"refunded_amount"`
	RecurringInterval *string    `json:"recurring_interval"`
	Recurring         *bool      `json:"recurring"`
	PaymentMethod     *string    `json:"payment_method"`
	Notes             *string    `json:"notes"`
}

// InvoiceResponse represents invoice response
type InvoiceResponse struct {
	ID                uint             `json:"id"`
	InvoiceNumber     string           `json:"invoice_number"`
	Amount            float64          `json:"amount"`
	PaidAmount        float64          `json:"paid_amount"`
	InvoiceDate       time.Time        `json:"invoice_date"`
	Category          string           `json:"category"`
	DueDate           time.Time        `json:"due_date"`
	PaymentStatus     string           `json:"payment_status"`
	RefundedAmount    float64          `json:"refunded_amount"`
	RecurringInterval string           `json:"recurring_interval"`
	Recurring         bool             `json:"recurring"`
	PaymentMethod     string           `json:"payment_method"`
	Notes             string           `json:"notes"`
	Tenant            UserResponse     `json:"tenant"`
	Property          PropertyResponse `json:"property"`
	Lease             *LeaseResponse   `json:"lease"`
	CreatedBy         UserResponse     `json:"created_by"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}

// ExpenseCreateRequest represents expense creation request
type ExpenseCreateRequest struct {
	PropertyID    uint      `json:"property_id" binding:"required"`
	Description   string    `json:"description" binding:"required"`
	Category      string    `json:"category" binding:"required"`
	Amount        float64   `json:"amount" binding:"required"`
	ExpenseDate   time.Time `json:"expense_date" binding:"required"`
	VendorName    string    `json:"vendor_name"`
	VendorEmail   string    `json:"vendor_email"`
	VendorPhone   string    `json:"vendor_phone"`
	PaymentMethod string    `json:"payment_method"`
	ReceiptURL    string    `json:"receipt_url"`
	Notes         string    `json:"notes"`
}

// ExpenseUpdateRequest represents expense update request
type ExpenseUpdateRequest struct {
	Description   *string    `json:"description"`
	Category      *string    `json:"category"`
	Amount        *float64   `json:"amount"`
	ExpenseDate   *time.Time `json:"expense_date"`
	VendorName    *string    `json:"vendor_name"`
	VendorEmail   *string    `json:"vendor_email"`
	VendorPhone   *string    `json:"vendor_phone"`
	PaymentMethod *string    `json:"payment_method"`
	ReceiptURL    *string    `json:"receipt_url"`
	Notes         *string    `json:"notes"`
}

// ExpenseResponse represents expense response
type ExpenseResponse struct {
	ID            uint             `json:"id"`
	ExpenseNumber string           `json:"expense_number"`
	Description   string           `json:"description"`
	Category      string           `json:"category"`
	Amount        float64          `json:"amount"`
	ExpenseDate   time.Time        `json:"expense_date"`
	VendorName    string           `json:"vendor_name"`
	VendorEmail   string           `json:"vendor_email"`
	VendorPhone   string           `json:"vendor_phone"`
	PaymentMethod string           `json:"payment_method"`
	ReceiptURL    string           `json:"receipt_url"`
	Notes         string           `json:"notes"`
	Property      PropertyResponse `json:"property"`
	CreatedBy     UserResponse     `json:"created_by"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

// ToResponse converts Invoice to InvoiceResponse
func (i *Invoice) ToResponse() InvoiceResponse {
	response := InvoiceResponse{
		ID:                i.ID,
		InvoiceNumber:     i.InvoiceNumber,
		Amount:            i.Amount,
		PaidAmount:        i.PaidAmount,
		InvoiceDate:       i.InvoiceDate,
		Category:          i.Category,
		DueDate:           i.DueDate,
		PaymentStatus:     i.PaymentStatus,
		RefundedAmount:    i.RefundedAmount,
		RecurringInterval: i.RecurringInterval,
		Recurring:         i.Recurring,
		PaymentMethod:     i.PaymentMethod,
		Notes:             i.Notes,
		Tenant:            i.Tenant.ToResponse(),
		Property:          i.Property.ToResponse(),
		CreatedBy:         i.CreatedBy.ToResponse(),
		CreatedAt:         i.CreatedAt,
		UpdatedAt:         i.UpdatedAt,
	}

	if i.Lease != nil {
		leaseResponse := i.Lease.ToResponse()
		response.Lease = &leaseResponse
	}

	return response
}

// ToResponse converts Expense to ExpenseResponse
func (e *Expense) ToResponse() ExpenseResponse {
	return ExpenseResponse{
		ID:            e.ID,
		ExpenseNumber: e.ExpenseNumber,
		Description:   e.Description,
		Category:      e.Category,
		Amount:        e.Amount,
		ExpenseDate:   e.ExpenseDate,
		VendorName:    e.VendorName,
		VendorEmail:   e.VendorEmail,
		VendorPhone:   e.VendorPhone,
		PaymentMethod: e.PaymentMethod,
		ReceiptURL:    e.ReceiptURL,
		Notes:         e.Notes,
		Property:      e.Property.ToResponse(),
		CreatedBy:     e.CreatedBy.ToResponse(),
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}

// Validate validates invoice creation request
func (req *InvoiceCreateRequest) Validate() error {
	errors := validator.CollectValidationErrors(
		validator.ValidatePositiveFloat(req.Amount, "amount"),
	)

	// Check if due date is after invoice date
	if !req.DueDate.After(req.InvoiceDate) {
		errors = append(errors, validator.ValidationError{
			Field:   "due_date",
			Message: "must be after invoice date",
			Value:   req.DueDate,
		})
	}

	// Validate category
	validCategories := []string{"rent", "utilities", "late_fee", "deposit", "maintenance", "other"}
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
			Message: "must be one of: rent, utilities, late_fee, deposit, maintenance, other",
			Value:   req.Category,
		})
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// Validate validates expense creation request
func (req *ExpenseCreateRequest) Validate() error {
	errors := validator.CollectValidationErrors(
		validator.ValidateRequired(req.Description, "description"),
		validator.ValidateMaxLength(req.Description, 500, "description"),
		validator.ValidatePositiveFloat(req.Amount, "amount"),
	)

	// Validate category
	validCategories := []string{"maintenance", "utilities", "taxes", "insurance", "repairs", "supplies", "other"}
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
			Message: "must be one of: maintenance, utilities, taxes, insurance, repairs, supplies, other",
			Value:   req.Category,
		})
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// IsOverdue checks if invoice is overdue
func (i *Invoice) IsOverdue() bool {
	return i.DueDate.Before(time.Now()) && i.PaymentStatus != "paid"
}

// IsPaid checks if invoice is fully paid
func (i *Invoice) IsPaid() bool {
	return i.PaymentStatus == "paid" || i.PaidAmount >= i.Amount
}

// BalanceRemaining returns the remaining balance
func (i *Invoice) BalanceRemaining() float64 {
	return i.Amount - i.PaidAmount
}

// DaysOverdue returns the number of days overdue
func (i *Invoice) DaysOverdue() int {
	if !i.IsOverdue() {
		return 0
	}
	duration := time.Since(i.DueDate)
	return int(duration.Hours() / 24)
}

// TableName returns the table name for Invoice model
func (Invoice) TableName() string {
	return "invoices"
}

// TableName returns the table name for Expense model
func (Expense) TableName() string {
	return "expenses"
}
