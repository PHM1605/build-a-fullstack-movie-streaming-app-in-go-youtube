import { Routes, Route, useNavigate } from 'react-router-dom'
import './App.css'
import Header from './components/header/Header'
import Home from './components/home/Home'
import Login from './components/login/Login'
import Register from "./components/register/Register"
import Layout from './components/Layout'
import RequiredAuth from './components/RequiredAuth'
import Recommended from './components/recommended/Recommended'
import Review from './components/review/Review'
import axiosClient from "./api/axiosConfig"
import useAuth from './hook/useAuth'
import StreamMovie from './components/stream/StreamMovie'

function App() {
  const navigate = useNavigate();
  const {auth, setAuth} = useAuth();
  
  const updateMovieReview = (imdb_id) => {
    navigate(`/review/${imdb_id}`);
  };
  
  const handleLogout = async () => {
    try {
      const response = await axiosClient.post("/logout", {user_id: auth.user_id });
      setAuth(null); // this will trigger the "useEffect()" in "AuthProvider.jsx" to remove "user" in "localStorage" 
      console.log("User logged out");
    } catch(error) {
      console.error("Error logging out: ", error);
    }
  }
  
  return (
    <>
    <Header handleLogout={handleLogout}/>
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<Home updateMovieReview={updateMovieReview}/>} />
        <Route path="/register" element={<Register />} />
        <Route path="/login" element={<Login />} />
        {/* Protected routes  */}
        <Route element={<RequiredAuth />}>
          <Route path="/recommended" element={<Recommended />} />
          <Route path="/review/:imdb_id" element={<Review />} />
          <Route path="/stream/:yt_id" element={<StreamMovie />}></Route>
        </Route>
      </Route>
    </Routes>
    </>
  )
}

export default App
