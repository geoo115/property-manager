package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/geoo115/property-manager/config"
	"github.com/geoo115/property-manager/db"
	"github.com/geoo115/property-manager/logger"
	"github.com/geoo115/property-manager/models"
	"github.com/geoo115/property-manager/utils"
)

func main() {
	// Initialize logger
	logger.InitLogger()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	if err := db.Init(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Seed all data
	seedUsers()
	seedProperties()
	seedLeases()
	seedMaintenanceRequests()
	seedInvoices()
	seedExpenses()

	fmt.Println("All seed data creation completed!")
}

func seedUsers() {
	fmt.Println("Creating demo users...")

	demoUsers := []models.User{
		{
			Username:  "admin",
			Email:     "admin@example.com",
			FirstName: "Admin",
			LastName:  "User",
			Role:      "admin",
			Password:  "Admin123!",
			Phone:     "1234567890",
			IsActive:  true,
		},
		{
			Username:  "landlord1",
			Email:     "landlord1@example.com",
			FirstName: "John",
			LastName:  "Smith",
			Role:      "landlord",
			Password:  "Landlord123!",
			Phone:     "1234567891",
			IsActive:  true,
		},
		{
			Username:  "landlord2",
			Email:     "landlord2@example.com",
			FirstName: "Sarah",
			LastName:  "Johnson",
			Role:      "landlord",
			Password:  "Landlord123!",
			Phone:     "1234567892",
			IsActive:  true,
		},
		{
			Username:  "tenant1",
			Email:     "tenant1@example.com",
			FirstName: "Alice",
			LastName:  "Brown",
			Role:      "tenant",
			Password:  "Tenant123!",
			Phone:     "1234567893",
			IsActive:  true,
		},
		{
			Username:  "tenant2",
			Email:     "tenant2@example.com",
			FirstName: "Bob",
			LastName:  "Wilson",
			Role:      "tenant",
			Password:  "Tenant123!",
			Phone:     "1234567894",
			IsActive:  true,
		},
		{
			Username:  "tenant3",
			Email:     "tenant3@example.com",
			FirstName: "Charlie",
			LastName:  "Davis",
			Role:      "tenant",
			Password:  "Tenant123!",
			Phone:     "1234567895",
			IsActive:  true,
		},
		{
			Username:  "maintenance1",
			Email:     "maintenance1@example.com",
			FirstName: "Mike",
			LastName:  "Repair",
			Role:      "maintenanceTeam",
			Password:  "Maintenance123!",
			Phone:     "1234567896",
			IsActive:  true,
		},
		{
			Username:  "maintenance2",
			Email:     "maintenance2@example.com",
			FirstName: "Tony",
			LastName:  "Fix",
			Role:      "maintenanceTeam",
			Password:  "Maintenance123!",
			Phone:     "1234567897",
			IsActive:  true,
		},
	}

	for _, user := range demoUsers {
		// Hash password
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			log.Printf("Failed to hash password for %s: %v", user.Username, err)
			continue
		}
		user.Password = hashedPassword

		// Check if user already exists
		var existingUser models.User
		if err := db.DB.Where("email = ? OR username = ?", user.Email, user.Username).First(&existingUser).Error; err == nil {
			fmt.Printf("User %s already exists, skipping...\n", user.Username)
			continue
		}

		// Create user
		if err := db.DB.Create(&user).Error; err != nil {
			log.Printf("Failed to create user %s: %v", user.Username, err)
			continue
		}

		fmt.Printf("Created demo user: %s (%s)\n", user.Username, user.Email)
	}
}

