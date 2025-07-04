# Property Management System - API Documentation

## API Overview

The Property Management System (PMS) provides a RESTful API for managing properties, leases, maintenance requests, and financial operations. The API is built with Go and Gin framework, featuring JWT authentication, role-based access control, and comprehensive error handling.

### Base URL
```
https://your-domain.com/api/v1
```

### Authentication
All protected endpoints require JWT authentication via the `Authorization` header:
```
Authorization: Bearer <jwt_token>
```

### Response Format
All API responses follow a consistent format:

**Success Response:**
```json
{
  "success": true,
  "data": { ... },
  "message": "Operation completed successfully",
  "timestamp": "2025-01-15T10:00:00Z"
}
```

**Error Response:**
```json
{
  "success": false,
  "error": [
    {
      "field": "email",
      "message": "must be a valid email address",
      "value": "invalid-email"
    }
  ],
  "message": "Validation failed",
  "timestamp": "2025-01-15T10:00:00Z",
  "request_id": "req_123456789"
}
```

### HTTP Status Codes

- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Access denied
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource conflict
- `422 Unprocessable Entity` - Validation error
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error

### Rate Limiting

The API implements rate limiting to prevent abuse:

- **Public endpoints**: 50 requests per minute per IP
- **Authenticated endpoints**: 200 requests per minute per user
- **Authentication endpoints**: 5 requests per minute per IP

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 200
X-RateLimit-Remaining: 199
X-RateLimit-Reset: 1642248000
```

## Authentication Endpoints

### Register User
Create a new user account.

**Endpoint:** `POST /register`

**Request Body:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+1-555-123-4567",
  "role": "tenant"
}
```

**Field Validation:**
- `username`: Required, 3-50 characters, alphanumeric and underscores only
- `email`: Required, valid email format
- `password`: Required, minimum 8 characters, must contain uppercase, lowercase, number, and special character
- `first_name`: Optional, maximum 50 characters
- `last_name`: Optional, maximum 50 characters
- `phone`: Optional, valid phone number format
- `role`: Required, one of: "admin", "landlord", "tenant", "maintenanceTeam"

**Success Response (201):**
```json
{
  "success": true,
  "data": {
    "id": 123,
    "username": "john_doe",
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "tenant",
    "is_active": true,
    "created_at": "2025-01-15T10:00:00Z"
  },
  "message": "User created successfully"
}
```

### Login
Authenticate user and receive JWT tokens.

**Endpoint:** `POST /login`

**Request Body:**
```json
{
  "username": "john_doe",
  "password": "SecurePass123!"
}
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "user": {
      "id": 123,
      "username": "john_doe",
      "email": "john@example.com",
      "role": "tenant"
    }
  },
  "message": "Login successful"
}
```

### Refresh Token
Get a new access token using refresh token.

**Endpoint:** `POST /refresh`

**Headers:**
```
Authorization: Bearer <refresh_token>
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600
  },
  "message": "Token refreshed successfully"
}
```

### Logout
Invalidate current tokens.

**Endpoint:** `POST /logout`

**Headers:**
```
Authorization: Bearer <access_token>
```

**Success Response (200):**
```json
{
  "success": true,
  "message": "Logout successful"
}
```

## User Management Endpoints

### Get All Users (Admin Only)
Retrieve a list of all users.

**Endpoint:** `GET /admin/users`

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Number of users per page (default: 10, max: 100)
- `role`: Filter by role
- `active`: Filter by active status (true/false)
- `search`: Search in username, email, first_name, last_name

**Example:** `GET /admin/users?page=1&limit=10&role=tenant&active=true`

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "users": [
      {
        "id": 123,
        "username": "john_doe",
        "email": "john@example.com",
        "first_name": "John",
        "last_name": "Doe",
        "role": "tenant",
        "is_active": true,
        "created_at": "2025-01-15T10:00:00Z",
        "updated_at": "2025-01-15T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 25,
      "pages": 3
    }
  },
  "message": "Users retrieved successfully"
}
```

### Get User by ID
Retrieve a specific user by ID.

**Endpoint:** `GET /admin/users/:id`

**Path Parameters:**
- `id`: User ID (integer)

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "id": 123,
    "username": "john_doe",
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1-555-123-4567",
    "role": "tenant",
    "is_active": true,
    "created_at": "2025-01-15T10:00:00Z",
    "updated_at": "2025-01-15T10:00:00Z"
  },
  "message": "User retrieved successfully"
}
```

