import { createContext, useState } from "react";

const AuthContext = createContext({});

export const AuthProvider = ({children}) => {
  // NOTE: lazy initialize "auth" with a Callback
  const [auth, setAuth] = useState(() => JSON.parse(localStorage.getItem("user")) || undefined);
  
  return (
    <AuthContext.Provider value={{auth, setAuth}}>
      {children}      
    </AuthContext.Provider>
  )
}

export default AuthContext;