func seedProperties() {
	fmt.Println("Creating demo properties...")

	// Get landlords
	var landlords []models.User
	db.DB.Where("role = ?", "landlord").Find(&landlords)
	if len(landlords) == 0 {
		log.Println("No landlords found, skipping property creation")
		return
	}

	demoProperties := []models.Property{
		{
			Name:         "Sunset Apartments Unit 1A",
			Description:  "Modern 2-bedroom apartment with city views and balcony",
			Bedrooms:     2,
			Bathrooms:    2,
			Price:        1200.00,
			SquareFeet:   850,
			Address:      "123 Main Street, Unit 1A",
			City:         "London",
			State:        "England",
			PostCode:     "SW1A 1AA",
			Country:      "UK",
			PropertyType: "apartment",
			OwnerID:      landlords[0].ID,
			Available:    false, // Will be leased
			Images:       []string{"apartment1_1.jpg", "apartment1_2.jpg"},
			Amenities:    []string{"Balcony", "Parking", "Gym", "Elevator"},
		},
		{
			Name:         "Garden View House",
			Description:  "Spacious 3-bedroom house with garden and garage",
			Bedrooms:     3,
			Bathrooms:    2,
			Price:        1800.00,
			SquareFeet:   1200,
			Address:      "456 Oak Avenue",
			City:         "Manchester",
			State:        "England",
			PostCode:     "M1 1AA",
			Country:      "UK",
			PropertyType: "house",
			OwnerID:      landlords[0].ID,
			Available:    false, // Will be leased
			Images:       []string{"house1_1.jpg", "house1_2.jpg", "house1_3.jpg"},
			Amenities:    []string{"Garden", "Garage", "Fireplace"},
		},
		{
			Name:         "City Center Studio",
			Description:  "Compact studio apartment in the heart of the city",
			Bedrooms:     1,
			Bathrooms:    1,
			Price:        800.00,
			SquareFeet:   400,
			Address:      "789 High Street, Unit 5B",
			City:         "Birmingham",
			State:        "England",
			PostCode:     "B1 1AA",
			Country:      "UK",
			PropertyType: "studio",
			OwnerID:      landlords[1].ID,
			Available:    false, // Will be leased
			Images:       []string{"studio1_1.jpg", "studio1_2.jpg"},
			Amenities:    []string{"24/7 Security", "Concierge", "Rooftop Terrace"},
		},
		{
			Name:         "Riverside Penthouse",
			Description:  "Luxury penthouse with river views and premium amenities",
			Bedrooms:     4,
			Bathrooms:    3,
			Price:        3500.00,
			SquareFeet:   2000,
			Address:      "101 River Walk, Penthouse",
			City:         "London",
			State:        "England",
			PostCode:     "SE1 9AA",
			Country:      "UK",
			PropertyType: "penthouse",
			OwnerID:      landlords[1].ID,
			Available:    true, // Available for lease
			Images:       []string{"penthouse1_1.jpg", "penthouse1_2.jpg", "penthouse1_3.jpg"},
			Amenities:    []string{"River View", "Private Terrace", "Concierge", "Parking", "Gym"},
		},
		{
			Name:         "Cozy Cottage",
			Description:  "Charming 2-bedroom cottage with private garden",
			Bedrooms:     2,
			Bathrooms:    1,
			Price:        1100.00,
			SquareFeet:   700,
			Address:      "25 Cottage Lane",
			City:         "Brighton",
			State:        "England",
			PostCode:     "BN1 1AA",
			Country:      "UK",
			PropertyType: "cottage",
			OwnerID:      landlords[0].ID,
			Available:    true, // Available for lease
			Images:       []string{"cottage1_1.jpg", "cottage1_2.jpg"},
			Amenities:    []string{"Private Garden", "Fireplace", "Parking"},
		},
	}

	for _, property := range demoProperties {
		// Check if property already exists
		var existingProperty models.Property
		if err := db.DB.Where("address = ?", property.Address).First(&existingProperty).Error; err == nil {
			fmt.Printf("Property at %s already exists, skipping...\n", property.Address)
			continue
		}

		// Create property
		if err := db.DB.Create(&property).Error; err != nil {
			log.Printf("Failed to create property %s: %v", property.Name, err)
			continue
		}

		fmt.Printf("Created demo property: %s\n", property.Name)
	}
}

