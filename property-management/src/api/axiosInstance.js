//api/axiosInstance.js
import axios from "axios";

const axiosInstance = axios.create({
  baseURL: process.env.REACT_APP_API_URL || "http://localhost:8080",
  withCredentials: true,
});


export const setupAxiosInterceptors = (logout, refreshAccessToken) => {
  // Request interceptor
  axiosInstance.interceptors.request.use((config) => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  });

  axiosInstance.interceptors.response.use(
    response => response,
    async (error) => {
      const originalRequest = error.config;
  
      if (error.response?.status === 429) {
        console.warn("Rate limited. Retrying after delay...");
        await new Promise(res => setTimeout(res, 3000)); // Wait 3 seconds before retrying
        return axiosInstance(originalRequest);
      }
  
      if (error.response?.status === 401 && !originalRequest._retry) {
        originalRequest._retry = true;
        try {
          await refreshAccessToken();
          return axiosInstance(originalRequest);
        } catch (refreshError) {
          logout();
          return Promise.reject(refreshError);
        }
      }
  
      return Promise.reject(error);
    }
  );
  
};

export default axiosInstance;