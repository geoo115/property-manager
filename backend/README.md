# Property Management System - Backend API

[![Go Version](https://img.shields.io/badge/Go-1.20+-blue.svg)](https://golang.org)
[![Gin Framework](https://img.shields.io/badge/Gin-v1.9+-green.svg)](https://github.com/gin-gonic/gin)
[![GORM](https://img.shields.io/badge/GORM-v1.25+-orange.svg)](https://gorm.io)
[![Security](https://img.shields.io/badge/security-hardened-red.svg)](/)
[![Production Ready](https://img.shields.io/badge/production-ready-brightgreen.svg)](/)

## Overview

This is the production-hardened backend API for the Property Management System (PMS). Built with Go, Gin, and GORM, it provides a secure, scalable, and maintainable RESTful API for managing properties, leases, maintenance requests, and financial operations.

## Table of Contents

- [Quick Start](#quick-start)
- [Architecture](#architecture)
- [Security Features](#security-features)
- [Configuration](#configuration)
- [Database Management](#database-management)
- [API Documentation](#api-documentation)
- [Monitoring & Health](#monitoring--health)
- [Development Guide](#development-guide)
- [Production Deployment](#production-deployment)
- [Troubleshooting](#troubleshooting)

## Quick Start

### Prerequisites

- **Go 1.20+** - [Download](https://golang.org/dl/)
- **PostgreSQL** - Database engine
- **Redis** - Caching and rate limiting (optional but recommended)
- **Kafka** - Event streaming (optional)

### Installation

1. **Clone and navigate to backend**
   ```bash
   git clone https://github.com/geoo115/property-manager.git
   cd property-manager/backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start services** (optional but recommended)
   ```bash
   docker-compose up -d
   ```

5. **Run the application**
   ```bash
   go run cmd/main.go
   ```

The API will be available at `http://localhost:8080`

### Health Check

Verify the installation:
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "timestamp": "2025-01-15T10:00:00Z",
  "version": "1.0.0",
  "database": "connected",
  "redis": "connected"
}
```

## Architecture

### Project Structure

```
backend/
├── api/                    # API handlers/controllers
│   ├── auth/              # Authentication endpoints
│   ├── property/          # Property management
│   ├── lease/             # Lease management
│   ├── maintenance/       # Maintenance requests
│   ├── accounting/        # Financial operations
│   └── user/              # User management
├── cmd/                   # Application entry points
│   └── main.go           # Main application
├── config/               # Configuration management
│   └── config.go         # Environment config
├── db/                   # Database layer
│   ├── database.go       # Database initialization
│   └── redis.go          # Redis client
├── events/               # Event handling
│   ├── kafka.go          # Kafka configuration
│   ├── producer.go       # Event production
│   └── consumer.go       # Event consumption
├── middleware/           # HTTP middleware
│   ├── auth.go           # Authentication
│   ├── jwt.go            # JWT handling
│   ├── role.go           # Role-based access
│   ├── rate_limit.go     # Rate limiting
│   ├── cors.go           # CORS headers
│   ├── security.go       # Security headers
│   └── error_handler.go  # Error handling
├── models/               # Data models
│   ├── user.go           # User model
│   ├── property.go       # Property model
│   ├── lease.go          # Lease model
│   ├── maintenance.go    # Maintenance model
│   └── accounting.go     # Financial models
├── router/               # HTTP routing
│   ├── router.go         # Main router
│   └── *_router.go       # Feature routers
├── utils/                # Utility functions
│   └── hash.go           # Password hashing
├── tests/                # Test files
├── logger/               # Logging configuration
├── .env.example          # Environment template
└── docker-compose.yml    # Docker services
```

### Request Flow

```
HTTP Request
    ↓
Gin Router
    ↓
Middleware Pipeline:
├─ CORS Headers
├─ Security Headers
├─ Rate Limiting
├─ Authentication
├─ Authorization
└─ Error Recovery
    ↓
Controller/Handler
    ↓
Business Logic
    ↓
Database Layer (GORM)
    ↓
Response Formation
    ↓
JSON Response
```

## Security Features

### 🔐 Authentication & Authorization

#### JWT Implementation
- **Access Tokens**: 1-hour expiration with secure claims
- **Refresh Tokens**: 24-hour expiration with automatic rotation
- **Token Storage**: HTTP-only cookies for web clients
- **Secure Headers**: CSRF protection and secure cookie attributes

#### Role-Based Access Control (RBAC)
```go
// User roles with specific permissions
type Role string

const (
    RoleAdmin           Role = "admin"           // Full system access
    RoleLandlord        Role = "landlord"        // Property management
    RoleTenant          Role = "tenant"          // Limited access
    RoleMaintenanceTeam Role = "maintenanceTeam" // Maintenance operations
)
```

### 🛡️ Security Hardening

#### Password Security
- **Hashing**: bcrypt with cost factor 12
- **Validation**: Minimum 8 chars, mixed case, numbers, symbols
- **Storage**: Never store plain text passwords

#### Brute Force Protection
- **Rate Limiting**: 5 login attempts per minute per IP
- **Account Lockout**: Exponential backoff on failed attempts
- **IP Tracking**: Redis-based attempt monitoring

#### Security Headers
```go
// Automatically applied security headers
Content-Security-Policy: default-src 'self'
X-Frame-Options: DENY
X-Content-Type-Options: nosniff
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000
```

### 🔒 Data Protection

#### Input Validation
- **Struct Validation**: Gin binding with validation tags
- **SQL Injection Prevention**: GORM parameterized queries
- **XSS Protection**: HTML escaping and sanitization

#### Error Handling
- **Secure Responses**: Never expose internal errors
- **Structured Logging**: Detailed logs for debugging
- **Panic Recovery**: Graceful error recovery

## Configuration

### Environment Variables

Create a `.env` file in the backend directory:

```env
# Database Configuration
DB_TYPE=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=property_manager
DB_SSL_MODE=disable

# Redis Configuration
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here
JWT_EXPIRES_IN=1h
REFRESH_TOKEN_EXPIRES_IN=24h

# Server Configuration
PORT=8080
GIN_MODE=release  # debug, release, test
LOG_LEVEL=info    # debug, info, warn, error

# Kafka Configuration (optional)
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC_MAINTENANCE=maintenance-requests

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60s
AUTH_RATE_LIMIT=5
```

### Configuration Structure

```go
type Config struct {
    Database DatabaseConfig
    Redis    RedisConfig
    JWT      JWTConfig
    Server   ServerConfig
    Kafka    KafkaConfig
    Logging  LoggingConfig
}
```

### Loading Configuration

The application loads configuration from:
1. Environment variables
2. `.env` file in backend directory
3. `.env` file in parent directory (fallback)
4. Default values

## Database Management

### Supported Databases

#### PostgreSQL (Recommended for Production)
```env
DB_TYPE=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=property_manager
DB_SSL_MODE=disable
```

### Database Initialization

The application automatically:
1. **Connects** to PostgreSQL database on startup
2. **Creates tables** if they don't exist
3. **Runs migrations** to update schema
4. **Creates indexes** for performance
5. **Handles existing data** gracefully

### Using Docker Compose

For easy setup, use the provided Docker Compose configuration:

```bash
# Start PostgreSQL and other services
docker-compose up -d

# Or start only PostgreSQL
docker-compose up -d postgres

# Check if services are running
docker-compose ps
```

### Migration System

#### Automatic Migrations
- **Schema Updates**: Automatic table creation and updates
- **Data Preservation**: Existing data is maintained
- **Index Creation**: Performance indexes are created
- **Rollback Support**: Manual rollback procedures documented

#### Migration Features
```go
// Pre-migration data fixes
func (db *Database) preMigrationFixes() error {
    // Handle NULL username values
    // Update missing required fields
    // Data consistency checks
}

// Post-migration optimizations
func (db *Database) postMigrationOptimizations() error {
    // Create performance indexes
    // Update statistics
    // Cleanup unused data
}
```

### Database Connection Management

#### Connection Pooling
```go
// Optimized connection pool settings
MaxOpenConns:    25
MaxIdleConns:    5
ConnMaxLifetime: 1 * time.Hour
ConnMaxIdleTime: 10 * time.Minute
```

#### Health Monitoring
- **Connection Health**: Automatic health checks
- **Reconnection**: Automatic reconnection on failure
- **Metrics**: Connection pool metrics

## API Documentation

### Base URL
```
http://localhost:8080
```

### Authentication

All protected endpoints require JWT authentication:
```bash
curl -H "Authorization: Bearer <token>" <endpoint>
```

### Health Check Endpoint

```http
GET /health
```

Response:
```json
{
  "status": "healthy",
  "timestamp": "2025-01-15T10:00:00Z",
  "version": "1.0.0",
  "database": "connected",
  "redis": "connected",
  "uptime": "2h45m30s"
}
```

### Authentication Endpoints

#### Register User
```http
POST /register
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "role": "tenant"
}
```

#### Login
```http
POST /login
Content-Type: application/json

{
  "username": "john_doe",
  "password": "SecurePass123!"
}
```

#### Refresh Token
```http
POST /refresh
Authorization: Bearer <refresh_token>
```

#### Logout
```http
POST /logout
Authorization: Bearer <token>
```

### Role-Based Endpoints

#### Admin Endpoints
- `GET /admin/users` - List all users
- `POST /admin/users` - Create user
- `PUT /admin/users/:id` - Update user
- `DELETE /admin/users/:id` - Delete user
- `GET /admin/properties` - List all properties
- `GET /admin/leases` - List all leases
- `GET /admin/maintenance` - List all maintenance requests
- `GET /admin/accounting/invoices` - List all invoices
- `GET /admin/accounting/expenses` - List all expenses

#### Landlord Endpoints
- `GET /landlord/properties` - List owned properties
- `POST /landlord/properties` - Create property
- `PUT /landlord/properties/:id` - Update property
- `GET /landlord/leases` - List property leases
- `POST /landlord/leases` - Create lease
- `GET /landlord/maintenance` - List maintenance requests
- `GET /landlord/accounting/invoices` - List invoices
- `GET /landlord/accounting/expenses` - List expenses

#### Tenant Endpoints
- `GET /tenant/leases` - List tenant leases
- `GET /tenant/maintenance` - List maintenance requests
- `POST /tenant/maintenance` - Create maintenance request
- `GET /tenant/invoices` - List tenant invoices

#### Maintenance Team Endpoints
- `GET /maintenance/requests` - List maintenance requests
- `PUT /maintenance/requests/:id` - Update maintenance request
- `GET /maintenance/users` - List users
- `GET /maintenance/properties` - List properties

### Error Responses

All endpoints return consistent error responses:

```json
{
  "success": false,
  "message": "Validation failed",
  "error": [
    {
      "field": "email",
      "message": "must be a valid email address",
      "value": "invalid-email"
    }
  ],
  "timestamp": "2025-01-15T10:00:00Z",
  "request_id": "req_123456789"
}
```

### HTTP Status Codes

- `200 OK` - Success
- `201 Created` - Resource created
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Access denied
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource conflict
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error

## Monitoring & Health

### Health Check System

The application provides comprehensive health monitoring:

#### Health Endpoint
```http
GET /health
```

Checks:
- **Database Connection**: PostgreSQL connectivity
- **Redis Connection**: Cache layer availability
- **Memory Usage**: Current memory consumption
- **Uptime**: Application runtime
- **Version**: Application version info

#### Metrics Collection

Built-in metrics tracking:
- **Request Counts**: Total requests by endpoint
- **Response Times**: Average response times
- **Error Rates**: Error percentage by endpoint
- **Active Connections**: Current database connections
- **Cache Hit Rates**: Redis cache performance

### Logging

#### Structured Logging
```go
// Example log entry
{
  "level": "info",
  "timestamp": "2025-01-15T10:00:00Z",
  "message": "User authenticated",
  "user_id": 123,
  "request_id": "req_123456789",
  "ip": "192.168.1.100",
  "user_agent": "Mozilla/5.0...",
  "duration": "45ms"
}
```

#### Log Levels
- **DEBUG**: Detailed debugging information
- **INFO**: General application information
- **WARN**: Warning messages
- **ERROR**: Error conditions
- **FATAL**: Critical errors that cause shutdown

#### Log Rotation
- **File Size**: 100MB rotation
- **Retention**: 30 days
- **Compression**: Gzip compression for old logs

## Development Guide

### Local Development Setup

1. **Prerequisites**
   ```bash
   go version  # Should be 1.20+
   docker --version
   ```

2. **Database Setup**
   ```bash
   # Using Docker
   docker-compose up -d postgres redis
   
   # Or use the provided script
   ./start_database.sh
   ```

3. **Run in Development Mode**
   ```bash
   # Set development environment
   export GIN_MODE=debug
   export LOG_LEVEL=debug
   
   # Run with hot reload (install air)
   go install github.com/cosmtrek/air@latest
   air
   ```

### Testing

#### Unit Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test file
go test ./tests/auth_test.go
```

#### Integration Tests
```bash
# Run integration tests
go test -tags=integration ./tests/...
```

#### Test Database
Tests use a separate PostgreSQL database:
```bash
# Test database configuration
DB_TYPE=postgres
DB_NAME=property_manager_test
```

### Code Style

#### Go Formatting
```bash
# Format code
go fmt ./...

# Lint code
golangci-lint run
```

#### Naming Conventions
- **Files**: snake_case (user_model.go)
- **Functions**: CamelCase (CreateUser)
- **Variables**: camelCase (userID)
- **Constants**: UPPER_SNAKE_CASE (MAX_RETRY_ATTEMPTS)

### Debugging

#### Debug Mode
```bash
# Enable debug logging
export GIN_MODE=debug
export LOG_LEVEL=debug

# Run with debugger
dlv debug cmd/main.go
```

#### Common Issues

1. **Database Connection**
   ```bash
   # Check database status
   docker-compose ps
   
   # View logs
   docker-compose logs postgres
   ```

2. **Redis Connection**
   ```bash
   # Test Redis connection
   redis-cli ping
   
   # Check Redis logs
   docker-compose logs redis
   ```

3. **JWT Token Issues**
   ```bash
   # Verify JWT secret is set
   echo $JWT_SECRET
   
   # Check token in requests
   curl -H "Authorization: Bearer <token>" http://localhost:8080/health
   ```

## Production Deployment

### Production Checklist

#### Security
- [ ] JWT secret is cryptographically secure
- [ ] Database credentials are secure
- [ ] HTTPS/TLS is enabled
- [ ] Security headers are configured
- [ ] Rate limiting is active
- [ ] Input validation is comprehensive

#### Performance
- [ ] Database connection pool is optimized
- [ ] Redis caching is enabled
- [ ] Static file serving is configured
- [ ] Gzip compression is enabled
- [ ] Database indexes are created

#### Monitoring
- [ ] Health checks are configured
- [ ] Logging is structured and centralized
- [ ] Metrics collection is active
- [ ] Alerting is configured
- [ ] Log rotation is set up

#### Reliability
- [ ] Database migrations are tested
- [ ] Backup procedures are in place
- [ ] Graceful shutdown is implemented
- [ ] Error recovery is tested
- [ ] Load balancing is configured

### Docker Deployment

#### Production Dockerfile
```dockerfile
FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env.example .env

EXPOSE 8080
CMD ["./main"]
```

#### Docker Compose Production
```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - GIN_MODE=release
      - DB_HOST=postgres
      - REDIS_ADDR=redis:6379
    depends_on:
      - postgres
      - redis
    restart: unless-stopped

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: property_manager
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    restart: unless-stopped

volumes:
  postgres_data:
```

### Environment Configuration

#### Production .env
```env
# Production database
DB_TYPE=postgres
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secure-password-here
DB_NAME=property_manager
DB_SSL_MODE=require

# Redis
REDIS_ADDR=redis:6379

# JWT (use strong secret)
JWT_SECRET=your-very-secure-jwt-secret-here

# Server
PORT=8080
GIN_MODE=release
LOG_LEVEL=info

# Rate limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60s
```

### Scaling Considerations

#### Horizontal Scaling
- **Load Balancer**: Nginx or HAProxy
- **Multiple Instances**: Docker Swarm or Kubernetes
- **Session Storage**: Redis for shared sessions
- **Database**: Read replicas for scaling reads

#### Vertical Scaling
- **CPU**: Increase container CPU limits
- **Memory**: Increase memory allocation
- **Database**: Upgrade database instance
- **Connection Pool**: Optimize pool size

## Troubleshooting

### Common Issues

#### Database Connection Issues
```bash
# Check database connectivity
docker-compose ps postgres
docker-compose logs postgres

# Test connection manually
psql -h localhost -U postgres -d property_manager
```

#### Redis Connection Issues
```bash
# Check Redis status
docker-compose ps redis
redis-cli ping

# Clear Redis cache
redis-cli flushall
```

#### JWT Token Issues
```bash
# Verify JWT secret
echo $JWT_SECRET

# Check token expiration
# Use jwt.io to decode tokens
```

#### Migration Issues
```bash
# Check migration logs
grep -i "migration" logs/app.log

# Manual migration reset (caution!)
# Drop tables and restart application
```

### Performance Issues

#### Slow Database Queries
```sql
-- Check slow queries (PostgreSQL)
SELECT query, mean_exec_time, calls
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;

-- Add missing indexes
CREATE INDEX idx_users_email ON users(email);
```

#### Memory Issues
```bash
# Check memory usage
docker stats

# Optimize garbage collection
export GOGC=100
```

### Debug Tools

#### Database Debugging
```bash
# PostgreSQL query logging
echo "log_statement = 'all'" >> postgresql.conf

# Connect to database directly
psql -h localhost -U postgres -d property_manager
```

#### HTTP Debugging
```bash
# Enable request logging
export GIN_MODE=debug

# Use curl with verbose output
curl -v http://localhost:8080/health
```

## Support

### Getting Help

1. **Documentation**: Check this README first
2. **Issues**: Create GitHub issues for bugs
3. **Discussions**: Use GitHub discussions for questions
4. **Contact**: Email support for urgent issues

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

### License

This project is licensed under the MIT License. See LICENSE file for details.

---

**Built with ❤️ using Go, Gin, and GORM**
