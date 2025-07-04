# Property Management System (PMS)

## 🏠 Overview
A comprehensive full-stack property management system built with Go backend and React frontend. This system provides complete property, tenant, lease, maintenance, and accounting management with role-based access control.

## ✨ Features

### Core Features
- **User Management**: Admin, landlord, tenant, and maintenance team roles
- **Property Management**: Create, read, update, and delete properties
- **Lease Management**: Handle lease agreements and tenant-property relationships
- **Maintenance Tracking**: Request, assign, and track maintenance tasks
- **Accounting System**: Invoice generation and expense tracking
- **Authentication**: JWT-based authentication with role-based access control
- **Dashboard**: Real-time analytics and KPIs for different user roles
- **Modern UI/UX**: Responsive design with modern components and navigation

### Role-Based Access Control
- **Admin**: Full system access, user management, system-wide reports
- **Landlord**: Property management, tenant relations, financial tracking
- **Tenant**: View personal information, submit maintenance requests
- **Maintenance Team**: View and update maintenance requests

## 🛠 Technology Stack

### Backend
- **Language**: Go 1.24
- **Framework**: Gin (HTTP web framework)
- **Database**: PostgreSQL with GORM ORM
- **Cache**: Redis v8
- **Authentication**: JWT tokens (golang-jwt/jwt/v4)
- **Message Queue**: Apache Kafka (segmentio/kafka-go)
- **Logging**: Sirupsen Logrus with Lumberjack rotation
- **Configuration**: Godotenv for environment variables
- **Security**: Bcrypt for password hashing
- **CORS**: Gin-contrib/cors middleware
- **Testing**: Go testing framework

### Frontend
- **Framework**: React 18
- **Routing**: React Router DOM v7
- **HTTP Client**: Axios
- **UI Components**: Custom modern components
- **State Management**: React Context API
- **Styling**: CSS3 with modern design system
- **Authentication**: JWT with automatic refresh
- **Notifications**: React Toastify
- **Testing**: React Testing Library
- **Prop Validation**: PropTypes

## 📁 Project Structure

```
property-manager/
├── backend/                    # Go backend application
│   ├── api/                   # HTTP handlers organized by domain
│   │   ├── accounting/        # Invoice and expense handlers
│   │   ├── auth/             # Authentication handlers
│   │   ├── dashboard/        # Dashboard and analytics handlers
│   │   ├── lease/            # Lease management handlers
│   │   ├── maintenance/      # Maintenance request handlers
│   │   ├── property/         # Property management handlers
│   │   └── user/             # User management handlers
│   ├── cmd/                  # Application entry points
│   ├── config/               # Configuration management
│   ├── db/                   # Database connections
│   ├── middleware/           # HTTP middleware
│   ├── models/               # Data models
│   ├── router/               # HTTP routing
│   └── tests/                # Test files
├── frontend/                  # React frontend application
│   ├── public/               # Static assets
│   ├── src/
│   │   ├── api/              # API service layer
│   │   ├── components/       # Reusable UI components
│   │   │   ├── common/       # Common components (Button, DataTable, etc.)
│   │   │   └── layout/       # Layout components (Header, Sidebar, etc.)
│   │   ├── constants/        # Application constants and routes
│   │   ├── context/          # React Context providers
│   │   ├── hooks/            # Custom React hooks
│   │   ├── pages/            # Page components
│   │   └── utils/            # Utility functions
│   ├── package.json          # Dependencies and scripts
│   └── README.md             # Frontend-specific documentation
└── docs/                     # Project documentation
```

## 🚀 Installation and Setup

### Prerequisites
- Go 1.24 or later
- Node.js 16 or later
- PostgreSQL 12 or later
- Redis 6 or later
- Apache Kafka (optional, for events)

### Step 1: Clone the Repository
```bash
git clone https://github.com/geoo115/property-manager.git
cd property-manager
```

### Step 2: Backend Setup

#### Install Dependencies
```bash
cd backend
go mod download
```

#### Environment Configuration
Create a `.env` file in the backend directory:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=property_management_db

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here
JWT_EXPIRY=24h
REFRESH_TOKEN_EXPIRY=168h

# Server Configuration
PORT=8080
APP_ENV=development

# Kafka Configuration (optional)
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC_MAINTENANCE=maintenance-requests
KAFKA_TOPIC_INVOICES=invoice-notifications
```

#### Database Setup
1. Create the PostgreSQL database:
```sql
CREATE DATABASE property_management_db;
```

2. Run the database migrations (auto-migrated on startup):
```bash
go run cmd/main.go
```

3. Seed the database with initial data:
```bash
go run scripts/comprehensive_seed.go
```

#### Start the Backend
```bash
go run cmd/main.go
```

The API will be available at `http://localhost:8080`

