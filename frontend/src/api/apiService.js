// api/apiService.js
import axios from 'axios';

const apiClient = axios.create({
  baseURL: process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1',
  timeout: parseInt(process.env.REACT_APP_API_TIMEOUT) || 10000,
  withCredentials: true,
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
      localStorage.removeItem('role');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Authentication API calls
export const authAPI = {
  login: (credentials) => apiClient.post('/login', credentials),
  register: (userData) => apiClient.post('/register', userData),
  logout: () => apiClient.post('/logout'),
  refreshToken: () => apiClient.post('/refresh-token'),
};

// User API calls
export const userAPI = {
  getUsers: () => apiClient.get('/admin/users'),
  getUserById: (id) => apiClient.get(`/admin/users/${id}`),
  createUser: (userData) => apiClient.post('/admin/users', userData),
  updateUser: (id, userData) => apiClient.put(`/admin/users/${id}`, userData),
  deleteUser: (id) => apiClient.delete(`/admin/users/${id}`),
};

// Property API calls
export const propertyAPI = {
  getProperties: (role) => {
    const endpoint = role === 'admin' ? '/admin/properties' : '/landlord/properties';
    return apiClient.get(endpoint);
  },
  getPropertyById: (id, role) => {
    const endpoint = role === 'admin' ? `/admin/properties/${id}` : `/landlord/properties/${id}`;
    return apiClient.get(endpoint);
  },
  createProperty: (propertyData) => apiClient.post('/admin/properties', propertyData),
  updateProperty: (id, propertyData) => apiClient.put(`/admin/properties/${id}`, propertyData),
  deleteProperty: (id) => apiClient.delete(`/admin/properties/${id}`),
};

// Lease API calls
export const leaseAPI = {
  getLeases: (role, params) => {
    const endpoint = role === 'admin' ? '/admin/leases' : 
                    role === 'landlord' ? '/landlord/leases' : '/tenant/leases';
    return apiClient.get(endpoint, { params });
  },
  getLeaseById: (id, role) => {
    const endpoint = role === 'admin' ? `/admin/leases/${id}` : 
                    role === 'landlord' ? `/landlord/leases/${id}` : `/tenant/leases/${id}`;
    return apiClient.get(endpoint);
  },
  createLease: (leaseData) => apiClient.post('/admin/leases', leaseData),
  updateLease: (id, leaseData) => apiClient.put(`/admin/leases/${id}`, leaseData),
  deleteLease: (id) => apiClient.delete(`/admin/leases/${id}`),
};

// Maintenance API calls
export const maintenanceAPI = {
  getMaintenances: (role, params) => {
    const endpoint = role === 'admin' ? '/admin/maintenances' : 
                    role === 'maintenanceTeam' ? '/maintenanceTeam/maintenances' : null;
    return apiClient.get(endpoint, { params });
  },
  getMaintenanceById: (id, role) => {
    const endpoint = role === 'admin' ? `/admin/maintenances/${id}` : 
                    role === 'maintenanceTeam' ? `/maintenanceTeam/maintenance/${id}` : null;
    return apiClient.get(endpoint);
  },
  createMaintenance: (maintenanceData, role, resourceId) => {
    let endpoint;
    if (role === 'tenant') {
      endpoint = `/tenant/leases/${resourceId}/maintenance`;
    } else if (role === 'landlord') {
      endpoint = `/landlord/properties/${resourceId}/maintenances`;
    } else {
      endpoint = '/admin/maintenances';
    }
    return apiClient.post(endpoint, maintenanceData);
  },
  updateMaintenance: (id, maintenanceData, role) => {
    const endpoint = role === 'admin' ? `/admin/maintenances/${id}` : 
                    role === 'maintenanceTeam' ? `/maintenanceTeam/maintenance/${id}` : null;
    return apiClient.put(endpoint, maintenanceData);
  },
  deleteMaintenance: (id) => apiClient.delete(`/admin/maintenances/${id}`),
};

// Accounting API calls
export const accountingAPI = {
  getInvoices: (role) => {
    const endpoint = role === 'admin' ? '/admin/accounting/invoices' : 
                    role === 'landlord' ? '/landlord/invoices' : '/tenant/invoices';
    return apiClient.get(endpoint);
  },
  getInvoiceById: (id) => apiClient.get(`/admin/accounting/invoices/${id}`),
  createInvoice: (invoiceData) => apiClient.post('/admin/accounting/invoices', invoiceData),
  updateInvoice: (id, invoiceData) => apiClient.put(`/admin/accounting/invoices/${id}`, invoiceData),
  deleteInvoice: (id) => apiClient.delete(`/admin/accounting/invoices/${id}`),
  
  getExpenses: (role) => {
    const endpoint = role === 'admin' ? '/admin/accounting/expenses' : '/landlord/expenses';
    return apiClient.get(endpoint);
  },
  getExpenseById: (id) => apiClient.get(`/admin/accounting/expense/${id}`),
  createExpense: (expenseData) => apiClient.post('/admin/accounting/expense', expenseData),
  updateExpense: (id, expenseData) => apiClient.put(`/admin/accounting/expense/${id}`, expenseData),
  deleteExpense: (id) => apiClient.delete(`/admin/accounting/expense/${id}`),
};

// Dashboard API calls
export const dashboardAPI = {
  getStats: () => apiClient.get('/admin/dashboard/stats'),
  getActivities: () => apiClient.get('/admin/dashboard/activities'),
};

export default apiClient;