### Create User (Admin Only)
Create a new user account.

**Endpoint:** `POST /admin/users`

**Request Body:** Same as registration endpoint

**Success Response (201):** Same as registration endpoint

### Update User
Update user information.

**Endpoint:** `PUT /admin/users/:id`

**Path Parameters:**
- `id`: User ID (integer)

**Request Body:**
```json
{
  "first_name": "John Updated",
  "last_name": "Doe Updated",
  "phone": "+1-555-987-6543",
  "is_active": true
}
```

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "id": 123,
    "username": "john_doe",
    "email": "john@example.com",
    "first_name": "John Updated",
    "last_name": "Doe Updated",
    "phone": "+1-555-987-6543",
    "role": "tenant",
    "is_active": true,
    "updated_at": "2025-01-15T11:00:00Z"
  },
  "message": "User updated successfully"
}
```

### Delete User (Admin Only)
Delete a user account.

**Endpoint:** `DELETE /admin/users/:id`

**Path Parameters:**
- `id`: User ID (integer)

**Success Response (200):**
```json
{
  "success": true,
  "message": "User deleted successfully"
}
```

## Property Management Endpoints

### Get All Properties
Retrieve properties based on user role.

**Endpoint:** 
- `GET /admin/properties` (Admin - all properties)
- `GET /landlord/properties` (Landlord - owned properties)

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Number of properties per page (default: 10, max: 100)
- `city`: Filter by city
- `available`: Filter by availability (true/false)
- `min_price`: Minimum price filter
- `max_price`: Maximum price filter
- `bedrooms`: Filter by number of bedrooms
- `bathrooms`: Filter by number of bathrooms

**Example:** `GET /landlord/properties?city=New York&available=true&min_price=1000&max_price=3000`

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "properties": [
      {
        "id": 1,
        "name": "Downtown Apartment",
        "address": "123 Main St",
        "city": "New York",
        "post_code": "10001",
        "owner_id": 123,
        "price": 2000.00,
        "bedrooms": 2,
        "bathrooms": 1,
        "square_feet": 850,
        "available": true,
        "created_at": "2025-01-15T10:00:00Z",
        "updated_at": "2025-01-15T10:00:00Z",
        "owner": {
          "id": 123,
          "username": "landlord_user",
          "email": "landlord@example.com"
        }
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 5,
      "pages": 1
    }
  },
  "message": "Properties retrieved successfully"
}
```

### Get Property by ID
Retrieve a specific property by ID.

**Endpoint:** 
- `GET /admin/properties/:id` (Admin)
- `GET /landlord/properties/:id` (Landlord - owned properties only)