### Step 3: Frontend Setup

#### Install Dependencies
```bash
cd frontend
npm install
```

#### Environment Configuration (Optional)
Create a `.env` file in the frontend directory:

```env
REACT_APP_API_URL=http://localhost:8080
REACT_APP_ENV=development
```

#### Start the Frontend
```bash
npm start
```

The frontend will be available at `http://localhost:3000`

## 🎨 Frontend Architecture

### Component Structure

#### Layout Components
- **AppLayout**: Main application wrapper with sidebar and header
- **Sidebar**: Navigation menu with role-based menu items
- **Header**: Top navigation with user info and logout functionality

#### Common Components
- **Button**: Reusable button component with multiple variants
- **DataTable**: Advanced table with sorting, filtering, and pagination
- **FormField**: Consistent form input component with validation
- **LoadingSpinner**: Loading states for async operations
- **Alert**: Notification and alert component

#### Page Components
- **Dashboard**: Role-specific dashboard with KPIs and quick actions
- **Properties**: Property management with CRUD operations
- **Leases**: Lease management and tenant relations
- **Users**: User management (admin only)
- **Maintenance**: Maintenance request tracking
- **Invoices**: Invoice management and generation
- **Expenses**: Expense tracking and reporting

### Routing Structure

#### Public Routes
- `/login` - User authentication
- `/register` - User registration

#### Protected Routes (Role-based)
- **Admin Routes**:
  - `/admin/dashboard` - Admin dashboard
  - `/admin/properties` - All properties management
  - `/admin/leases` - All leases management
  - `/admin/users` - User management
  - `/admin/maintenances` - All maintenance requests
  - `/admin/accounting/invoices` - All invoices
  - `/admin/accounting/expenses` - All expenses

- **Landlord Routes**:
  - `/landlord/properties` - Landlord's properties
  - `/landlord/leases` - Landlord's leases
  - `/accounting/invoices` - Landlord's invoices
  - `/accounting/expenses` - Landlord's expenses

- **Tenant Routes**:
  - `/tenant/leases` - Tenant's leases
  - `/tenant/maintenance` - Tenant's maintenance requests

- **Maintenance Team Routes**:
  - `/maintenanceTeam/maintenances` - Assigned maintenance requests

### UI/UX Features

#### Modern Design System
- **Color Palette**: Professional blue and gray theme
- **Typography**: Clean, readable font hierarchy
- **Spacing**: Consistent spacing using CSS custom properties
- **Responsive Design**: Mobile-first approach with breakpoints
- **Accessibility**: WCAG 2.1 compliant components

#### Interactive Features
- **Real-time Updates**: Automatic data refresh for critical information
- **Form Validation**: Client-side validation with error messages
- **Loading States**: Skeleton screens and spinners
- **Toast Notifications**: Success, error, and info messages
- **Confirmation Dialogs**: Destructive action confirmations

#### Navigation
- **Sidebar Navigation**: Collapsible sidebar with role-based menu items
- **Breadcrumbs**: Context-aware navigation breadcrumbs
- **Quick Actions**: Dashboard shortcuts for common tasks
- **Search Functionality**: Global search across entities

## 🔐 Authentication & Security

### Frontend Authentication Flow
1. **Login Process**: User credentials sent to `/login` endpoint
2. **Token Storage**: JWT tokens stored in localStorage
3. **Auto-refresh**: Automatic token refresh before expiration
4. **Route Protection**: Protected routes based on user roles
5. **Logout**: Secure token cleanup and redirect

### Security Features
- **JWT Token Management**: Automatic refresh and secure storage
- **Role-based Route Protection**: Component-level access control
- **Input Validation**: Client-side validation with backend verification
- **CORS Configuration**: Proper cross-origin resource sharing
- **XSS Prevention**: Sanitized user inputs and outputs

## 📊 Dashboard Features

### Admin Dashboard
- **System Overview**: Total users, properties, active leases
- **Financial Metrics**: Monthly revenue, outstanding invoices
- **Maintenance Status**: Open, in-progress, completed requests
- **Recent Activity**: Latest system activities and notifications
- **Quick Actions**: Direct access to common admin tasks

### Landlord Dashboard
- **Property Portfolio**: Property count, occupancy rates
- **Financial Performance**: Rental income, expenses, profit margins
- **Maintenance Overview**: Property maintenance status
- **Tenant Information**: Active tenants, lease expirations
- **Quick Actions**: Add property, create lease, view reports

