// This is an improved version of "useAxiosPrivate" - to use Refresh Token to refresh the Access Token if expired
import axios from "axios"
import useAuth from "./useAuth";
import { useEffect } from "react";

const apiUrl = import.meta.env.VITE_API_BASE_URL;

const useAxiosPrivate = () => {
  const axiosAuth = axios.create({
    baseURL: apiUrl,
    withCredentials: true, // to use Cookies
  })
  
  // Handle Refresh Token
  const {auth, setAuth} = useAuth();
  
  let isRefreshing = false
  let failedQueue = []
  
  const processQueue = (error, response=null) => {
    failedQueue.forEach(prom => {
      if (error) {
        prom.reject(error);
      } else {
        prom.resolve(response);
      }
    });
    failedQueue = [];
  }
  
  useEffect(() => {
    // Interceptor of Response from Server
    axiosAuth.interceptors.response.use(
      response => response,
      // Error handler; coming from "/refresh" Endpoint or other Endpoints
      async error => {
        console.log("Interceptor caught error: ", error);
        const originalRequest = error.config;
        
        // If Error comes from "/refresh" endpoint => Refresh Token expires or invalid => no retry
        if (originalRequest.url.includes("/refresh") && error.response.status === 401) {
          console.error("Refresh token has expired or is invalid.");
          return Promise.reject(error);
        }
        
        // If Error comes from other endpoints e.g. when Access Token has
        if (error.response && error.response.status === 401 && !originalRequest._retry) {
          if (isRefreshing) {
            return new Promise((resolve, reject) => {
              failedQueue.push({resolve, reject});
            })
            .then(() => axiosAuth(originalRequest))
            .catch(err => Promise.reject(err));
          }
          // When Access Token expires => Retry by POST to /refresh endpoint to get new Access Token
          originalRequest._retry = true;
          isRefreshing = true;
          return new Promise((resolve, reject) => {
            // send to /refresh Endpoint
            axiosAuth(originalRequest)
            .post("/refresh")
            .then(() => {
              processQueue(null) // no error 
              
              axiosAuth(originalRequest)
              .then(resolve)
              .catch(reject);
            })
            // error coming from /refresh endpoint e.g. refresh token has expired
            .catch(refreshError => {
              processQueue(refreshError, null);
              localStorage.removeItem("user");
              setAuth(null); // clear auth state
              reject(refreshError);
            })
            .finally(() => {
              isRefreshing = false;
            });
          });
        }
        return Promise.reject(error);
      }
    );
  }, [auth])
  
  return axiosAuth
}

export default useAxiosPrivate;