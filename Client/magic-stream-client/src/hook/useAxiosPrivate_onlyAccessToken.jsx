import axios from "axios";
import useAuth from "./useAuth";
const apiUrl = import.meta.env.VITE_API_BASE_URL;

const useAxiosPrivate = () => {
  // // Unsafe way
  // const {auth, setAuth} = useAuth();
  // const axiosAuth = axios.create({
  //   baseURL: apiUrl,
  // })
  // // Add "interceptor" for Request; to add {Authorization: "Bearer xxxyyy"}
  // axiosAuth.interceptors.request.use((config) => {
  //   if (auth) {
  //     config.headers.Authorization = `Bearer ${auth.token}`;
  //   }
  //   return config;
  // });
  
  // Safe way with Cookies
  const axiosAuth = axios.create({
    baseURL: apiUrl,
    withCredentials: true,
  })
  
  return axiosAuth;
}

export default useAxiosPrivate;