### Tenant Dashboard
- **Lease Information**: Current lease details, payment history
- **Maintenance Requests**: Submit and track maintenance requests
- **Payment Status**: Outstanding invoices, payment history
- **Property Details**: Current residence information
- **Quick Actions**: Submit maintenance, view lease, contact landlord

### Maintenance Team Dashboard
- **Work Queue**: Assigned maintenance requests by priority
- **Completion Status**: Daily, weekly, monthly completion rates
- **Priority Distribution**: High, medium, low priority requests
- **Performance Metrics**: Response times, completion rates
- **Quick Actions**: Update status, add notes, schedule work

## 📱 Frontend Components

### Button Component
```jsx
// Usage examples
<Button variant="primary" size="medium" onClick={handleClick}>
  Primary Action
</Button>

<Button variant="secondary" size="small" disabled>
  Secondary Action
</Button>

<Button variant="danger" size="large" icon="trash">
  Delete Item
</Button>
```

### DataTable Component
```jsx
// Usage example
<DataTable
  data={properties}
  columns={[
    { key: 'address', title: 'Address', sortable: true },
    { key: 'type', title: 'Type', filterable: true },
    { key: 'rent_amount', title: 'Rent', sortable: true, format: 'currency' }
  ]}
  onEdit={handleEdit}
  onDelete={handleDelete}
  pagination={true}
  searchable={true}
/>
```

### FormField Component
```jsx
// Usage example
<FormField
  label="Property Address"
  type="text"
  value={formData.address}
  onChange={(value) => setFormData({...formData, address: value})}
  required
  validation={addressValidation}
  placeholder="Enter property address"
/>
```

## 🔧 API Integration

### API Service Layer
The frontend uses a centralized API service layer for backend communication:

```javascript
// api/apiService.js
import axios from 'axios';

const apiClient = axios.create({
  baseURL: process.env.REACT_APP_API_URL || 'http://localhost:8080',
  timeout: 10000,
});

// Request interceptor for authentication
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor for error handling
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Handle token expiration
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);
```

### API Endpoints

#### Authentication
- `POST /login` - User authentication
- `POST /register` - User registration
- `POST /refresh-token` - Token refresh
- `POST /logout` - User logout

#### Properties
- `GET /admin/properties` - Get all properties
- `GET /admin/properties/:id` - Get property by ID
- `POST /admin/properties` - Create property
- `PUT /admin/properties/:id` - Update property
- `DELETE /admin/properties/:id` - Delete property

#### Leases
- `GET /admin/leases` - Get all leases
- `GET /admin/leases/:id` - Get lease by ID
- `GET /tenant/leases` - Get tenant leases
- `POST /admin/leases` - Create lease
- `PUT /admin/leases/:id` - Update lease
- `DELETE /admin/leases/:id` - Delete lease

#### Maintenance
- `GET /admin/maintenances` - Get all maintenance requests
- `GET /admin/maintenances/:id` - Get maintenance by ID
- `POST /admin/maintenances` - Create maintenance request
- `PUT /admin/maintenances/:id` - Update maintenance request
- `DELETE /admin/maintenances/:id` - Delete maintenance request

#### Accounting
- `GET /admin/accounting/invoices` - Get all invoices
- `GET /admin/accounting/expenses` - Get all expenses
- `POST /admin/accounting/invoices` - Create invoice
- `POST /admin/accounting/expenses` - Create expense

## 🧪 Testing

### Backend Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test file
go test ./tests/auth_test.go

# Run with verbose output
go test -v ./...
```

### Frontend Testing
```bash
# Run all tests
npm test

# Run tests with coverage
npm test -- --coverage

# Run tests in watch mode
npm test -- --watch

# Run specific test file
npm test -- --testNamePattern="Button"
```

### Test Structure
- **Backend**: `tests/` directory with Go test files
- **Frontend**: `src/` directory with `.test.js` files alongside components

## 📦 Build and Deployment

### Backend Deployment
```bash
# Build binary
go build -o bin/server cmd/main.go

# Docker build
docker build -t property-manager-backend .

# Run with Docker Compose
docker-compose up -d
```

### Frontend Deployment
```bash
# Build for production
npm run build

