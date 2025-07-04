// utils/rbac.js
import { USER_ROLES } from '../constants';

/**
 * Role-Based Access Control (RBAC) utility functions
 */

// Define permissions for each role
export const ROLE_PERMISSIONS = {
  [USER_ROLES.ADMIN]: {
    // Full system access - Admin: Full system access, user management, system-wide reports
    users: ['create', 'read', 'update', 'delete'],
    properties: ['create', 'read', 'update', 'delete'],
    leases: ['create', 'read', 'update', 'delete'],
    maintenance: ['create', 'read', 'update', 'delete'],
    invoices: ['create', 'read', 'update', 'delete'],
    expenses: ['create', 'read', 'update', 'delete'],
    reports: ['read', 'export', 'system-wide'],
    dashboard: ['read', 'system-wide'],
    settings: ['read', 'update'],
    system: ['read', 'update', 'manage'],
  },
  [USER_ROLES.LANDLORD]: {
    // Landlord: Property management, tenant relations, financial tracking
    properties: ['read', 'update'], // Only their properties
    leases: ['create', 'read', 'update'], // Only their properties' leases
    maintenance: ['read', 'update'], // Only their properties' maintenance
    invoices: ['create', 'read', 'update'], // Only their properties' invoices
    expenses: ['create', 'read', 'update'], // Only their expenses
    reports: ['read'], // Only their data
    dashboard: ['read', 'landlord-view'],
    tenants: ['read', 'communicate'], // Tenant relations
    financial: ['read', 'track'], // Financial tracking
  },
  [USER_ROLES.TENANT]: {
    // Tenant: View personal information, submit maintenance requests
    profile: ['read', 'update'],
    leases: ['read'], // Only their leases
    maintenance: ['create', 'read'], // Only their maintenance requests
    invoices: ['read'], // Only their invoices
    dashboard: ['read', 'tenant-view'],
    payments: ['read'], // View payment history
    personal: ['read', 'update'], // Personal information
  },
  [USER_ROLES.MAINTENANCE_TEAM]: {
    // Maintenance Team: View and update maintenance requests
    maintenance: ['read', 'update'], // Assigned maintenance requests
    properties: ['read'], // To understand maintenance context
    dashboard: ['read', 'maintenance-view'],
    reports: ['read'], // Maintenance reports
    tasks: ['read', 'update'], // Maintenance tasks
  },
};

/**
 * Check if a user has a specific permission
 * @param {string} userRole - The user's role
 * @param {string} resource - The resource (e.g., 'properties', 'users')
 * @param {string} action - The action (e.g., 'create', 'read', 'update', 'delete')
 * @returns {boolean} - Whether the user has the permission
 */
export const hasPermission = (userRole, resource, action) => {
  const rolePermissions = ROLE_PERMISSIONS[userRole];
  if (!rolePermissions) return false;
  
  const resourcePermissions = rolePermissions[resource];
  if (!resourcePermissions) return false;
  
  return resourcePermissions.includes(action);
};

/**
 * Check if a user can access a specific feature
 * @param {string} userRole - The user's role
 * @param {string} feature - The feature name
 * @returns {boolean} - Whether the user can access the feature
 */
export const canAccess = (userRole, feature) => {
  const rolePermissions = ROLE_PERMISSIONS[userRole];
  return rolePermissions && rolePermissions[feature];
};

/**
 * Get allowed actions for a user on a specific resource
 * @param {string} userRole - The user's role
 * @param {string} resource - The resource
 * @returns {string[]} - Array of allowed actions
 */
export const getAllowedActions = (userRole, resource) => {
  const rolePermissions = ROLE_PERMISSIONS[userRole];
  if (!rolePermissions) return [];
  
  return rolePermissions[resource] || [];
};

/**
 * Filter navigation items based on user role
 * @param {Array} navigationItems - Array of navigation items
 * @param {string} userRole - The user's role
 * @returns {Array} - Filtered navigation items
 */
