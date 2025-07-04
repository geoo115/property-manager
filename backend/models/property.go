package models

import (
	"time"

	"github.com/geoo115/property-manager/validator"
)

type Property struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	Name         string     `json:"name" gorm:"not null;index"`
	Description  string     `json:"description"`
	Bedrooms     uint       `json:"bedrooms" gorm:"not null"`
	Bathrooms    uint       `json:"bathrooms" gorm:"not null"`
	Price        float64    `json:"price" gorm:"not null"`
	SquareFeet   uint       `json:"square_feet"`
	Address      string     `json:"address" gorm:"not null"`
	City         string     `json:"city" gorm:"not null;index"`
	State        string     `json:"state" gorm:"index"`
	PostCode     string     `json:"post_code" gorm:"index"`
	Country      string     `json:"country" gorm:"default:'UK'"`
	PropertyType string     `json:"property_type" gorm:"default:'apartment'"`
	OwnerID      uint       `json:"owner_id" gorm:"not null;index"`
	Owner        User       `json:"owner" gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE;"`
	Available    bool       `json:"available" gorm:"default:true"`
	TenantID     *uint      `json:"tenant_id" gorm:"index"`
	Tenant       *User      `json:"tenant" gorm:"foreignKey:TenantID;constraint:OnDelete:SET NULL;"`
	Images       []string   `json:"images" gorm:"serializer:json"`
	Amenities    []string   `json:"amenities" gorm:"serializer:json"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at" gorm:"index"`

	// Relationships
	Units               []Unit        `json:"units,omitempty" gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE;"`
	Leases              []Lease       `json:"leases,omitempty" gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE;"`
	MaintenanceRequests []Maintenance `json:"maintenance_requests,omitempty" gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE;"`
	Invoices            []Invoice     `json:"invoices,omitempty" gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE;"`
	Expenses            []Expense     `json:"expenses,omitempty" gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE;"`
}

type Unit struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	PropertyID  uint       `json:"property_id" gorm:"not null;index"`
	Property    Property   `json:"property" gorm:"foreignKey:PropertyID;constraint:OnDelete:CASCADE;"`
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description"`
	Price       float64    `json:"price" gorm:"not null"`
	Available   bool       `json:"available" gorm:"default:true"`
	TenantID    *uint      `json:"tenant_id" gorm:"index"`
	Tenant      *User      `json:"tenant" gorm:"foreignKey:TenantID;constraint:OnDelete:SET NULL;"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`
}

// PropertyCreateRequest represents property creation request
type PropertyCreateRequest struct {
	Name         string   `json:"name" binding:"required"`
	Description  string   `json:"description"`
	Bedrooms     uint     `json:"bedrooms" binding:"required"`
	Bathrooms    uint     `json:"bathrooms" binding:"required"`
	Price        float64  `json:"price" binding:"required"`
	SquareFeet   uint     `json:"square_feet"`
	Address      string   `json:"address" binding:"required"`
	City         string   `json:"city" binding:"required"`
	State        string   `json:"state"`
	PostCode     string   `json:"post_code"`
	Country      string   `json:"country"`
	PropertyType string   `json:"property_type"`
	OwnerID      uint     `json:"owner_id" binding:"required"`
	Available    bool     `json:"available"`
	Images       []string `json:"images"`
	Amenities    []string `json:"amenities"`
}

// PropertyUpdateRequest represents property update request
type PropertyUpdateRequest struct {
	Name         *string  `json:"name"`
	Description  *string  `json:"description"`
	Bedrooms     *uint    `json:"bedrooms"`
	Bathrooms    *uint    `json:"bathrooms"`
	Price        *float64 `json:"price"`
	SquareFeet   *uint    `json:"square_feet"`
	Address      *string  `json:"address"`
	City         *string  `json:"city"`
	State        *string  `json:"state"`
	PostCode     *string  `json:"post_code"`
	Country      *string  `json:"country"`
	PropertyType *string  `json:"property_type"`
	Available    *bool    `json:"available"`
	Images       []string `json:"images"`
	Amenities    []string `json:"amenities"`
}