**Path Parameters:**
- `id`: Property ID (integer)

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Downtown Apartment",
    "address": "123 Main St",
    "city": "New York",
    "post_code": "10001",
    "owner_id": 123,
    "price": 2000.00,
    "bedrooms": 2,
    "bathrooms": 1,
    "square_feet": 850,
    "available": true,
    "created_at": "2025-01-15T10:00:00Z",
    "updated_at": "2025-01-15T10:00:00Z",
    "owner": {
      "id": 123,
      "username": "landlord_user",
      "email": "landlord@example.com"
    },
    "leases": [
      {
        "id": 1,
        "tenant_id": 456,
        "start_date": "2025-01-01",
        "end_date": "2025-12-31",
        "monthly_rent": 1800.00,
        "status": "active"
      }
    ]
  },
  "message": "Property retrieved successfully"
}
```

### Create Property
Create a new property.

**Endpoint:** 
- `POST /admin/properties` (Admin)
- `POST /landlord/properties` (Landlord)

**Request Body:**
```json
{
  "name": "Downtown Apartment",
  "address": "123 Main St",
  "city": "New York",
  "post_code": "10001",
  "owner_id": 123,
  "price": 2000.00,
  "bedrooms": 2,
  "bathrooms": 1,
  "square_feet": 850,
  "available": true
}
```

**Field Validation:**
- `name`: Required, maximum 100 characters
- `address`: Required, maximum 255 characters
- `city`: Required, maximum 50 characters
- `post_code`: Optional, maximum 10 characters
- `owner_id`: Required, valid user ID with landlord role
- `price`: Required, positive decimal value
- `bedrooms`: Optional, positive integer
- `bathrooms`: Optional, positive integer
- `square_feet`: Optional, positive integer
- `available`: Optional, boolean (default: true)

**Success Response (201):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Downtown Apartment",
    "address": "123 Main St",
    "city": "New York",
    "post_code": "10001",
    "owner_id": 123,
    "price": 2000.00,
    "bedrooms": 2,
    "bathrooms": 1,
    "square_feet": 850,
    "available": true,
    "created_at": "2025-01-15T10:00:00Z",
    "updated_at": "2025-01-15T10:00:00Z"
  },
  "message": "Property created successfully"
}
```

### Update Property
Update property information.

**Endpoint:** 
- `PUT /admin/properties/:id` (Admin)
- `PUT /landlord/properties/:id` (Landlord - owned properties only)

**Path Parameters:**
- `id`: Property ID (integer)

**Request Body:** Same as create property (all fields optional)

**Success Response (200):** Same as create property with updated values

### Delete Property
Delete a property.

**Endpoint:** 
- `DELETE /admin/properties/:id` (Admin)
- `DELETE /landlord/properties/:id` (Landlord - owned properties only)

**Path Parameters:**
- `id`: Property ID (integer)

**Success Response (200):**
```json
{
  "success": true,
  "message": "Property deleted successfully"
}
```

## Lease Management Endpoints

### Get All Leases
Retrieve leases based on user role.