export const filterNavigation = (navigationItems, userRole) => {
  return navigationItems.filter(item => {
    if (item.roles && item.roles.length > 0) {
      return item.roles.includes(userRole);
    }
    return true;
  });
};

/**
 * Get dashboard configuration based on user role
 * @param {string} userRole - The user's role
 * @returns {Object} - Dashboard configuration
 */
export const getDashboardConfig = (userRole) => {
  const configs = {
    [USER_ROLES.ADMIN]: {
      showSystemStats: true,
      showAllProperties: true,
      showAllUsers: true,
      showFinancialOverview: true,
      showSystemReports: true,
      allowUserManagement: true,
      allowSystemSettings: true,
    },
    [USER_ROLES.LANDLORD]: {
      showPropertyStats: true,
      showTenantInfo: true,
      showFinancialTracking: true,
      showMaintenanceOverview: true,
      showPropertyReports: true,
      allowPropertyManagement: true,
      allowTenantCommunication: true,
    },
    [USER_ROLES.TENANT]: {
      showPersonalInfo: true,
      showLeaseDetails: true,
      showMaintenanceRequests: true,
      showPaymentHistory: true,
      allowMaintenanceSubmission: true,
      allowProfileUpdate: true,
    },
    [USER_ROLES.MAINTENANCE_TEAM]: {
      showMaintenanceQueue: true,
      showAssignedTasks: true,
      showCompletionStats: true,
      showMaintenanceReports: true,
      allowMaintenanceUpdate: true,
      allowTaskManagement: true,
    },
  };
  
  return configs[userRole] || {};
};

/**
 * Check if user role is valid
 * @param {string} role - The role to validate
 * @returns {boolean} - Whether the role is valid
 */
export const isValidRole = (role) => {
  return Object.values(USER_ROLES).includes(role);
};

/**
 * Get user role display name
 * @param {string} role - The role
 * @returns {string} - Display name for the role
 */
export const getRoleDisplayName = (role) => {
  const displayNames = {
    [USER_ROLES.ADMIN]: 'Administrator',
    [USER_ROLES.LANDLORD]: 'Landlord',
    [USER_ROLES.TENANT]: 'Tenant',
    [USER_ROLES.MAINTENANCE_TEAM]: 'Maintenance Team',
  };
  
  return displayNames[role] || role;
};

/**
 * Get role-specific routes
 * @param {string} userRole - The user's role
 * @returns {Object} - Role-specific routes
 */
export const getRoleRoutes = (userRole) => {
  const routes = {
    [USER_ROLES.ADMIN]: {
      dashboard: '/dashboard',
      properties: '/admin/properties',
      leases: '/admin/leases',
      users: '/admin/users',
      maintenance: '/admin/maintenances',
      invoices: '/admin/accounting/invoices',
      expenses: '/admin/accounting/expenses',
    },
    [USER_ROLES.LANDLORD]: {
      dashboard: '/dashboard',
      properties: '/landlord/properties',
      leases: '/landlord/leases',
      maintenance: '/landlord/maintenance',
      invoices: '/accounting/invoices',
      expenses: '/accounting/expenses',
    },
    [USER_ROLES.TENANT]: {
      dashboard: '/dashboard',
      profile: '/tenant/profile',
      leases: '/tenant/leases',
      maintenance: '/tenant/maintenance',
      invoices: '/tenant/invoices',
      payments: '/tenant/payments',
    },
    [USER_ROLES.MAINTENANCE_TEAM]: {
      dashboard: '/dashboard',
      maintenance: '/maintenanceTeam/maintenances',
      reports: '/maintenanceTeam/reports',
    },
  };
  
  return routes[userRole] || {};
};

export default {
  hasPermission,
  canAccess,
  getAllowedActions,
  filterNavigation,
  getDashboardConfig,
  isValidRole,
  getRoleDisplayName,
  getRoleRoutes,
  ROLE_PERMISSIONS,
};
