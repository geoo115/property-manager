import axiosInstance from '../api/axiosInstance';

export const getUsers = async () => {
  const response = await axiosInstance.get('/admin/users');
  return response.data;
};

export const getUserByID = async (id) => {
  const response = await axiosInstance.get(`/admin/users/${id}`);
  return response.data;
};

export const createUser = async (userData) => {
  const response = await axiosInstance.post('/admin/users', userData);
  return response.data;
};

export const updateUser = async (id, userData) => {
  const response = await axiosInstance.put(`/admin/users/${id}`, userData);
  return response.data;
};

export const deleteUser = async (id) => {
  const response = await axiosInstance.delete(`/admin/users/${id}`);
  return response.data;
};

export const activeUser = async () => {
  const response = await axiosInstance.put('/admin/users/active');
  return response.data;
};