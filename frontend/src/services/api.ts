import axios, { type AxiosInstance } from "axios";

const api: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_URL || "/api",
  headers: {
    "Content-Type": "application/json",
  },
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    // Centralized error handling
    console.error("API Error:", error.response?.data || error.message);
    return Promise.reject(error);
  },
);

export default api;
