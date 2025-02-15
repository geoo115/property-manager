package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	FirstName string    `json:"first_name" required:"true"`
	LastName  string    `json:"last_name" required:"true"`
	Password  string    `json:"password" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Role      string    `json:"role" binding:"required,oneof=admin tenant landlord maintenanceTeam"`
	Phone     string    `json:"phone" gorm:"unique;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
