export const USER_ROLES = {
  ADMIN: 'admin',
  LANDLORD: 'landlord',
  TENANT: 'tenant',
  MAINTENANCE_TEAM: 'maintenanceTeam',
};

export const ROLE_OPTIONS = [
  { value: USER_ROLES.TENANT, label: 'Tenant' },
  { value: USER_ROLES.LANDLORD, label: 'Landlord' },
  { value: USER_ROLES.MAINTENANCE_TEAM, label: 'Maintenance Team' },
];

export const ROUTES = {
  LOGIN: '/login',
  REGISTER: '/register',
  HOME: '/',
  
  // Admin routes
  ADMIN_DASHBOARD: '/admin/dashboard',
  ADMIN_PROPERTIES: '/admin/properties',
  ADMIN_LEASES: '/admin/leases',
  ADMIN_USERS: '/admin/users',
  ADMIN_MAINTENANCE: '/admin/maintenances',
  ADMIN_INVOICES: '/admin/accounting/invoices',
  ADMIN_EXPENSES: '/admin/accounting/expenses',
  
  // Landlord routes
  LANDLORD_PROPERTIES: '/landlord/properties',
  LANDLORD_LEASES: '/landlord/leases',
  LANDLORD_INVOICES: '/accounting/invoices',
  LANDLORD_EXPENSES: '/accounting/expenses',
  
  // Tenant routes
  TENANT_LEASES: '/tenant/leases',
  TENANT_MAINTENANCE: '/tenant/maintenance',
  
  // Maintenance team routes
  MAINTENANCE_TEAM_MAINTENANCE: '/maintenanceTeam/maintenances',
};

export const ROLE_ROUTES = {
  [USER_ROLES.ADMIN]: ROUTES.ADMIN_DASHBOARD,
  [USER_ROLES.LANDLORD]: ROUTES.LANDLORD_PROPERTIES,
  [USER_ROLES.TENANT]: ROUTES.TENANT_LEASES,
  [USER_ROLES.MAINTENANCE_TEAM]: ROUTES.MAINTENANCE_TEAM_MAINTENANCE,
};

export const API_ENDPOINTS = {
  LOGIN: '/login',
  REGISTER: '/register',
  REFRESH_TOKEN: '/refresh-token',
  LOGOUT: '/logout',
  
  // Dashboard
  DASHBOARD_STATS: '/admin/dashboard/stats',
  DASHBOARD_ACTIVITIES: '/admin/dashboard/activities',
  
  // Users
  USERS: '/admin/users',
  USER_BY_ID: (id) => `/admin/users/${id}`,
  
  // Properties
  PROPERTIES: '/admin/properties',
  PROPERTY_BY_ID: (id) => `/admin/properties/${id}`,
  LANDLORD_PROPERTIES: '/landlord/properties',
  
  // Leases
  LEASES: '/admin/leases',
  LEASE_BY_ID: (id) => `/admin/leases/${id}`,
  LANDLORD_LEASES: '/landlord/leases',
  TENANT_LEASES: '/tenant/leases',
  TENANT_ACTIVE_LEASE: '/tenant/leases/active',
  
  // Maintenance
  MAINTENANCE: '/admin/maintenances',
  MAINTENANCE_BY_ID: (id) => `/admin/maintenances/${id}`,
  TENANT_MAINTENANCE: (leaseId) => `/tenant/leases/${leaseId}/maintenance`,
  TENANT_MAINTENANCE_CREATE: (leaseId) => `/tenant/leases/${leaseId}/maintenance`,
  MAINTENANCE_TEAM_MAINTENANCE: '/maintenanceTeam/maintenances',
  MAINTENANCE_TEAM_MAINTENANCE_BY_ID: (id) => `/maintenanceTeam/maintenance/${id}`,
  
  // Accounting
  INVOICES: '/admin/accounting/invoices',
  INVOICE_BY_ID: (id) => `/admin/accounting/invoices/${id}`,
  EXPENSES: '/admin/accounting/expenses',
  EXPENSE_BY_ID: (id) => `/admin/accounting/expense/${id}`,
  LANDLORD_INVOICES: '/landlord/invoices',
  LANDLORD_EXPENSES: '/landlord/expenses',
  TENANT_INVOICES: '/tenant/invoices',
};
