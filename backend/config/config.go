package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration values
type Config struct {
	// Database Configuration
	Database DatabaseConfig

	// Redis Configuration
	Redis RedisConfig

	// JWT Configuration
	JWT JWTConfig

	// Server Configuration
	Server ServerConfig

	// Kafka Configuration
	Kafka KafkaConfig

	// Email Configuration
	Email EmailConfig

	// Rate Limiting Configuration
	RateLimit RateLimitConfig

	// File Upload Configuration
	FileUpload FileUploadConfig

	// Security Configuration
	Security SecurityConfig

	// Monitoring Configuration
	Monitoring MonitoringConfig
}

type DatabaseConfig struct {
	Type     string
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret               string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

type ServerConfig struct {
	Host    string
	Port    int
	GinMode string
}

type KafkaConfig struct {
	Broker string
	Topic  string
}

type EmailConfig struct {
	SMTPHost             string
	SMTPPort             int
	SMTPUser             string
	SMTPPass             string
	MaintenanceTeamEmail string
}

type RateLimitConfig struct {
	Requests int
	Duration time.Duration
}

type FileUploadConfig struct {
	MaxFileSize int64
	UploadPath  string
}

type SecurityConfig struct {
	BcryptCost    int
	CORSOrigins   []string
	SecureCookies bool
}

type MonitoringConfig struct {
	EnableMetrics bool
	MetricsPort   int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	// Try loading from current directory first, then parent directory
	if err := godotenv.Load(); err != nil {
		// If not found in current directory, try parent directory
		if err := godotenv.Load("../.env"); err != nil {
			log.Println("No .env file found, using environment variables")
		}
	}

	config := &Config{
		Database: DatabaseConfig{
			Type:     getEnv("DB_TYPE", "postgres"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			Name:     getEnv("DB_NAME", "property_management"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:               getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"),
			AccessTokenDuration:  getEnvDuration("JWT_ACCESS_TOKEN_DURATION", time.Hour),
			RefreshTokenDuration: getEnvDuration("JWT_REFRESH_TOKEN_DURATION", 24*time.Hour),
		},
		Server: ServerConfig{
			Host:    getEnv("SERVER_HOST", "localhost"),
			Port:    getEnvInt("SERVER_PORT", 8080),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Kafka: KafkaConfig{
			Broker: getEnv("KAFKA_BROKER", "localhost:9092"),
			Topic:  getEnv("KAFKA_TOPIC", "property-events"),
		},
		Email: EmailConfig{
			SMTPHost:             getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:             getEnvInt("SMTP_PORT", 587),
			SMTPUser:             getEnv("SMTP_USER", ""),
			SMTPPass:             getEnv("SMTP_PASS", ""),
			MaintenanceTeamEmail: getEnv("MAINTENANCE_TEAM_EMAIL", "maintenance@yourcompany.com"),
		},
		RateLimit: RateLimitConfig{
			Requests: getEnvInt("RATE_LIMIT_REQUESTS", 100),
			Duration: getEnvDuration("RATE_LIMIT_DURATION", time.Minute),
		},
		FileUpload: FileUploadConfig{
			MaxFileSize: getEnvInt64("MAX_FILE_SIZE", 10*1024*1024), // 10MB
			UploadPath:  getEnv("UPLOAD_PATH", "./uploads"),
		},
		Security: SecurityConfig{
			BcryptCost:    getEnvInt("BCRYPT_COST", 12),
			CORSOrigins:   getEnvStringSlice("CORS_ORIGINS", []string{"http://localhost:3000"}),
			SecureCookies: getEnvBool("SECURE_COOKIES", false),
		},
		Monitoring: MonitoringConfig{
			EnableMetrics: getEnvBool("ENABLE_METRICS", true),
			MetricsPort:   getEnvInt("METRICS_PORT", 9090),
		},
	}

	return config, nil
}

// Helper functions for environment variable parsing
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvStringSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Split by comma and trim spaces
		parts := strings.Split(value, ",")
		result := make([]string, len(parts))
		for i, part := range parts {
			result[i] = strings.TrimSpace(part)
		}
		return result
	}
	return defaultValue
}
