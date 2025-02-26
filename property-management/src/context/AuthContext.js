import React, { createContext, useState, useEffect, useCallback } from 'react';
import axiosInstance from '../api/axiosInstance';
import {jwtDecode} from 'jwt-decode'; // Use default import for jwt-decode
import { setupAxiosInterceptors } from '../api/axiosInstance';

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [token, setToken] = useState(localStorage.getItem('token'));
  const [loading, setLoading] = useState(true);

  const getRefreshTokenFromCookie = useCallback(() => {
    const cookies = document.cookie.split("; ");
    const refreshTokenCookie = cookies.find(cookie => cookie.startsWith("refresh_token="));
    return refreshTokenCookie ? refreshTokenCookie.split("=")[1] : null;
  }, []);

  const logout = useCallback(() => {
    setToken(null);
    localStorage.removeItem('token');
    setUser(null);
  }, []);

  let isRefreshing = false;
let refreshSubscribers = [];

const refreshAccessToken = useCallback(async () => {
  try {
    // Call the refresh endpoint without sending the refresh token in the body.
    // The backend will read the refresh token from the HTTP-only cookie.
    const response = await axiosInstance.post('/refresh-token', {}, { withCredentials: true });
    const newAccessToken = response.data.access_token;
    setToken(newAccessToken);
    localStorage.setItem('token', newAccessToken);
    axiosInstance.defaults.headers.common["Authorization"] = `Bearer ${newAccessToken}`;
    return newAccessToken;
  } catch (error) {
    console.error("Error refreshing token", error);
    logout();
    return null;
  }
}, [logout]);



  // Set up axios interceptors to attach token and refresh when needed
  useEffect(() => {
    setupAxiosInterceptors(logout, refreshAccessToken);
  }, [logout, refreshAccessToken]);

  useEffect(() => {
    if (token) {
      localStorage.setItem('token', token);
      try {
        const decoded = jwtDecode(token);
        if (decoded.exp * 1000 < Date.now()) {
          logout();
        } else {
          // Save the role in localStorage for use in API calls
          localStorage.setItem('role', decoded.role);
          setUser({
            username: decoded.username,
            role: decoded.role,
            user_id: decoded.user_id,
          });
        }
      } catch (error) {
        logout();
      }
    } else {
      localStorage.removeItem('token');
      localStorage.removeItem('role');
      setUser(null);
    }
    setLoading(false);
  }, [token, logout]);

  const login = async (credentials) => {
    try {
      const response = await axiosInstance.post('/login', credentials);
      if (response.data.access_token) {
        setToken(response.data.access_token);
        localStorage.setItem('token', response.data.access_token);
      }
      return response.data;
    } catch (error) {
      throw new Error('Login failed');
    }
  };

  const register = async (userData) => {
    try {
      const response = await axiosInstance.post('/register', userData);
      if (response.data.access_token) {
        setToken(response.data.access_token);
        localStorage.setItem('token', response.data.access_token);
      }
      return response.data;
    } catch (error) {
      throw new Error('Registration failed');
    }
  };

  return (
    <AuthContext.Provider value={{ user, token, login, register, logout }}>
      {loading ? <div>Loading...</div> : children}
    </AuthContext.Provider>
  );
};