import axiosInstance from '../api/axiosInstance';

const API_URL = 'http://localhost:8080/admin/users';

export const getUsers = async () => {
  const response = await axiosInstance.get(API_URL);
  return response.data;
};

export const getUserByID = async (id) => {
  const response = await axiosInstance.get(`${API_URL}/${id}`);
  return response.data;
};

export const createUser = async (userData) => {
  const response = await axiosInstance.post(API_URL, userData);
  return response.data;
};

export const updateUser = async (id, userData) => {
  const response = await axiosInstance.put(`${API_URL}/${id}`, userData);
  return response.data;
};

export const deleteUser = async (id) => {
  const response = await axiosInstance.delete(`${API_URL}/${id}`);
  return response.data;
};

export const activeUser = async () => {
  const response = await axiosInstance.put(`${API_URL}/active`);
  return response.data;
}