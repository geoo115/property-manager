# Property Management System (PMS) Documentation

## Table of Contents
1. [Introduction](#introduction)
2. [User Roles](#user-roles)
3. [Technical Architecture](#technical-architecture)
4. [API Reference](#api-reference)
5. [Security Features](#security-features)
6. [Installation and Setup](#installation-and-setup)
7. [Getting Started Guide](#getting-started-guide)
8. [Advanced Features](#advanced-features)
9. [Future Enhancements](#future-enhancements)
10. [Support and Resources](#support-and-resources)

## Introduction

The Property Management System (PMS) is a comprehensive RESTful API designed to streamline property management operations. This system provides a robust solution for managing rental properties, leases, maintenance requests, and accounting tasks.

### Core Functionality

- **User Management**: Role-based user operations with secure authentication
- **Property Management**: Complete lifecycle management for rental properties
- **Lease Administration**: Digital lease creation and management
- **Maintenance Request Handling**: Event-driven workflow for maintenance issues
- **Financial Operations**: Invoice and expense tracking with role-appropriate access

### Technology Stack

- **Backend**: Go, Gin framework
- **Database**: SQLite with GORM
- **Caching**: Redis
- **Messaging**: Apache Kafka
- **Authentication**: JWT (JSON Web Tokens)

## User Roles

The PMS implements role-based access control with four distinct user types:

### Admin
- Full system access with unrestricted capabilities
- User, property, lease, and accounting management
- System-wide configuration and monitoring

### Landlord
- Management of owned properties
- Lease creation and oversight
- Maintenance request initiation and tracking
- Financial record access (invoices and expenses)

### Tenant
- Lease viewing and management
- Maintenance request submission
- Invoice access

### Maintenance Team
- Maintenance request updates
- Property and user information access
- Task management and reporting

## Technical Architecture

The PMS implements a layered architecture designed for modularity, maintainability, and scalability:

### Web Layer
- Built with Gin framework
- RESTful API endpoints with efficient routing
- Middleware support for authentication, RBAC, and rate limiting

### Database Layer
- SQLite database accessed via GORM
- Structured data models for properties, leases, maintenance requests, and accounting
- Transactional operations ensuring data integrity

### Caching Layer
- Redis implementation for performance optimization
- 10-minute TTL balancing freshness and performance
- Automatic cache invalidation on write operations

### Messaging Layer
- Apache Kafka for event-driven architecture
- Asynchronous processing of maintenance requests
- Decoupled communication between system components

### Security Layer
- JWT-based authentication
- Role-based access control (RBAC)
- Rate limiting to prevent abuse

## API Reference

### Authentication Endpoints

| Method | Endpoint | Description | Sample Request |
|--------|----------|-------------|----------------|
| POST | /login | Authenticate user | `{"username": "john_doe", "password": "securepassword123"}` |
| POST | /register | Register new user | `{"username": "john_doe", "password": "securepassword123", "email": "john@example.com", "role": "tenant"}` |
| POST | /refresh | Refresh access token | - |
| POST | /logout | Invalidate JWT token | - |

### Role-Specific Endpoints

#### Admin Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /admin/users | List all users |
| GET | /admin/users/:id | Get user by ID |
| POST | /admin/users | Create a new user |
| PUT | /admin/users/:id | Update user |
| DELETE | /admin/users/:id | Delete user |
| GET | /admin/properties | List all properties |
| GET | /admin/properties/:id | Get property by ID |
| POST | /admin/properties | Create a property |
| PUT | /admin/properties/:id | Update property |
| DELETE | /admin/properties/:id | Delete property |
| GET | /admin/leases | List all leases |
| GET | /admin/leases/:id | Get lease by ID |
| POST | /admin/leases | Create a lease |
| PUT | /admin/leases/:id | Update lease |
| DELETE | /admin/leases/:id | Delete lease |
| GET | /admin/maintenances | List all maintenance requests |
| GET | /admin/maintenance/:id | Get maintenance request by ID |
| POST | /admin/properties/:propertyID/maintenances | Create maintenance request |
| PUT | /admin/maintenance/:id | Update maintenance request |
| DELETE | /admin/maintenance/:id | Delete maintenance request |
| GET | /admin/accounting/invoices | List all invoices |
| POST | /admin/accounting/invoices | Create invoice |
| PUT | /admin/accounting/invoices/:id | Update invoice |
| DELETE | /admin/accounting/invoices/:id | Delete invoice |
| GET | /admin/accounting/expenses | List all expenses |
| POST | /admin/accounting/expense | Create expense |
| PUT | /admin/accounting/expense/:id | Update expense |
| DELETE | /admin/accounting/expense/:id | Delete expense |

#### Landlord Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /landlord/properties | List owned properties |
| GET | /landlord/properties/:id | Get owned property by ID |
| POST | /landlord/properties | Create a property |
| PUT | /landlord/properties/:id | Update owned property |
| DELETE | /landlord/properties/:id | Delete owned property |
| GET | /landlord/leases | List leases for owned properties |
| GET | /landlord/leases/:id | Get lease by ID for owned property |
| POST | /landlord/leases | Create a lease for owned property |
| GET | /landlord/properties/:id/maintenances | List maintenance requests for property |
| POST | /landlord/properties/:id/maintenances | Create maintenance request |
| GET | /landlord/invoices | List invoices for owned properties |
| GET | /landlord/expenses | List expenses for owned properties |
| POST | /landlord/expenses | Create expense for owned property |

#### Tenant Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /tenant/leases | List tenant's leases |
| GET | /tenant/leases/:id | Get lease by ID |
| GET | /tenant/leases/active | Get active lease |
| GET | /tenant/leases/:id/maintenance | List maintenance requests for lease |
| POST | /tenant/leases/:id/maintenance | Create maintenance request |
| GET | /tenant/invoices | List tenant's invoices |

#### Maintenance Team Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /maintenanceTeam/maintenances | List all maintenance requests |
| GET | /maintenanceTeam/maintenance/:id | Get maintenance request by ID |
| PUT | /maintenanceTeam/maintenance/:id | Update maintenance request |
| GET | /maintenanceTeam/users | List all users |
| GET | /maintenanceTeam/properties | List all properties |

### Sample Request Payloads

#### User Registration
```json
{
  "username": "john_doe",
  "first_name": "John",
  "last_name": "Doe",
  "password": "securepassword123",
  "email": "john@example.com",
  "role": "tenant",
  "phone": "+1-555-123-4567"
}
```

#### Property Creation
```json
{
  "name": "Downtown Apartment",
  "address": "123 Main St",
  "city": "New York",
  "price": 2000.00,
  "owner_id": 1,
  "available": true,
  "bedrooms": 2,
  "bathrooms": 1,
  "square_feet": 850,
  "post_code": "10001"
}
```

#### Lease Creation
```json
{
  "property_id": 1,
  "tenant_id": 2,
  "start_date": "2023-01-01",
  "end_date": "2024-01-01",
  "monthly_rent": 1500,
  "security_deposit": 3000
}
```

#### Maintenance Request
```json
{
  "description": "Fix leaky faucet in master bathroom"
}
```

#### Invoice Creation
```json
{
  "tenant_id": 2,
  "property_id": 1,
  "amount": 1500.00,
  "paid_amount": 0.00,
  "invoice_date": "2023-01-01",
  "due_date": "2023-02-01",
  "category": "rent",
  "payment_status": "unpaid"
}
```

#### Expense Recording
```json
{
  "property_id": 1,
  "description": "Plumbing repair",
  "category": "maintenance",
  "amount": 200.00,
  "expense_date": "2023-02-15"
}
```

## Security Features

### Authentication and Authorization

The PMS implements a robust security model using:

- **JWT Authentication**: JSON Web Tokens secure all protected endpoints
  - Access tokens valid for 1 hour
  - Refresh tokens valid for 24 hours
  - Tokens delivered via HTTP-only cookies

- **Role-Based Access Control (RBAC)**: Middleware enforces role-specific permissions
  - Role verification at request time
  - Resource ownership validation
  - Route-based access restrictions

### Additional Security Measures

- Password hashing using industry-standard algorithms
- HTTPS transport (assumed in production)
- No sensitive data in URL parameters
- Input validation and sanitization
- Rate limiting protection

### Rate Limiting

Protection against abuse through Redis-based rate limiting:

- API Endpoints: 100 requests per minute per IP
- Authentication Endpoints: 5 requests per minute per IP
- Implementation: Middleware integration with configurable limits

### Error Handling

The system implements consistent error handling throughout:

**HTTP Status Codes**
- 401 Unauthorized: Authentication issues
- 403 Forbidden: Authorization failures
- 404 Not Found: Resource does not exist
- 400 Bad Request: Invalid input data
- 429 Too Many Requests: Rate limit exceeded
- 500 Internal Server Error: Server-side failures

**Error Response Format**
```json
{
  "error": "Access denied",
  "details": "Required role: admin"
}
```

## Installation and Setup

### Prerequisites

Before installing the PMS, ensure you have the following:

- Go 1.20 or higher
- Docker for container management
- Git for repository access

### Installation Steps

1. **Clone the Repository**
   ```bash
   git clone https://github.com/geoo115/property-manager.git
   cd property-manager
   ```

2. **Start Dependencies with Docker**
   ```bash
   docker compose up -d
   ```
   This configures:
   - Redis for caching
   - Kafka for event messaging

3. **Configure Environment Variables**
   Create a `.env` file in the root directory:
   ```
   JWT_SECRET=your_jwt_secret_here
   KAFKA_BROKER=localhost:9092
   REDIS_ADDR=localhost:6379
   ```

4. **Run the Application**
   ```bash
   go run main.go
   ```
   The server starts on port 8080 by default.

### Verification

Verify your installation by checking:

- Redis container status: `docker ps | grep redis`
- Kafka container status: `docker ps | grep kafka`
- API accessibility: `curl http://localhost:8080/health` (returns 200 OK)

## Getting Started Guide

Follow these steps to quickly start using the PMS:

1. **Register a Landlord Account**
   ```bash
   curl -X POST http://localhost:8080/register \
     -H "Content-Type: application/json" \
     -d '{"username": "landlord1", "password": "pass123", "email": "landlord@example.com", "role": "landlord"}'
   ```

2. **Log In to Get JWT Token**
   ```bash
   curl -X POST http://localhost:8080/login \
     -H "Content-Type: application/json" \
     -d '{"username": "landlord1", "password": "pass123"}'
   ```

3. **Create a Property**
   ```bash
   curl -X POST http://localhost:8080/landlord/properties \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <your_token>" \
     -d '{"name": "Downtown Apt", "address": "123 Main St", "price": 2000.00, "owner_id": 1, "available": true}'
   ```

4. **Create a Lease**
   ```bash
   curl -X POST http://localhost:8080/landlord/leases \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <your_token>" \
     -d '{"property_id": 1, "tenant_id": 2, "start_date": "2023-01-01", "end_date": "2024-01-01", "monthly_rent": 1500}'
   ```

5. **Submit a Maintenance Request**
   ```bash
   curl -X POST http://localhost:8080/tenant/leases/1/maintenance \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <tenant_token>" \
     -d '{"description": "Fix leaky faucet"}'
   ```

## Advanced Features

### Caching Strategy

The caching implementation optimizes performance while maintaining data freshness:

**Redis Configuration**
- TTL: 10-minute expiration for cached items
- Key Patterns:
  - Properties: "properties:<id>"
  - Leases: "leases:<id>"
  - Users: "users:<id>"
  - Maintenance Requests: Role-specific keys (e.g., "maintenances:tenant:<user_id>:lease:<lease_id>")

**Cache Invalidation**
- Automatic on write operations (create/update/delete)
- Implemented via Redis DEL commands
- Ensures data consistency across the system

### Event-Driven Architecture

The PMS leverages Kafka for asynchronous, event-driven processing:

**Maintenance Request Workflow**
1. Request created via API endpoint
2. Event published to maintenance-requests Kafka topic
3. Consumer processes event, performing:
   - Notification to maintenance team
   - Status updates
   - Audit logging

**Message Format**
```json
{
  "id": 45,
  "property_id": 12,
  "description": "Broken faucet",
  "status": "pending",
  "requested_at": "2025-02-25T21:00:00Z"
}
```

## Future Enhancements

The PMS roadmap includes several planned enhancements:

### Payment Processing Integration
- Payment gateway support for rent collection
- Automatic receipt generation
- Payment status tracking

### Notification System
- Email notifications for critical events
- SMS alerts for maintenance updates
- Push notifications via mobile app

### Document Management
- PDF generation for leases and invoices
- Document storage and retrieval
- Signature integration for lease signing

### Analytics Dashboard
- Property performance metrics
- Maintenance history analytics
- Financial reporting and forecasting

### Mobile Application
- Native mobile apps for iOS and Android
- Responsive web interface
- Offline capability for key features

## Support and Resources

### Community and Assistance
- GitHub Repository: github.com/geoo115/property-manager
- Issue Tracker: Create issues for bugs or feature requests
- Documentation Wiki: Additional detailed documentation (future)
- Contact: For direct support, contact geoo115@gmail.com

### Contributing
Contributions to the PMS are welcome! Please:
- Fork the repository
- Create a feature branch
- Submit a pull request with detailed description
- Follow the project's coding standards

### License
This project is licensed under the MIT License. See the LICENSE file for details.