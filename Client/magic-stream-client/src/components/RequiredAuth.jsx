import { useLocation, Navigate, Outlet} from "react-router-dom";
import useAuth from "../hook/useAuth";
import Spinner from "./spinner/Spinner";

const RequiredAuth = () => {
  const { auth, loading } = useAuth();
  const location = useLocation();
  
  if (loading) {
    return (
      <Spinner />
    );
  }
  
  return auth ? (
    <Outlet />
  ): (
    // Before navigating, we store current Page in {from: location}
    <Navigate to="/login" state={{from:location}} replace/>
  )
};

export default RequiredAuth;