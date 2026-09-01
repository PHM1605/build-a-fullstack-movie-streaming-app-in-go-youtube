import { useState } from "react";
import Container from "react-bootstrap/Container";
import Button from "react-bootstrap/Button";
import Form from "react-bootstrap/Form";
import axiosClient from "../../api/axiosConfig";
import { useNavigate, Link, useLocation } from "react-router-dom";
import useAuth from "../../hook/useAuth";
import logo from "../../assets/MagicStreamLogo.png";

const Login = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(false);
  const location = useLocation(); // NOTE: give information that's given in the URL
  const {setAuth} = useAuth(); // NOTE: before sending Login Request we set our Authorization
  const navigate = useNavigate();
  
  // Check RequiredAuth.jsx to see where "from" is set
  const from = location.state?.from?.pathname || "/";
  
  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    
    try {
      const response = await axiosClient.post("/login", {email, password});
      console.log(response.data)
      if (response.data.error) {
        setError(response.data.error);
        return;
      }
      setAuth(response.data); // this will trigger the "useEffect()" in "AuthProvider.jsx" to insert "user" to "localStorage" 
      // handle successful login; return User to where he was previously
      navigate(from, {replace: true})
      // navigate("/");
    } catch(err) {
      console.error(err);
      setError("Invalid email or password");
    } finally {
      setLoading(false);
    }
  }
  
  return (
    <Container className="login-container d-flex align-items-center justify-content-center min-vh-100">
      <div className="login-card shadow p-4 rounded bg-white" style={{maxWidth: 400, width:"100%"}}>
        <div className="text-center mb-4">
          <img src={logo} alt="Logo" width={60} className="mb-2" />
          <h2 className="fw-bold">Sign In</h2>
          <p className="text-muted">Welcome back! Please login to your account.</p>
        </div>
        
        {error && <div className="alert alert-danger py-2">{error}</div>}
        
        <Form onSubmit={handleSubmit}>
          <Form.Group controlId="formBasicEmail" className="mb-3">
            <Form.Label>Email address</Form.Label>
            <Form.Control 
              type="email"
              placeholder="Enter email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              autoFocus 
            />
          </Form.Group>
          
          <Form.Group controlId="formBasicPassword" className="mb-3">
            <Form.Label>Password</Form.Label>
            <Form.Control 
              type="password" 
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </Form.Group>
          
          <Button
            variant="primary"
            type="submit"
            className="w-100 mb-2"
            disabled={loading}
            style={{fontWeight: 600, letterSpacing: 1}}
          >
            {loading ? (
              <>
              <span className="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true" />
              Logging in...
              </>
            ): "Login"
            }
          </Button>
        </Form>
        
        <div className="text-center mt-3">
          <span className="text-muted">Don't have an account?</span>
          <Link to="/register" className="fw-semibold">Register here</Link>
        </div>      
      </div>
    </Container>
  )
}

export default Login;