# The build folder contains optimized production files
# Deploy to your preferred hosting service (Netlify, Vercel, etc.)
```

### Environment Variables

#### Backend Environment Variables
| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | Server port |
| `APP_ENV` | development | Application environment |
| `DB_HOST` | localhost | PostgreSQL host |
| `DB_PORT` | 5432 | PostgreSQL port |
| `DB_USER` | postgres | Database user |
| `DB_PASSWORD` | - | Database password |
| `DB_NAME` | property_management_db | Database name |
| `REDIS_HOST` | localhost | Redis host |
| `REDIS_PORT` | 6379 | Redis port |
| `JWT_SECRET` | - | JWT signing secret |
| `JWT_EXPIRY` | 24h | JWT token expiry |

#### Frontend Environment Variables
| Variable | Default | Description |
|----------|---------|-------------|
| `REACT_APP_API_URL` | http://localhost:8080 | Backend API URL |
| `REACT_APP_ENV` | development | Application environment |

## 👥 Demo Accounts

### Demo Users
The system includes the following demo users for testing:

1. **Admin User**
   - Email: `admin@example.com`
   - Password: `admin123`
   - Role: Admin
   - Access: Full system access

2. **Landlord User**
   - Email: `landlord@example.com`
   - Password: `landlord123`
   - Role: Landlord
   - Access: Property and tenant management

3. **Tenant User**
   - Email: `tenant@example.com`
   - Password: `tenant123`
   - Role: Tenant
   - Access: Personal information and maintenance requests

4. **Maintenance Team User**
   - Email: `maintenance@example.com`
   - Password: `maintenance123`
   - Role: Maintenance Team
   - Access: Maintenance request management

## 🔍 Troubleshooting

### Common Frontend Issues

1. **API Connection Issues**
   - Check backend server is running on port 8080
   - Verify REACT_APP_API_URL environment variable
   - Check browser console for CORS errors

2. **Authentication Issues**
   - Clear localStorage and try logging in again
   - Check JWT token expiration in browser dev tools
   - Verify user credentials with demo accounts

3. **Component Rendering Issues**
   - Check browser console for React errors
   - Verify all required props are passed to components
   - Check for missing dependencies in package.json

### Common Backend Issues

1. **Database Connection Issues**
   - Check PostgreSQL service status
   - Verify connection credentials in .env file
   - Ensure database exists and migrations are run

2. **Authentication Failures**
   - Verify JWT secret configuration
   - Check token expiration settings
   - Validate user credentials in database

### Debug Mode
Enable debug logging:

**Backend:**
```bash
export LOG_LEVEL=debug
go run cmd/main.go
```

**Frontend:**
```bash
REACT_APP_LOG_LEVEL=debug npm start
```

## 📖 Development Guidelines

### Code Style
- **Backend**: Follow Go naming conventions, use `gofmt`
- **Frontend**: Follow React best practices, use ES6+ features
- **Comments**: Add comments for public functions and complex logic
- **Testing**: Write tests for new features and bug fixes

### Git Workflow
1. Create feature branch from main
2. Make changes with descriptive commit messages
3. Add tests for new functionality
4. Create pull request with detailed description
5. Code review and merge

### Component Development
1. Create reusable components in `src/components/common/`
2. Use PropTypes for type checking
3. Follow naming conventions (PascalCase for components)
4. Add CSS modules for styling
5. Include comprehensive documentation

## 📚 API Documentation

### Backend API Documentation
- **OpenAPI/Swagger**: Available at `http://localhost:8080/swagger`
- **Postman Collection**: Import from `docs/postman_collection.json`
- **API Reference**: See `docs/API_DOCUMENTATION.md`

### Frontend Component Documentation
- **Storybook**: Run `npm run storybook` for component documentation
- **PropTypes**: Component prop documentation in source files
- **Usage Examples**: See component files for usage examples

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Setup
1. Follow installation instructions above
2. Make sure all tests pass
3. Follow code style guidelines
4. Add tests for new features
5. Update documentation as needed

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For support and questions:
- **Email**: support@propertymanager.com
- **Documentation**: [API Documentation](docs/API_DOCUMENTATION.md)
- **Issues**: [GitHub Issues](https://github.com/yourusername/property-manager/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/property-manager/discussions)

## 🚧 Roadmap

### Upcoming Features
- [ ] Mobile app (React Native)
- [ ] Advanced reporting and analytics
- [ ] Document management system
- [ ] Payment integration (Stripe/PayPal)
- [ ] Multi-language support
- [ ] Advanced notification system
- [ ] API rate limiting dashboard
- [ ] Audit trail and logging improvements

### Version History

#### Version 1.0.0
- Initial release with full CRUD operations
- Authentication system with JWT
- Role-based access control
- Basic dashboard functionality
- Modern React frontend with responsive design

#### Version 1.1.0 (In Progress)
- Enhanced UI/UX with modern design system
- Advanced data table with filtering and sorting
- Improved form validation and error handling
- Real-time notifications
- Performance optimizations
- Comprehensive testing suite

## 🙏 Acknowledgments

- Go community for excellent libraries and tools
- React community for modern frontend development patterns
- PostgreSQL for robust database management
- Redis for efficient caching solutions
- All contributors and users of this project