func seedLeases() {
	fmt.Println("Creating demo leases...")

	// Get tenants and properties
	var tenants []models.User
	var properties []models.Property
	db.DB.Where("role = ?", "tenant").Find(&tenants)
	db.DB.Where("available = ?", false).Find(&properties)

	if len(tenants) == 0 || len(properties) == 0 {
		log.Println("Not enough tenants or properties found, skipping lease creation")
		return
	}

	now := time.Now()

	demoLeases := []models.Lease{
		{
			TenantID:        tenants[0].ID,
			PropertyID:      properties[0].ID,
			StartDate:       now.AddDate(0, -6, 0), // 6 months ago
			EndDate:         now.AddDate(1, 0, 0),  // 1 year from now
			MonthlyRent:     1200.00,
			SecurityDeposit: 2400.00,
			Status:          "active",
			LeaseType:       "fixed",
			RenewalTerms:    "Automatic renewal for 12 months unless 30 days notice given",
			SpecialTerms:    "Pet allowed with additional deposit",
		},
		{
			TenantID:        tenants[1].ID,
			PropertyID:      properties[1].ID,
			StartDate:       now.AddDate(0, -3, 0), // 3 months ago
			EndDate:         now.AddDate(0, 9, 0),  // 9 months from now
			MonthlyRent:     1800.00,
			SecurityDeposit: 3600.00,
			Status:          "active",
			LeaseType:       "fixed",
			RenewalTerms:    "Standard renewal terms apply",
			SpecialTerms:    "Garden maintenance included",
		},
		{
			TenantID:        tenants[2].ID,
			PropertyID:      properties[2].ID,
			StartDate:       now.AddDate(0, -1, 0), // 1 month ago
			EndDate:         now.AddDate(0, 11, 0), // 11 months from now
			MonthlyRent:     800.00,
			SecurityDeposit: 1600.00,
			Status:          "active",
			LeaseType:       "fixed",
			RenewalTerms:    "Month-to-month after initial term",
			SpecialTerms:    "No smoking policy strictly enforced",
		},
	}

	for i, lease := range demoLeases {
		// Check if lease already exists
		var existingLease models.Lease
		if err := db.DB.Where("tenant_id = ? AND property_id = ?", lease.TenantID, lease.PropertyID).First(&existingLease).Error; err == nil {
			fmt.Printf("Lease for tenant %d and property %d already exists, skipping...\n", lease.TenantID, lease.PropertyID)
			continue
		}

		// Create lease
		if err := db.DB.Create(&lease).Error; err != nil {
			log.Printf("Failed to create lease %d: %v", i+1, err)
			continue
		}

		// Update property tenant_id
		db.DB.Model(&properties[i]).Update("tenant_id", lease.TenantID)

		fmt.Printf("Created demo lease: Tenant %d -> Property %d\n", lease.TenantID, lease.PropertyID)
	}
}

