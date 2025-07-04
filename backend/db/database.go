package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/geoo115/property-manager/config"
	"github.com/geoo115/property-manager/logger"
	"github.com/geoo115/property-manager/models"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init(cfg *config.Config) error {
	var err error
	var dialector gorm.Dialector

	// Configure GORM logger
	var logLevel gormLogger.LogLevel
	switch cfg.Server.GinMode {
	case "debug":
		logLevel = gormLogger.Info
	case "release":
		logLevel = gormLogger.Error
	default:
		logLevel = gormLogger.Warn
	}

	gormConfig := &gorm.Config{
		Logger: gormLogger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		QueryFields: true,
	}

	// Configure PostgreSQL dialector
	if cfg.Database.Type != "postgres" {
		return fmt.Errorf("unsupported database type: %s. Only PostgreSQL is supported", cfg.Database.Type)
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)
	dialector = postgres.Open(dsn)

	DB, err = gorm.Open(dialector, gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.LogInfo("Database connection established", logrus.Fields{
		"database_type": cfg.Database.Type,
		"host":          cfg.Database.Host,
		"port":          cfg.Database.Port,
		"database":      cfg.Database.Name,
	})

	// Handle pre-migration data fixes
	if err := handlePreMigrationFixes(); err != nil {
		return fmt.Errorf("failed to handle pre-migration fixes: %w", err)
	}

	// Auto-migrate models
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("failed to auto-migrate: %w", err)
	}

	// Handle post-migration data fixes
	if err := handlePostMigrationFixes(); err != nil {
		return fmt.Errorf("failed to handle post-migration fixes: %w", err)
	}

	// Create indexes
	if err := createIndexes(); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	// Initialize Redis
	if err := InitRedis(cfg); err != nil {
		logger.LogError(err, "Failed to initialize Redis, continuing without it", logrus.Fields{
			"redis_addr": cfg.Redis.Addr,
		})
		// Continue without Redis - rate limiting will be gracefully degraded
	}

	return nil
}

func handlePreMigrationFixes() error {
	// Fix existing users with null usernames before migration
	if err := fixExistingUsersBeforeMigration(); err != nil {
		return fmt.Errorf("failed to fix existing users before migration: %w", err)
	}

	logger.LogInfo("Pre-migration fixes completed", nil)
	return nil
}

func fixExistingUsersBeforeMigration() error {
	// Check if users table exists
	tableExists, err := tableExists("users")
	if err != nil {
		return fmt.Errorf("failed to check if users table exists: %w", err)
	}

	if !tableExists {
		// Table doesn't exist yet, nothing to fix
		logger.LogInfo("Users table doesn't exist yet, skipping pre-migration fixes", nil)
		return nil
	}

	// Check if users table has data
	var count int64
	if err := DB.Table("users").Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}

	if count == 0 {
		// No users to fix
		logger.LogInfo("No users found, skipping pre-migration fixes", nil)
		return nil
	}

	// Check if username column exists
	columnExists, err := columnExists("users", "username")
	if err != nil {
		return fmt.Errorf("failed to check if username column exists: %w", err)
	}

	if !columnExists {
		// Username column doesn't exist, add it as nullable first
		logger.LogInfo("Adding username column as nullable before migration", nil)
		if err := DB.Exec("ALTER TABLE users ADD COLUMN username TEXT").Error; err != nil {
			// If column already exists, ignore the error
			if !strings.Contains(err.Error(), "duplicate column name") {
				return fmt.Errorf("failed to add username column: %w", err)
			}
			logger.LogInfo("Username column already exists, skipping add", nil)
		}
	}

	// Update users with NULL or empty usernames
	// Use PostgreSQL-specific syntax for id conversion
	updateSQL := `
		UPDATE users 
		SET username = COALESCE(NULLIF(username, ''), 'user_' || id::text)
		WHERE username IS NULL OR username = ''
	`

	result := DB.Exec(updateSQL)
	if result.Error != nil {
		return fmt.Errorf("failed to update users with missing usernames: %w", result.Error)
	}

	if result.RowsAffected > 0 {
		logger.LogInfo("Fixed existing users with missing usernames in pre-migration", logrus.Fields{
			"rows_affected": result.RowsAffected,
		})
	}

	return nil
}

func handlePostMigrationFixes() error {
	// Additional post-migration fixes can be added here
	// For example, data validation, cleanup, etc.

	logger.LogInfo("Post-migration fixes completed", nil)
	return nil
}

func autoMigrate() error {
	models := []interface{}{
		&models.User{},
		&models.Property{},
		&models.Unit{},
		&models.Lease{},
		&models.Maintenance{},
		&models.Invoice{},
		&models.Expense{},
		&models.AuditLog{},
	}

	for _, model := range models {
		if err := DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %w", model, err)
		}
	}

	logger.LogInfo("Database auto-migration completed", nil)
	return nil
}

func createIndexes() error {
	// Create additional indexes for better query performance
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_users_email_active ON users(email, is_active);",
		"CREATE INDEX IF NOT EXISTS idx_users_role_active ON users(role, is_active);",
		"CREATE INDEX IF NOT EXISTS idx_properties_owner_available ON properties(owner_id, available);",
		"CREATE INDEX IF NOT EXISTS idx_properties_city_available ON properties(city, available);",
		"CREATE INDEX IF NOT EXISTS idx_leases_dates ON leases(start_date, end_date);",
		"CREATE INDEX IF NOT EXISTS idx_leases_status ON leases(status);",
		"CREATE INDEX IF NOT EXISTS idx_maintenance_status_priority ON maintenance_requests(status, priority);",
		"CREATE INDEX IF NOT EXISTS idx_maintenance_dates ON maintenance_requests(requested_at, scheduled_at);",
		"CREATE INDEX IF NOT EXISTS idx_invoices_status_due ON invoices(payment_status, due_date);",
		"CREATE INDEX IF NOT EXISTS idx_invoices_dates ON invoices(invoice_date, due_date);",
		"CREATE INDEX IF NOT EXISTS idx_expenses_date_category ON expenses(expense_date, category);",
		"CREATE INDEX IF NOT EXISTS idx_audit_logs_user_action ON audit_logs(user_id, action);",
		"CREATE INDEX IF NOT EXISTS idx_audit_logs_entity ON audit_logs(entity_type, entity_id);",
	}

	for _, indexSQL := range indexes {
		if err := DB.Exec(indexSQL).Error; err != nil {
			logger.LogWarning("Failed to create index", logrus.Fields{
				"sql":   indexSQL,
				"error": err.Error(),
			})
			// Continue with other indexes even if one fails
		}
	}

	logger.LogInfo("Database indexes created", nil)
	return nil
}

// Helper functions for database operations
func tableExists(tableName string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA() AND table_name = ?)`
	err := DB.Raw(query, tableName).Scan(&exists).Error
	return exists, err
}

func columnExists(tableName, columnName string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema = CURRENT_SCHEMA() AND table_name = ? AND column_name = ?)`
	err := DB.Raw(query, tableName, columnName).Scan(&exists).Error
	return exists, err
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// HealthCheck checks the database connection
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	return sqlDB.Close()
}
