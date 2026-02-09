import axios from "axios";

// For dev: points directly to backend (localhost:8080)
// For production: uses /api prefix which gets rewritten by Next.js
const isProduction = process.env.NODE_ENV === "production";
const baseURL = isProduction ? "/api" : (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080");

export const apiClient = axios.create({
  baseURL,
  headers: {
    "Content-Type": "application/json",
  },
});

// Add response interceptor for error handling
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    // Handle 401 - unauthorized (token expired or invalid)
    if (error.response?.status === 401) {
      // Clear token cookie
      if (typeof window !== "undefined") {
        document.cookie = "token=;expires=Thu, 01 Jan 1970 00:00:00 UTC;path=/;";
        // Redirect to login if not already there
        if (window.location.pathname !== "/login") {
          window.location.href = "/login";
        }
      }
    }
    return Promise.reject(error);
  }
);