func seedMaintenanceRequests() {
	fmt.Println("Creating demo maintenance requests...")

	// Get users and properties
	var tenants []models.User
	var maintenanceTeam []models.User
	var properties []models.Property
	var leases []models.Lease

	db.DB.Where("role = ?", "tenant").Find(&tenants)
	db.DB.Where("role = ?", "maintenanceTeam").Find(&maintenanceTeam)
	db.DB.Find(&properties)
	db.DB.Find(&leases)

	if len(tenants) == 0 || len(properties) == 0 {
		log.Println("Not enough users or properties found, skipping maintenance request creation")
		return
	}

	now := time.Now()
	demoMaintenance := []models.Maintenance{
		{
			RequestedByID: tenants[0].ID,
			PropertyID:    properties[0].ID,
			LeaseID: func() *uint {
				if len(leases) > 0 {
					return &leases[0].ID
				} else {
					return nil
				}
			}(),
			AssignedToID: func() *uint {
				if len(maintenanceTeam) > 0 {
					return &maintenanceTeam[0].ID
				} else {
					return nil
				}
			}(),
			Title:         "Leaky Kitchen Faucet",
			Description:   "The kitchen faucet is dripping constantly and needs repair or replacement",
			Status:        "completed",
			Priority:      "medium",
			Category:      "plumbing",
			EstimatedCost: 150.00,
			ActualCost:    120.00,
			RequestedAt:   now.AddDate(0, 0, -15), // 15 days ago
			ScheduledAt:   func() *time.Time { t := now.AddDate(0, 0, -10); return &t }(),
			CompletedAt:   func() *time.Time { t := now.AddDate(0, 0, -8); return &t }(),
			Notes:         "Replaced faucet cartridge and checked all connections. No further issues expected.",
		},
		{
			RequestedByID: tenants[1].ID,
			PropertyID:    properties[1].ID,
			LeaseID: func() *uint {
				if len(leases) > 1 {
					return &leases[1].ID
				} else {
					return nil
				}
			}(),
			AssignedToID: func() *uint {
				if len(maintenanceTeam) > 1 {
					return &maintenanceTeam[1].ID
				} else if len(maintenanceTeam) > 0 {
					return &maintenanceTeam[0].ID
				} else {
					return nil
				}
			}(),
			Title:         "Heating System Not Working",
			Description:   "Central heating system not responding, house is getting cold",
			Status:        "in_progress",
			Priority:      "high",
			Category:      "heating",
			EstimatedCost: 300.00,
			ActualCost:    0.00,
			RequestedAt:   now.AddDate(0, 0, -5), // 5 days ago
			ScheduledAt:   func() *time.Time { t := now.AddDate(0, 0, -2); return &t }(),
			Notes:         "Initial diagnosis shows faulty thermostat. Replacement parts ordered.",
		},
		{
			RequestedByID: tenants[2].ID,
			PropertyID:    properties[2].ID,
			LeaseID: func() *uint {
				if len(leases) > 2 {
					return &leases[2].ID
				} else {
					return nil
				}
			}(),
			Title:         "Electrical Outlet Not Working",
			Description:   "Main bedroom electrical outlet has stopped working completely",
			Status:        "pending",
			Priority:      "medium",
			Category:      "electrical",
			EstimatedCost: 100.00,
			ActualCost:    0.00,
			RequestedAt:   now.AddDate(0, 0, -2), // 2 days ago
			Notes:         "Awaiting electrician availability",
		},
		{
			RequestedByID: tenants[0].ID,
			PropertyID:    properties[0].ID,
			LeaseID: func() *uint {
				if len(leases) > 0 {
					return &leases[0].ID
				} else {
					return nil
				}
			}(),
			Title:         "Washing Machine Making Noise",
			Description:   "Washing machine making loud banging noise during spin cycle",
			Status:        "pending",
			Priority:      "low",
			Category:      "appliances",
			EstimatedCost: 80.00,
			ActualCost:    0.00,
			RequestedAt:   now.AddDate(0, 0, -1), // 1 day ago
			Notes:         "Request submitted, pending initial assessment",
		},
		{
			RequestedByID: tenants[1].ID,
			PropertyID:    properties[1].ID,
			LeaseID: func() *uint {
				if len(leases) > 1 {
					return &leases[1].ID
				} else {
					return nil
				}
			}(),
			AssignedToID: func() *uint {
				if len(maintenanceTeam) > 0 {
					return &maintenanceTeam[0].ID
				} else {
					return nil
				}
			}(),
			Title:         "Emergency - Water Leak in Bathroom",
			Description:   "Large water leak under bathroom sink, water spreading to hallway",
			Status:        "completed",
			Priority:      "urgent",
			Category:      "emergency",
			EstimatedCost: 200.00,
			ActualCost:    180.00,
			RequestedAt:   now.AddDate(0, 0, -20),                                         // 20 days ago
			ScheduledAt:   func() *time.Time { t := now.AddDate(0, 0, -20); return &t }(), // Same day
			CompletedAt:   func() *time.Time { t := now.AddDate(0, 0, -19); return &t }(),
			Notes:         "Emergency repair completed. Replaced pipe joint and cleaned up water damage.",
		},
	}

	for i, maintenance := range demoMaintenance {
		// Check if maintenance request already exists
		var existingMaintenance models.Maintenance
		if err := db.DB.Where("title = ? AND property_id = ?", maintenance.Title, maintenance.PropertyID).First(&existingMaintenance).Error; err == nil {
			fmt.Printf("Maintenance request '%s' already exists, skipping...\n", maintenance.Title)
			continue
		}

		// Create maintenance request
		if err := db.DB.Create(&maintenance).Error; err != nil {
			log.Printf("Failed to create maintenance request %d: %v", i+1, err)
			continue
		}

		fmt.Printf("Created demo maintenance request: %s\n", maintenance.Title)
	}
}

