package models

import "time"

type Property struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Bedrooms    uint      `json:"bedrooms"`
	Bathrooms   uint      `json:"bathrooms"`
	Price       float64   `json:"price"`
	SquareFeet  uint      `json:"square_feet"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	PostCode    string    `json:"post_code"`
	OwnerID     uint      `json:"owner_id"` // Foreign key to User (landlord)
	Owner       User      `json:"owner" gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE;"`
	Available   bool      `json:"available"`
	TenantID    *uint     `json:"tenant_id"` // Tenant who rented the entire property (optional)
	Tenant      *User     `json:"tenant" gorm:"foreignKey:TenantID;constraint:OnDelete:SET NULL;"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Units []Unit `json:"units" gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE;"` // Units inside the property
}

type Unit struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	PropertyID  uint      `json:"property_id"` // Foreign key to Property
	Property    Property  `json:"property" gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE;"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Available   bool      `json:"available"`
	TenantID    *uint     `json:"tenant_id"` // Tenant renting this unit
	Tenant      *User     `json:"tenant" gorm:"foreignKey:TenantID;constraint:OnDelete:SET NULL;"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