**Endpoint:** 
- `GET /admin/leases` (Admin - all leases)
- `GET /landlord/leases` (Landlord - leases for owned properties)
- `GET /tenant/leases` (Tenant - tenant's leases)

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Number of leases per page (default: 10, max: 100)
- `status`: Filter by status (active, expired, terminated)
- `property_id`: Filter by property ID
- `tenant_id`: Filter by tenant ID
- `start_date`: Filter by start date (YYYY-MM-DD)
- `end_date`: Filter by end date (YYYY-MM-DD)

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "leases": [
      {
        "id": 1,
        "property_id": 1,
        "tenant_id": 456,
        "start_date": "2025-01-01",
        "end_date": "2025-12-31",
        "monthly_rent": 1800.00,
        "security_deposit": 3600.00,
        "status": "active",
        "created_at": "2025-01-15T10:00:00Z",
        "updated_at": "2025-01-15T10:00:00Z",
        "property": {
          "id": 1,
          "name": "Downtown Apartment",
          "address": "123 Main St",
          "city": "New York"
        },
        "tenant": {
          "id": 456,
          "username": "tenant_user",
          "email": "tenant@example.com",
          "first_name": "Jane",
          "last_name": "Smith"
        }
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 3,
      "pages": 1
    }
  },
  "message": "Leases retrieved successfully"
}
```

### Get Lease by ID
Retrieve a specific lease by ID.

**Endpoint:** 
- `GET /admin/leases/:id` (Admin)
- `GET /landlord/leases/:id` (Landlord - owned properties only)
- `GET /tenant/leases/:id` (Tenant - tenant's leases only)

**Path Parameters:**
- `id`: Lease ID (integer)

**Success Response (200):** Same as lease object in list response

### Create Lease
Create a new lease.

**Endpoint:** 
- `POST /admin/leases` (Admin)
- `POST /landlord/leases` (Landlord)

**Request Body:**
```json
{
  "property_id": 1,
  "tenant_id": 456,
  "start_date": "2025-01-01",
  "end_date": "2025-12-31",
  "monthly_rent": 1800.00,
  "security_deposit": 3600.00
}
```

**Field Validation:**
- `property_id`: Required, valid property ID
- `tenant_id`: Required, valid user ID with tenant role
- `start_date`: Required, valid date (YYYY-MM-DD)
- `end_date`: Required, valid date after start_date
- `monthly_rent`: Required, positive decimal value
- `security_deposit`: Optional, positive decimal value

**Success Response (201):** Same as lease object with generated ID

### Update Lease
Update lease information.

**Endpoint:** 
- `PUT /admin/leases/:id` (Admin)
- `PUT /landlord/leases/:id` (Landlord - owned properties only)

**Path Parameters:**
- `id`: Lease ID (integer)

**Request Body:** Same as create lease (all fields optional)

**Success Response (200):** Same as create lease with updated values

### Delete Lease
Delete a lease.

**Endpoint:** 
- `DELETE /admin/leases/:id` (Admin)
- `DELETE /landlord/leases/:id` (Landlord - owned properties only)

**Path Parameters:**
- `id`: Lease ID (integer)

**Success Response (200):**
```json
{
  "success": true,
  "message": "Lease deleted successfully"
}
```

## Maintenance Request Endpoints

### Get All Maintenance Requests
Retrieve maintenance requests based on user role.

**Endpoint:** 
- `GET /admin/maintenance` (Admin - all requests)
- `GET /landlord/properties/:property_id/maintenance` (Landlord - requests for owned properties)
- `GET /tenant/leases/:lease_id/maintenance` (Tenant - requests for tenant's leases)
- `GET /maintenance/requests` (Maintenance Team - all requests)

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Number of requests per page (default: 10, max: 100)
- `status`: Filter by status (pending, in_progress, completed, cancelled)
- `priority`: Filter by priority (low, medium, high, urgent)
- `property_id`: Filter by property ID
- `tenant_id`: Filter by tenant ID

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "maintenance_requests": [
      {
        "id": 1,
        "property_id": 1,
        "tenant_id": 456,
        "description": "Fix leaky faucet in master bathroom",
        "status": "pending",
        "priority": "medium",
        "created_at": "2025-01-15T10:00:00Z",
        "updated_at": "2025-01-15T10:00:00Z",
        "property": {
          "id": 1,
          "name": "Downtown Apartment",
          "address": "123 Main St"
        },
        "tenant": {
          "id": 456,
          "username": "tenant_user",
          "email": "tenant@example.com"
        }
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 2,
      "pages": 1
    }
  },
  "message": "Maintenance requests retrieved successfully"
}
```

### Get Maintenance Request by ID
Retrieve a specific maintenance request by ID.

**Endpoint:** 
- `GET /admin/maintenance/:id` (Admin)
- `GET /maintenance/requests/:id` (Maintenance Team)

**Path Parameters:**
- `id`: Maintenance request ID (integer)

**Success Response (200):** Same as maintenance request object in list response

### Create Maintenance Request
Create a new maintenance request.

**Endpoint:** 
- `POST /admin/properties/:property_id/maintenance` (Admin)
- `POST /landlord/properties/:property_id/maintenance` (Landlord)
- `POST /tenant/leases/:lease_id/maintenance` (Tenant)

**Path Parameters:**
- `property_id`: Property ID (integer) for admin/landlord
- `lease_id`: Lease ID (integer) for tenant

**Request Body:**
```json
{
  "description": "Fix leaky faucet in master bathroom",
  "priority": "medium"
}
```

**Field Validation:**
- `description`: Required, maximum 1000 characters
- `priority`: Optional, one of: "low", "medium", "high", "urgent" (default: "medium")

**Success Response (201):** Same as maintenance request object with generated ID

### Update Maintenance Request
Update maintenance request information.

**Endpoint:** 
- `PUT /admin/maintenance/:id` (Admin)
- `PUT /maintenance/requests/:id` (Maintenance Team)