func seedInvoices() {
	fmt.Println("Creating demo invoices...")

	// Get tenants, properties, and leases
	var tenants []models.User
	var landlords []models.User
	var properties []models.Property
	var leases []models.Lease

	db.DB.Where("role = ?", "tenant").Find(&tenants)
	db.DB.Where("role = ?", "landlord").Find(&landlords)
	db.DB.Find(&properties)
	db.DB.Find(&leases)

	if len(tenants) == 0 || len(landlords) == 0 || len(leases) == 0 {
		log.Println("Not enough users or leases found, skipping invoice creation")
		return
	}

	now := time.Now()
	invoiceNumber := 1000

	// Create rent invoices for each lease
	for i, lease := range leases {
		if i >= len(landlords) {
			break
		}

		// Monthly rent invoice (current month)
		rentInvoice := models.Invoice{
			TenantID:          lease.TenantID,
			PropertyID:        lease.PropertyID,
			LeaseID:           &lease.ID,
			CreatedByID:       landlords[i%len(landlords)].ID,
			InvoiceNumber:     fmt.Sprintf("INV-%d", invoiceNumber),
			Amount:            lease.MonthlyRent,
			PaidAmount:        lease.MonthlyRent,
			InvoiceDate:       now.AddDate(0, 0, -5), // 5 days ago
			Category:          "rent",
			DueDate:           now.AddDate(0, 0, 25), // 25 days from now
			PaymentStatus:     "paid",
			PaymentMethod:     "bank_transfer",
			Recurring:         true,
			RecurringInterval: "monthly",
			Notes:             "Monthly rent payment",
		}
		invoiceNumber++

		// Previous month rent invoice
		prevRentInvoice := models.Invoice{
			TenantID:          lease.TenantID,
			PropertyID:        lease.PropertyID,
			LeaseID:           &lease.ID,
			CreatedByID:       landlords[i%len(landlords)].ID,
			InvoiceNumber:     fmt.Sprintf("INV-%d", invoiceNumber),
			Amount:            lease.MonthlyRent,
			PaidAmount:        lease.MonthlyRent,
			InvoiceDate:       now.AddDate(0, -1, -5), // Last month
			Category:          "rent",
			DueDate:           now.AddDate(0, -1, 25), // Last month + 25 days
			PaymentStatus:     "paid",
			PaymentMethod:     "bank_transfer",
			Recurring:         true,
			RecurringInterval: "monthly",
			Notes:             "Monthly rent payment",
		}
		invoiceNumber++

		// Utilities invoice
		utilitiesInvoice := models.Invoice{
			TenantID:      lease.TenantID,
			PropertyID:    lease.PropertyID,
			LeaseID:       &lease.ID,
			CreatedByID:   landlords[i%len(landlords)].ID,
			InvoiceNumber: fmt.Sprintf("INV-%d", invoiceNumber),
			Amount:        150.00 + float64(rand.Intn(100)), // Random utilities amount
			PaidAmount:    0,
			InvoiceDate:   now.AddDate(0, 0, -10), // 10 days ago
			Category:      "utilities",
			DueDate:       now.AddDate(0, 0, 20), // 20 days from now
			PaymentStatus: "pending",
			Notes:         "Monthly utilities bill - electricity, gas, water",
		}
		invoiceNumber++

		invoices := []models.Invoice{rentInvoice, prevRentInvoice, utilitiesInvoice}

		for _, invoice := range invoices {
			// Check if invoice already exists
			var existingInvoice models.Invoice
			if err := db.DB.Where("invoice_number = ?", invoice.InvoiceNumber).First(&existingInvoice).Error; err == nil {
				fmt.Printf("Invoice %s already exists, skipping...\n", invoice.InvoiceNumber)
				continue
			}

			// Create invoice
			if err := db.DB.Create(&invoice).Error; err != nil {
				log.Printf("Failed to create invoice %s: %v", invoice.InvoiceNumber, err)
				continue
			}

			fmt.Printf("Created demo invoice: %s (%s)\n", invoice.InvoiceNumber, invoice.Category)
		}
	}

	// Add some overdue invoices
	if len(leases) > 0 && len(landlords) > 0 {
		overdueInvoice := models.Invoice{
			TenantID:      leases[0].TenantID,
			PropertyID:    leases[0].PropertyID,
			LeaseID:       &leases[0].ID,
			CreatedByID:   landlords[0].ID,
			InvoiceNumber: fmt.Sprintf("INV-%d", invoiceNumber),
			Amount:        50.00,
			PaidAmount:    0,
			InvoiceDate:   now.AddDate(0, 0, -40), // 40 days ago
			Category:      "late_fee",
			DueDate:       now.AddDate(0, 0, -10), // 10 days ago (overdue)
			PaymentStatus: "overdue",
			Notes:         "Late payment fee for delayed rent",
		}

		if err := db.DB.Create(&overdueInvoice).Error; err == nil {
			fmt.Printf("Created demo invoice: %s (overdue)\n", overdueInvoice.InvoiceNumber)
		}
	}
}

