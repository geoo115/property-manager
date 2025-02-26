import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

export const login = async (credentials) => {
  const response = await axios.post(`${API_URL}/login`, credentials);
  console.log('Login response:', response.data);
  if (response.data.access_token) {
    localStorage.setItem('token', response.data.access_token);
    console.log('Token stored');
  }
  return response.data;
};
export const register = async (userData) => {
  const response = await axios.post(`${API_URL}/register`, userData);
  return response.data;
};

export const logout = () => {
  localStorage.removeItem('token');
};