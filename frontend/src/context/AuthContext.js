import React, { createContext, useState, useEffect, useCallback, useMemo } from 'react';
import axiosInstance from '../api/axiosInstance';
import { jwtDecode } from 'jwt-decode';
import { setupAxiosInterceptors } from '../api/axiosInstance';
import { USER_ROLES } from '../constants';
import { hasPermission, canAccess, getDashboardConfig, getRoleRoutes } from '../utils/rbac';

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [token, setToken] = useState(localStorage.getItem('token'));
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Clear any authentication errors
  const clearError = useCallback(() => {
    setError(null);
  }, []);

  // Logout function
  const logout = useCallback(() => {
    setToken(null);
    setUser(null);
    localStorage.removeItem('token');
    localStorage.removeItem('role');
    document.cookie = 'refresh_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
    clearError();
  }, [clearError]);

  // Refresh token function
  const refreshAccessToken = useCallback(async () => {
    try {
      const response = await axiosInstance.post('/refresh-token', {}, { 
        withCredentials: true 
      });
      
      // Extract data from the response wrapper
      const responseData = response.data.data || response.data;
      const newAccessToken = responseData.access_token;
      
      setToken(newAccessToken);
      localStorage.setItem('token', newAccessToken);
      axiosInstance.defaults.headers.common["Authorization"] = `Bearer ${newAccessToken}`;
      
      return newAccessToken;
    } catch (error) {
      console.error("Error refreshing token:", error);
      logout();
      return null;
    }
  }, [logout]);

  // Set up axios interceptors
  useEffect(() => {
    setupAxiosInterceptors(logout, refreshAccessToken);
  }, [logout, refreshAccessToken]);

  // Decode token and set user
  const decodeAndSetUser = useCallback((token) => {
    try {
      const decoded = jwtDecode(token);
      
      // Check if token is expired
      if (decoded.exp * 1000 < Date.now()) {
        logout();
        return null;
      }
      
      const userData = {
        id: decoded.user_id,
        username: decoded.username,
        email: decoded.email,
        role: decoded.role,
        firstName: decoded.first_name,
        lastName: decoded.last_name,
        name: `${decoded.first_name || ''} ${decoded.last_name || ''}`.trim() || decoded.username,
      };
      
      localStorage.setItem('role', decoded.role);
      setUser(userData);
      return userData;
    } catch (error) {
      console.error('Error decoding token:', error);
      logout();
      return null;
    }
  }, [logout]);

  // Handle token changes
  useEffect(() => {
    if (token) {
      localStorage.setItem('token', token);
      axiosInstance.defaults.headers.common["Authorization"] = `Bearer ${token}`;
      decodeAndSetUser(token);
    } else {
      localStorage.removeItem('token');
      localStorage.removeItem('role');
      delete axiosInstance.defaults.headers.common["Authorization"];
      setUser(null);
    }
    setLoading(false);
  }, [token, decodeAndSetUser]);

  // Login function
  const login = useCallback(async (credentials) => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await axiosInstance.post('/login', credentials);
      console.log('Full login response:', response);
      console.log('Response data:', response.data);
      
      // Extract data from the response wrapper
      const responseData = response.data.data || response.data;
      console.log('Extracted responseData:', responseData);
      
      if (responseData.access_token) {
        setToken(responseData.access_token);
        const userData = decodeAndSetUser(responseData.access_token);
        return { ...responseData, user: userData };
      }
      
      throw new Error('No access token received');
    } catch (error) {
      console.error('Login error:', error);
      
      let errorMessage = 'Login failed. Please check your connection.';
      
      if (error.response) {
        const { status, data } = error.response;
        
        switch (status) {
          case 401:
            errorMessage = 'Invalid email or password';
            break;
          case 429:
            errorMessage = 'Too many login attempts. Please try again later.';
            break;
          case 500:
            errorMessage = 'Server error. Please try again later.';
            break;
          default:
            errorMessage = data?.message || errorMessage;
        }
      } else if (error.message) {
        errorMessage = error.message;
      }
      
      setError(errorMessage);
      throw new Error(errorMessage);
    } finally {
      setLoading(false);
    }
  }, [decodeAndSetUser]);

  // Register function
  const register = useCallback(async (userData) => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await axiosInstance.post('/register', userData);
      
      // Extract data from the response wrapper
      const responseData = response.data.data || response.data;
      
      if (responseData.access_token) {
        setToken(responseData.access_token);
        const decodedUser = decodeAndSetUser(responseData.access_token);
        return { ...responseData, user: decodedUser };
      }
      
      return responseData;
    } catch (error) {
      console.error('Registration error:', error);
      
      let errorMessage = 'Registration failed. Please check your connection.';
      
      if (error.response) {
        const { status, data } = error.response;
        
        switch (status) {
          case 400:
            errorMessage = 'Invalid user data. Please check your information.';
            break;
          case 409:
            errorMessage = 'Username or email already exists.';
            break;
          case 429:
            errorMessage = 'Too many registration attempts. Please try again later.';
            break;
          case 500:
            errorMessage = 'Server error. Please try again later.';
            break;
          default:
            errorMessage = data?.message || errorMessage;
        }
      }
      
      setError(errorMessage);
      throw new Error(errorMessage);
    } finally {
      setLoading(false);
    }
  }, [decodeAndSetUser]);

  // Check if user has specific role
  const hasRole = useCallback((role) => {
    return user?.role === role;
  }, [user]);

  // Check if user has any of the specified roles
  const hasAnyRole = useCallback((roles) => {
    return roles.some(role => hasRole(role));
  }, [hasRole]);

  // Check if user is admin
  const isAdmin = useCallback(() => {
    return hasRole(USER_ROLES.ADMIN);
  }, [hasRole]);

  // RBAC functions
  const hasUserPermission = useCallback((resource, action) => {
    return user ? hasPermission(user.role, resource, action) : false;
  }, [user]);

  const canUserAccess = useCallback((feature) => {
    return user ? canAccess(user.role, feature) : false;
  }, [user]);

  const getUserDashboardConfig = useCallback(() => {
    return user ? getDashboardConfig(user.role) : {};
  }, [user]);

  const getUserRoutes = useCallback(() => {
    return user ? getRoleRoutes(user.role) : {};
  }, [user]);

  // Get user's full name
  const getUserFullName = useCallback(() => {
    if (!user) return '';
    return user.name || `${user.firstName || ''} ${user.lastName || ''}`.trim() || user.username;
  }, [user]);

  // Memoize the context value to prevent unnecessary re-renders
  const contextValue = useMemo(() => ({
    user,
    token,
    loading,
    error,
    login,
    register,
    logout,
    clearError,
    hasRole,
    hasAnyRole,
    isAdmin,
    hasUserPermission,
    canUserAccess,
    getUserDashboardConfig,
    getUserRoutes,
    getUserFullName,
  }), [
    user,
    token,
    loading,
    error,
    login,
    register,
    logout,
    clearError,
    hasRole,
    hasAnyRole,
    isAdmin,
    hasUserPermission,
    canUserAccess,
    getUserDashboardConfig,
    getUserRoutes,
    getUserFullName,
  ]);

  return (
    <AuthContext.Provider value={contextValue}>
      {children}
    </AuthContext.Provider>
  );
};