**Path Parameters:**
- `id`: Maintenance request ID (integer)

**Request Body:**
```json
{
  "status": "in_progress",
  "priority": "high",
  "notes": "Started work on the issue"
}
```

**Field Validation:**
- `status`: Optional, one of: "pending", "in_progress", "completed", "cancelled"
- `priority`: Optional, one of: "low", "medium", "high", "urgent"
- `notes`: Optional, maximum 1000 characters

**Success Response (200):** Same as maintenance request object with updated values

### Delete Maintenance Request
Delete a maintenance request.

**Endpoint:** 
- `DELETE /admin/maintenance/:id` (Admin)

**Path Parameters:**
- `id`: Maintenance request ID (integer)

**Success Response (200):**
```json
{
  "success": true,
  "message": "Maintenance request deleted successfully"
}
```

## Financial Management Endpoints

### Invoices

#### Get All Invoices
Retrieve invoices based on user role.

**Endpoint:** 
- `GET /admin/accounting/invoices` (Admin - all invoices)
- `GET /landlord/accounting/invoices` (Landlord - invoices for owned properties)
- `GET /tenant/invoices` (Tenant - tenant's invoices)

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Number of invoices per page (default: 10, max: 100)
- `payment_status`: Filter by payment status (paid, unpaid, overdue, partial)
- `category`: Filter by category (rent, utilities, fees, other)
- `property_id`: Filter by property ID
- `tenant_id`: Filter by tenant ID
- `start_date`: Filter by invoice date (YYYY-MM-DD)
- `end_date`: Filter by invoice date (YYYY-MM-DD)

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "invoices": [
      {
        "id": 1,
        "tenant_id": 456,
        "property_id": 1,
        "amount": 1800.00,
        "paid_amount": 1800.00,
        "invoice_date": "2025-01-01",
        "due_date": "2025-01-31",
        "category": "rent",
        "payment_status": "paid",
        "created_at": "2025-01-15T10:00:00Z",
        "updated_at": "2025-01-15T10:00:00Z",
        "property": {
          "id": 1,
          "name": "Downtown Apartment",
          "address": "123 Main St"
        },
        "tenant": {
          "id": 456,
          "username": "tenant_user",
          "email": "tenant@example.com"
        }
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 12,
      "pages": 2
    }
  },
  "message": "Invoices retrieved successfully"
}
```

#### Create Invoice
Create a new invoice.

**Endpoint:** 
- `POST /admin/accounting/invoices` (Admin)
- `POST /landlord/accounting/invoices` (Landlord)

**Request Body:**
```json
{
  "tenant_id": 456,
  "property_id": 1,
  "amount": 1800.00,
  "invoice_date": "2025-01-01",
  "due_date": "2025-01-31",
  "category": "rent",
  "description": "Monthly rent for January 2025"
}
```

**Field Validation:**
- `tenant_id`: Required, valid user ID with tenant role
- `property_id`: Required, valid property ID
- `amount`: Required, positive decimal value
- `invoice_date`: Required, valid date (YYYY-MM-DD)
- `due_date`: Required, valid date after invoice_date
- `category`: Required, one of: "rent", "utilities", "fees", "other"
- `description`: Optional, maximum 255 characters

**Success Response (201):** Same as invoice object with generated ID

### Expenses

#### Get All Expenses
Retrieve expenses based on user role.

**Endpoint:** 
- `GET /admin/accounting/expenses` (Admin - all expenses)
- `GET /landlord/accounting/expenses` (Landlord - expenses for owned properties)

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Number of expenses per page (default: 10, max: 100)
- `category`: Filter by category (maintenance, utilities, insurance, other)
- `property_id`: Filter by property ID
- `start_date`: Filter by expense date (YYYY-MM-DD)
- `end_date`: Filter by expense date (YYYY-MM-DD)

**Success Response (200):**
```json
{
  "success": true,
  "data": {
    "expenses": [
      {
        "id": 1,
        "property_id": 1,
        "description": "Plumbing repair - bathroom faucet",
        "category": "maintenance",
        "amount": 150.00,
        "expense_date": "2025-01-10",
        "created_at": "2025-01-15T10:00:00Z",
        "updated_at": "2025-01-15T10:00:00Z",
        "property": {
          "id": 1,
          "name": "Downtown Apartment",
          "address": "123 Main St"
        }
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 8,
      "pages": 1
    }
  },
  "message": "Expenses retrieved successfully"
}
```

#### Create Expense
Create a new expense.

**Endpoint:** 
- `POST /admin/accounting/expenses` (Admin)
- `POST /landlord/accounting/expenses` (Landlord)

**Request Body:**
```json
{
  "property_id": 1,
  "description": "Plumbing repair - bathroom faucet",
  "category": "maintenance",
  "amount": 150.00,
  "expense_date": "2025-01-10"
}
```

**Field Validation:**
- `property_id`: Required, valid property ID
- `description`: Required, maximum 255 characters
- `category`: Required, one of: "maintenance", "utilities", "insurance", "other"
- `amount`: Required, positive decimal value
- `expense_date`: Required, valid date (YYYY-MM-DD)

**Success Response (201):** Same as expense object with generated ID

## Health Check Endpoint

### System Health
Get system health status.

**Endpoint:** `GET /health`

**Success Response (200):**
```json
{
  "status": "healthy",
  "timestamp": "2025-01-15T10:00:00Z",
  "version": "1.0.0",
  "database": "connected",
  "redis": "connected",
  "uptime": "2h45m30s",
  "memory": {
    "used": "45MB",
    "total": "512MB"
  },
  "database_connections": {
    "active": 5,
    "idle": 3,
    "max": 25
  }
}
```

## WebSocket Endpoints (Future)

### Real-time Notifications
Connect to real-time notification stream.

**Endpoint:** `WS /ws/notifications`

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Message Format:**
```json
{
  "type": "maintenance_request_created",
  "data": {
    "id": 1,
    "property_id": 1,
    "description": "Fix leaky faucet",
    "status": "pending",
    "created_at": "2025-01-15T10:00:00Z"
  },
  "timestamp": "2025-01-15T10:00:00Z"
}
```

## Error Handling

### Common Error Responses

#### Validation Error (400)
```json
{
  "success": false,
  "error": [
    {
      "field": "email",
      "message": "must be a valid email address",
      "value": "invalid-email"
    },
    {
      "field": "password",
      "message": "must be at least 8 characters long",
      "value": "short"
    }
  ],
  "message": "Validation failed",
  "timestamp": "2025-01-15T10:00:00Z",
  "request_id": "req_123456789"
}
```

#### Authentication Error (401)
```json
{
  "success": false,
  "error": "Authentication required",
  "message": "Please provide a valid access token",
  "timestamp": "2025-01-15T10:00:00Z",
  "request_id": "req_123456789"
}
```

#### Authorization Error (403)
```json
{
  "success": false,
  "error": "Access denied",
  "message": "Required role: admin",
  "timestamp": "2025-01-15T10:00:00Z",
  "request_id": "req_123456789"
}
```

#### Resource Not Found (404)
```json
{
  "success": false,
  "error": "Resource not found",
  "message": "Property with ID 999 not found",
  "timestamp": "2025-01-15T10:00:00Z",
  "request_id": "req_123456789"
}
```

#### Rate Limit Exceeded (429)
```json
{
  "success": false,
  "error": "Rate limit exceeded",
  "message": "Too many requests. Try again in 60 seconds.",
  "timestamp": "2025-01-15T10:00:00Z",
  "request_id": "req_123456789"
}
```

## Postman Collection

### Import Collection
A Postman collection with all endpoints is available at:
```
https://raw.githubusercontent.com/geoo115/property-manager/main/docs/postman_collection.json
```

### Environment Variables
Set up these environment variables in Postman:
- `base_url`: https://your-domain.com/api/v1
- `access_token`: (set after login)
- `refresh_token`: (set after login)

## SDK Examples

### JavaScript/Node.js
```javascript
const axios = require('axios');

class PropertyManagerAPI {
  constructor(baseURL, accessToken) {
    this.client = axios.create({
      baseURL,
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'Content-Type': 'application/json'
      }
    });
  }

  async getProperties(params = {}) {
    const response = await this.client.get('/landlord/properties', { params });
    return response.data;
  }

  async createProperty(propertyData) {
    const response = await this.client.post('/landlord/properties', propertyData);
    return response.data;
  }

  async createMaintenanceRequest(propertyId, requestData) {
    const response = await this.client.post(`/landlord/properties/${propertyId}/maintenance`, requestData);
    return response.data;
  }
}