// PropertyResponse represents property response
type PropertyResponse struct {
	ID           uint          `json:"id"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Bedrooms     uint          `json:"bedrooms"`
	Bathrooms    uint          `json:"bathrooms"`
	Price        float64       `json:"price"`
	SquareFeet   uint          `json:"square_feet"`
	Address      string        `json:"address"`
	City         string        `json:"city"`
	State        string        `json:"state"`
	PostCode     string        `json:"post_code"`
	Country      string        `json:"country"`
	PropertyType string        `json:"property_type"`
	Available    bool          `json:"available"`
	Images       []string      `json:"images"`
	Amenities    []string      `json:"amenities"`
	Owner        UserResponse  `json:"owner"`
	Tenant       *UserResponse `json:"tenant"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// ToResponse converts Property to PropertyResponse
func (p *Property) ToResponse() PropertyResponse {
	response := PropertyResponse{
		ID:           p.ID,
		Name:         p.Name,
		Description:  p.Description,
		Bedrooms:     p.Bedrooms,
		Bathrooms:    p.Bathrooms,
		Price:        p.Price,
		SquareFeet:   p.SquareFeet,
		Address:      p.Address,
		City:         p.City,
		State:        p.State,
		PostCode:     p.PostCode,
		Country:      p.Country,
		PropertyType: p.PropertyType,
		Available:    p.Available,
		Images:       p.Images,
		Amenities:    p.Amenities,
		Owner:        p.Owner.ToResponse(),
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}

	if p.Tenant != nil {
		tenantResponse := p.Tenant.ToResponse()
		response.Tenant = &tenantResponse
	}

	return response
}

// Validate validates property creation request
func (req *PropertyCreateRequest) Validate() error {
	errors := validator.CollectValidationErrors(
		validator.ValidateRequired(req.Name, "name"),
		validator.ValidateMaxLength(req.Name, 200, "name"),
		validator.ValidateRequired(req.Address, "address"),
		validator.ValidateMaxLength(req.Address, 255, "address"),
		validator.ValidateRequired(req.City, "city"),
		validator.ValidateMaxLength(req.City, 100, "city"),
		validator.ValidatePositiveFloat(req.Price, "price"),
		validator.ValidatePositiveInt(int(req.Bedrooms), "bedrooms"),
		validator.ValidatePositiveInt(int(req.Bathrooms), "bathrooms"),
	)

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// Validate validates property update request
func (req *PropertyUpdateRequest) Validate() error {
	var errors validator.ValidationErrors

	if req.Name != nil {
		if err := validator.ValidateRequired(*req.Name, "name"); err != nil {
			errors = append(errors, *err)
		}
		if err := validator.ValidateMaxLength(*req.Name, 200, "name"); err != nil {
			errors = append(errors, *err)
		}
	}

	if req.Address != nil {
		if err := validator.ValidateRequired(*req.Address, "address"); err != nil {
			errors = append(errors, *err)
		}
		if err := validator.ValidateMaxLength(*req.Address, 255, "address"); err != nil {
			errors = append(errors, *err)
		}
	}

	if req.City != nil {
		if err := validator.ValidateRequired(*req.City, "city"); err != nil {
			errors = append(errors, *err)
		}
		if err := validator.ValidateMaxLength(*req.City, 100, "city"); err != nil {
			errors = append(errors, *err)
		}
	}

	if req.Price != nil {
		if err := validator.ValidatePositiveFloat(*req.Price, "price"); err != nil {
			errors = append(errors, *err)
		}
	}

	if req.Bedrooms != nil {
		if err := validator.ValidatePositiveInt(int(*req.Bedrooms), "bedrooms"); err != nil {
			errors = append(errors, *err)
		}
	}

	if req.Bathrooms != nil {
		if err := validator.ValidatePositiveInt(int(*req.Bathrooms), "bathrooms"); err != nil {
			errors = append(errors, *err)
		}
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// GetFullAddress returns the full address of the property
func (p *Property) GetFullAddress() string {
	return p.Address + ", " + p.City + ", " + p.PostCode
}

// IsOccupied checks if property is occupied
func (p *Property) IsOccupied() bool {
	return p.TenantID != nil
}

// TableName returns the table name for Property model
func (Property) TableName() string {
	return "properties"
}

// TableName returns the table name for Unit model
func (Unit) TableName() string {
	return "units"
}
