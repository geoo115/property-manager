import axiosInstance from './axiosInstance';

export const login = async (credentials) => {
  const response = await axiosInstance.post('/login', credentials);
  console.log('Login response:', response.data);
  
  // Handle the response structure that wraps data in a 'data' property
  const responseData = response.data;
  const accessToken = responseData.data?.access_token;
  
  if (accessToken) {
    localStorage.setItem('token', accessToken);
    console.log('Token stored');
  }
  
  return responseData;
};

export const register = async (userData) => {
  const response = await axiosInstance.post('/register', userData);
  return response.data;
};

export const logout = () => {
  localStorage.removeItem('token');
};