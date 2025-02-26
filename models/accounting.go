// models/accounting.go
package models

import (
	"time"
)

// Invoice represents an invoice for a tenant (could be for rent, utilities, etc.)
type Invoice struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	TenantID          uint      `json:"tenant_id"`
	PropertyID        uint      `json:"property_id"`
	Amount            float64   `json:"amount"`
	PaidAmount        float64   `json:"paid_amount"`
	InvoiceDate       time.Time `json:"invoice_date"`
	Category          string    `json:"category"` // Example: "rent", "utilities", "late fee"
	DueDate           time.Time `json:"due_date"`
	PaymentStatus     string    `json:"payment_status"`     // e.g., "paid", "pending", "overdue"
	RefundedAmount    float64   `json:"refunded_amount"`    // Amount refunded to tenant
	RecurringInterval string    `json:"recurring_interval"` // e.g., "monthly", "yearly", ""
	Recurring         bool      `json:"recurring"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	Tenant   User     `json:"tenant" gorm:"foreignKey:TenantID"`
	Property Property `json:"property" gorm:"foreignKey:PropertyID"`
}

// Expense represents an expense related to a property or business
type Expense struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	PropertyID  uint      `json:"property_id"`
	Description string    `json:"description"`
	Category    string    `json:"category"` // Example: "maintenance", "utilities", "taxes"
	Amount      float64   `json:"amount"`
	ExpenseDate time.Time `json:"expense_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Property Property `json:"property" gorm:"foreignKey:PropertyID"`
}
