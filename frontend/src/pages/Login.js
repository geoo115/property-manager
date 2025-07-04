import React, { useState, useContext } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { toast } from 'react-toastify';
import { AuthContext } from '../context/AuthContext';
import FormInput from '../components/forms/FormInput';
import Button from '../components/common/Button';
import Alert from '../components/common/Alert';
import Card from '../components/common/Card';

const Login = () => {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
  });
  const [errors, setErrors] = useState({});
  const [isLoading, setIsLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);

  const { login } = useContext(AuthContext);
  const navigate = useNavigate();
  const location = useLocation();

  const from = location.state?.from?.pathname || '/dashboard';

  const validateForm = () => {
    const newErrors = {};

    if (!formData.email) {
      newErrors.email = 'Email is required';
    } else if (!/\S+@\S+\.\S+/.test(formData.email)) {
      newErrors.email = 'Please enter a valid email address';
    }

    if (!formData.password) {
      newErrors.password = 'Password is required';
    } else if (formData.password.length < 6) {
      newErrors.password = 'Password must be at least 6 characters';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
    
    // Clear error when user starts typing
    if (errors[name]) {
      setErrors(prev => ({
        ...prev,
        [name]: ''
      }));
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    setIsLoading(true);
    
    try {
      // Send credentials as an object
      const credentials = {
        email: formData.email,
        password: formData.password
      };
      
      await login(credentials);
      toast.success('Login successful! Welcome back.');
      navigate(from, { replace: true });
    } catch (error) {
      const message = error.message || 'Login failed. Please try again.';
      toast.error(message);
      setErrors({ general: message });
    } finally {
      setIsLoading(false);
    }
  };

  const handleDemoLogin = async (role) => {
    setIsLoading(true);
    
    try {
      // Demo credentials based on role
      const demoCredentials = {
        admin: { email: 'admin@example.com', password: 'Admin123!' },
        landlord: { email: 'landlord@example.com', password: 'Landlord123!' },
        tenant: { email: 'tenant@example.com', password: 'Tenant123!' },
        maintenance: { email: 'maintenance@example.com', password: 'Maintenance123!' }
      };

      const credentials = demoCredentials[role];
      await login(credentials);
      toast.success(`Demo login successful! Logged in as ${role}.`);
      navigate('/dashboard', { replace: true });
    } catch (error) {
      toast.error('Demo login failed. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="login-page">
      <div className="login-container">
        <Card 
          className="login-card"
          elevation="lg"
        >
          <div className="login-header">
            <div className="login-logo">
              <span role="img" aria-label="Property Manager" className="logo-icon">
                🏢
              </span>
              <h1 className="login-title">Property Manager</h1>
            </div>
            <p className="login-subtitle">Sign in to your account</p>
          </div>

          {errors.general && (
            <Alert 
              type="error" 
              title="Login Error"
              dismissible
              onDismiss={() => setErrors(prev => ({ ...prev, general: '' }))}
            >
              {errors.general}
            </Alert>
          )}

          <form onSubmit={handleSubmit} className="login-form">
            <FormInput
              label="Email Address"
              name="email"
              type="email"
              value={formData.email}
              onChange={handleInputChange}
              placeholder="Enter your email"
              error={errors.email}
              required
              icon="📧"
              autoComplete="email"
              autoFocus
            />

            <FormInput
              label="Password"
              name="password"
              type={showPassword ? 'text' : 'password'}
              value={formData.password}
              onChange={handleInputChange}
              placeholder="Enter your password"
              error={errors.password}
              required
              icon="🔒"
              autoComplete="current-password"
            />

            <div className="password-toggle">
              <Button
                type="button"
                variant="ghost"
                size="small"
                onClick={() => setShowPassword(!showPassword)}
                icon={showPassword ? "👁️" : "👁️‍🗨️"}
              >
                {showPassword ? 'Hide' : 'Show'} password
              </Button>
            </div>

            <Button
              type="submit"
              variant="primary"
              size="large"
              loading={isLoading}
              fullWidth
              className="login-button"
            >
              {isLoading ? 'Signing in...' : 'Sign In'}
            </Button>
          </form>

          <div className="login-divider">
            <span>or</span>
          </div>

          <div className="demo-login">
            <h3 className="demo-title">Try Demo Accounts</h3>
            <div className="demo-buttons">
              <Button
                variant="outline"
                size="small"
                onClick={() => handleDemoLogin('admin')}
                disabled={isLoading}
                icon="👨‍💼"
              >
                Admin
              </Button>
              <Button
                variant="outline"
                size="small"
                onClick={() => handleDemoLogin('landlord')}
                disabled={isLoading}
                icon="🏠"
              >
                Landlord
              </Button>
              <Button
                variant="outline"
                size="small"
                onClick={() => handleDemoLogin('tenant')}
                disabled={isLoading}
                icon="👤"
              >
                Tenant
              </Button>
              <Button
                variant="outline"
                size="small"
                onClick={() => handleDemoLogin('maintenance')}
                disabled={isLoading}
                icon="🔧"
              >
                Maintenance
              </Button>
            </div>
          </div>

          <div className="login-footer">
            <p className="register-link">
              Don't have an account?{' '}
              <Link to="/register" className="text-primary">
                Sign up here
              </Link>
            </p>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default Login;