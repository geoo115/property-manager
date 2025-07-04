# Property Management System - Frontend

[![React](https://img.shields.io/badge/React-18.2.0-blue.svg)](https://reactjs.org)
[![Node.js](https://img.shields.io/badge/Node.js-16+-green.svg)](https://nodejs.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

## Overview

This is the React frontend for the Property Management System (PMS). It provides a modern, responsive web interface for managing properties, leases, maintenance requests, and financial operations with role-based access control.

## Quick Start

### Prerequisites

- **Node.js 16+** - [Download](https://nodejs.org/)
- **npm or yarn** - Package manager
- **Running Backend** - PMS backend API server at `http://localhost:8080`

### Installation

1. **Navigate to frontend directory**
   ```bash
   cd property-manager/frontend
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Configure environment**
   ```bash
   # Create .env file with API configuration
   echo "REACT_APP_API_URL=http://localhost:8080/api/v1" > .env
   ```

4. **Start development server**
   ```bash
   npm start
   ```

The application will open at `http://localhost:3000`

## Available Scripts

In the project directory, you can run:

### `npm start`

Runs the app in development mode with hot reload.\
Open [http://localhost:3000](http://localhost:3000) to view it in your browser.

**Note**: Ensure the backend API is running at `http://localhost:8080` before starting the frontend.

### `npm test`

Launches the test runner in interactive watch mode.\
See the section about [running tests](https://facebook.github.io/create-react-app/docs/running-tests) for more information.

### `npm run build`

Builds the app for production to the `build` folder.\
The build is minified and optimized for deployment.

## Features

### 🔐 **Role-Based Interface**
- **Admin Dashboard**: Full system management capabilities
- **Landlord Portal**: Property and lease management  
- **Tenant Portal**: Lease viewing and maintenance requests
- **Maintenance Team**: Request management and updates

### 🎨 **User Experience**
- **Responsive Design**: Works on desktop, tablet, and mobile
- **Modern UI**: Clean, intuitive interface with React components
- **Real-time Updates**: Dynamic data updates with API integration
- **Error Handling**: Comprehensive error messages and validation

### 🛡️ **Security Features**
- **JWT Authentication**: Secure token-based authentication
- **Protected Routes**: Role-based route protection
- **Automatic Token Refresh**: Seamless token renewal
- **CORS Support**: Proper cross-origin request handling

## Configuration

### Environment Variables

Create a `.env` file in the frontend directory:

```env
# API Configuration
REACT_APP_API_URL=http://localhost:8080/api/v1
REACT_APP_NAME="Property Management System"

# Development Settings
GENERATE_SOURCEMAP=true
```

### Production Configuration

For production deployment:

```env
# Production API URL
REACT_APP_API_URL=https://api.yourdomain.com/api/v1
REACT_APP_NAME="Property Management System"

# Disable source maps in production
GENERATE_SOURCEMAP=false
```

## Authentication

The frontend implements JWT-based authentication with:

- **Login/Register Pages**: User authentication forms
- **Token Management**: Automatic token refresh
- **Protected Routes**: Role-based access control
- **Session Handling**: Secure token storage

### User Roles

- **Admin**: Full system access
- **Landlord**: Property and lease management
- **Tenant**: View leases and submit maintenance requests
- **Maintenance Team**: Manage maintenance requests

## API Integration

The frontend communicates with the backend API using:

- **Axios HTTP Client**: Configured with interceptors
- **Base URL**: `http://localhost:8080/api/v1`
- **Authentication**: Bearer token in headers
- **Error Handling**: Automatic retry and error recovery

### API Endpoints Used

- `POST /login` - User authentication
- `POST /register` - User registration
- `GET /admin/*` - Admin operations
- `GET /landlord/*` - Landlord operations
- `GET /tenant/*` - Tenant operations
- `GET /maintenanceTeam/*` - Maintenance operations

## Development

### Project Structure

```
frontend/
├── public/           # Static assets
├── src/
│   ├── api/         # API integration
│   ├── components/  # Reusable components
│   ├── context/     # React contexts
│   ├── pages/       # Page components
│   └── styles/      # CSS styles
├── package.json     # Dependencies
└── .env            # Environment config
```

### Common Development Tasks

1. **Add New Component**
   ```bash
   # Create component file
   touch src/components/NewComponent.js
   
   # Add to exports if needed
   ```

2. **Add New Page**
   ```bash
   # Create page component
   touch src/pages/NewPage.js
   
   # Add route in App.js
   ```

3. **Update API Integration**
   ```bash
   # Modify API files in src/api/
   # Update axiosInstance configuration
   ```

## Deployment

### Build for Production

```bash
npm run build
```

### Deploy to Static Hosting

The built files in the `build/` directory can be deployed to:

- **Netlify**: Drag and drop the build folder
- **Vercel**: Connect GitHub repository
- **AWS S3**: Upload build files to S3 bucket
- **GitHub Pages**: Use `gh-pages` package

## Troubleshooting

### Common Issues

1. **API Connection Errors**
   - Ensure backend is running on port 8080
   - Check CORS configuration
   - Verify API URL in .env file

2. **Authentication Issues**
   - Clear browser localStorage
   - Check token expiration
   - Verify JWT secret configuration

3. **Route Protection Problems**
   - Check user role in localStorage
   - Verify ProtectedRoute component
   - Ensure proper route configuration

### Debug Tools

- **React DevTools**: Browser extension for component inspection
- **Network Tab**: Check API requests and responses
- **Console Logs**: Review error messages and warnings

## Available Scripts (Continued)

### `npm run build`

Builds the app for production to the `build` folder.\
It correctly bundles React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.\
Your app is ready to be deployed!

See the section about [deployment](https://facebook.github.io/create-react-app/docs/deployment) for more information.

### `npm run eject`

**Note: this is a one-way operation. Once you `eject`, you can't go back!**

If you aren't satisfied with the build tool and configuration choices, you can `eject` at any time. This command will remove the single build dependency from your project.

Instead, it will copy all the configuration files and the transitive dependencies (webpack, Babel, ESLint, etc) right into your project so you have full control over them. All of the commands except `eject` will still work, but they will point to the copied scripts so you can tweak them. At this point you're on your own.

You don't have to ever use `eject`. The curated feature set is suitable for small and middle deployments, and you shouldn't feel obligated to use this feature. However we understand that this tool wouldn't be useful if you couldn't customize it when you are ready for it.

## Learn More

You can learn more in the [Create React App documentation](https://facebook.github.io/create-react-app/docs/getting-started).

To learn React, check out the [React documentation](https://reactjs.org/).

### Code Splitting

This section has moved here: [https://facebook.github.io/create-react-app/docs/code-splitting](https://facebook.github.io/create-react-app/docs/code-splitting)

### Analyzing the Bundle Size

This section has moved here: [https://facebook.github.io/create-react-app/docs/analyzing-the-bundle-size](https://facebook.github.io/create-react-app/docs/analyzing-the-bundle-size)

### Making a Progressive Web App

This section has moved here: [https://facebook.github.io/create-react-app/docs/making-a-progressive-web-app](https://facebook.github.io/create-react-app/docs/making-a-progressive-web-app)

### Advanced Configuration

This section has moved here: [https://facebook.github.io/create-react-app/docs/advanced-configuration](https://facebook.github.io/create-react-app/docs/advanced-configuration)

### Deployment

This section has moved here: [https://facebook.github.io/create-react-app/docs/deployment](https://facebook.github.io/create-react-app/docs/deployment)

### `npm run build` fails to minify

This section has moved here: [https://facebook.github.io/create-react-app/docs/troubleshooting#npm-run-build-fails-to-minify](https://facebook.github.io/create-react-app/docs/troubleshooting#npm-run-build-fails-to-minify)