func seedExpenses() {
	fmt.Println("Creating demo expenses...")

	// Get landlords and properties
	var landlords []models.User
	var properties []models.Property

	db.DB.Where("role = ?", "landlord").Find(&landlords)
	db.DB.Find(&properties)

	if len(landlords) == 0 || len(properties) == 0 {
		log.Println("Not enough landlords or properties found, skipping expense creation")
		return
	}

	now := time.Now()
	expenseNumber := 2000

	demoExpenses := []models.Expense{
		{
			PropertyID:    properties[0].ID,
			CreatedByID:   landlords[0].ID,
			ExpenseNumber: fmt.Sprintf("EXP-%d", expenseNumber),
			Description:   "Monthly property insurance premium",
			Category:      "insurance",
			Amount:        200.00,
			ExpenseDate:   now.AddDate(0, 0, -5),
			VendorName:    "PropertySure Insurance",
			PaymentMethod: "bank_transfer",
			Notes:         "Annual insurance premium - monthly installment",
		},
		{
			PropertyID:    properties[0].ID,
			CreatedByID:   landlords[0].ID,
			ExpenseNumber: fmt.Sprintf("EXP-%d", expenseNumber+1),
			Description:   "Plumbing repair - kitchen faucet replacement",
			Category:      "maintenance",
			Amount:        120.00,
			ExpenseDate:   now.AddDate(0, 0, -8),
			VendorName:    "QuickFix Plumbing",
			PaymentMethod: "card",
			Notes:         "Emergency repair for leaky faucet",
		},
		{
			PropertyID:    properties[1].ID,
			CreatedByID:   landlords[0].ID,
			ExpenseNumber: fmt.Sprintf("EXP-%d", expenseNumber+2),
			Description:   "Garden maintenance and landscaping",
			Category:      "maintenance",
			Amount:        300.00,
			ExpenseDate:   now.AddDate(0, 0, -15),
			VendorName:    "GreenThumb Landscaping",
			PaymentMethod: "bank_transfer",
			Notes:         "Monthly garden maintenance service",
		},
		{
			PropertyID:    properties[1].ID,
			CreatedByID:   landlords[1].ID,
			ExpenseNumber: fmt.Sprintf("EXP-%d", expenseNumber+3),
			Description:   "Property management software subscription",
			Category:      "other",
			Amount:        89.99,
			ExpenseDate:   now.AddDate(0, 0, -3),
			VendorName:    "PropTech Solutions",
			PaymentMethod: "card",
			Notes:         "Monthly SaaS subscription",
		},
		{
			PropertyID:    properties[2].ID,
			CreatedByID:   landlords[1].ID,
			ExpenseNumber: fmt.Sprintf("EXP-%d", expenseNumber+4),
			Description:   "Building maintenance supplies",
			Category:      "supplies",
			Amount:        156.78,
			ExpenseDate:   now.AddDate(0, 0, -12),
			VendorName:    "Builder's Warehouse",
			PaymentMethod: "card",
			Notes:         "Paint, brushes, and cleaning supplies",
		},
		{
			PropertyID:    properties[0].ID,
			CreatedByID:   landlords[0].ID,
			ExpenseNumber: fmt.Sprintf("EXP-%d", expenseNumber+5),
			Description:   "Annual property tax payment",
			Category:      "taxes",
			Amount:        2400.00,
			ExpenseDate:   now.AddDate(0, -2, 0), // 2 months ago
			VendorName:    "Local Council",
			PaymentMethod: "bank_transfer",
			Notes:         "Annual council tax payment",
		},
		{
			PropertyID:    properties[1].ID,
			CreatedByID:   landlords[0].ID,
			ExpenseNumber: fmt.Sprintf("EXP-%d", expenseNumber+6),
			Description:   "Heating system repair",
			Category:      "repairs",
			Amount:        450.00,
			ExpenseDate:   now.AddDate(0, 0, -2),
			VendorName:    "HeatPro Services",
			PaymentMethod: "bank_transfer",
			Notes:         "Thermostat replacement and system check",
		},
	}

	for _, expense := range demoExpenses {
		// Check if expense already exists
		var existingExpense models.Expense
		if err := db.DB.Where("expense_number = ?", expense.ExpenseNumber).First(&existingExpense).Error; err == nil {
			fmt.Printf("Expense %s already exists, skipping...\n", expense.ExpenseNumber)
			continue
		}

		// Create expense
		if err := db.DB.Create(&expense).Error; err != nil {
			log.Printf("Failed to create expense %s: %v", expense.ExpenseNumber, err)
			continue
		}

		fmt.Printf("Created demo expense: %s (%s)\n", expense.ExpenseNumber, expense.Category)
	}
}