// Usage
const api = new PropertyManagerAPI('https://your-domain.com/api/v1', 'your_access_token');
const properties = await api.getProperties({ city: 'New York' });
```

### Python
```python
import requests
from typing import Dict, Optional

class PropertyManagerAPI:
    def __init__(self, base_url: str, access_token: str):
        self.base_url = base_url
        self.session = requests.Session()
        self.session.headers.update({
            'Authorization': f'Bearer {access_token}',
            'Content-Type': 'application/json'
        })
    
    def get_properties(self, params: Optional[Dict] = None) -> Dict:
        response = self.session.get(f'{self.base_url}/landlord/properties', params=params)
        response.raise_for_status()
        return response.json()
    
    def create_property(self, property_data: Dict) -> Dict:
        response = self.session.post(f'{self.base_url}/landlord/properties', json=property_data)
        response.raise_for_status()
        return response.json()
    
    def create_maintenance_request(self, property_id: int, request_data: Dict) -> Dict:
        response = self.session.post(
            f'{self.base_url}/landlord/properties/{property_id}/maintenance',
            json=request_data
        )
        response.raise_for_status()
        return response.json()

# Usage
api = PropertyManagerAPI('https://your-domain.com/api/v1', 'your_access_token')
properties = api.get_properties({'city': 'New York'})
```

## Rate Limiting Details

### Rate Limit Tiers
- **Public endpoints**: 50 requests/minute/IP
- **Authenticated endpoints**: 200 requests/minute/user
- **Authentication endpoints**: 5 requests/minute/IP
- **Admin endpoints**: 500 requests/minute/user

### Rate Limit Headers
All responses include rate limit information:
```
X-RateLimit-Limit: 200
X-RateLimit-Remaining: 199
X-RateLimit-Reset: 1642248000
X-RateLimit-Retry-After: 60
```

### Rate Limit Exceeded Response
```json
{
  "success": false,
  "error": "Rate limit exceeded",
  "message": "Too many requests. Try again in 60 seconds.",
  "timestamp": "2025-01-15T10:00:00Z",
  "request_id": "req_123456789",
  "retry_after": 60
}
```

## Webhooks (Future Enhancement)

### Webhook Configuration
Configure webhooks to receive real-time notifications:

**Endpoint:** `POST /admin/webhooks`

**Request Body:**
```json
{
  "url": "https://your-app.com/webhooks/property-manager",
  "events": ["maintenance_request_created", "lease_expired", "invoice_overdue"],
  "secret": "your_webhook_secret"
}
```

### Webhook Payload
```json
{
  "event": "maintenance_request_created",
  "data": {
    "id": 1,
    "property_id": 1,
    "tenant_id": 456,
    "description": "Fix leaky faucet",
    "status": "pending",
    "created_at": "2025-01-15T10:00:00Z"
  },
  "timestamp": "2025-01-15T10:00:00Z",
  "signature": "sha256=abc123..."
}
```

## Support

For API support, please:
1. Check this documentation first
2. Create an issue on GitHub
3. Contact support at api-support@example.com

## Changelog

### v1.0.0 (2025-01-15)
- Initial API release
- User management endpoints
- Property management endpoints
- Lease management endpoints
- Maintenance request endpoints
- Financial management endpoints
- JWT authentication
- Role-based access control
- Rate limiting
- Comprehensive error handling
