import { createContext, useEffect, useState } from "react";

const AuthContext = createContext({});

export const AuthProvider = ({children}) => {
  const [auth, setAuth] = useState();
  const [loading, setLoading] = useState(true); // did we finish authentication checking?
  
  // Scanning for "user" in localStorage first
  useEffect(() => {
    try {
      const storedUser = localStorage.getItem("user");
      if (storedUser) {
        const parsedUser = JSON.parse(storedUser);
        setAuth(parsedUser);
      } 
    } catch(error) {
      console.error("Failed to parse user from localStorage")
    } finally {
      setLoading(false);
    }
  }, [])
  
  // Remove or add "user" from localStorage based on "auth" state
  useEffect(() => {
    if (auth) {
      localStorage.setItem("user", JSON.stringify(auth));
    } else {
      localStorage.removeItem("user")
    }
  }, [auth])
  
  return (
    <AuthContext.Provider value={{auth, setAuth, loading}}>
      {children}      
    </AuthContext.Provider>
  )
}

export default AuthContext;