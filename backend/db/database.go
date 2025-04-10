package db

import (
	"fmt"

	"github.com/geoo115/property-manager/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Database connection successfully opened")

	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Property{})
	DB.AutoMigrate(&models.Lease{})
	DB.AutoMigrate(&models.Unit{})
	DB.AutoMigrate(&models.Maintenance{})
	DB.AutoMigrate(&models.Invoice{})
	DB.AutoMigrate(&models.Expense{})
	DB.AutoMigrate(&models.AuditLog{})
}
