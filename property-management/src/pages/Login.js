import React, { useContext, useState } from 'react';
import { AuthContext } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';

const Login = () => {
  const { login } = useContext(AuthContext);
  const navigate = useNavigate();
  const [credentials, setCredentials] = useState({ username: '', password: '' });
  const [error, setError] = useState(null);

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await login(credentials);
      navigate('/');
    } catch (err) {
      setError("Invalid username or password");
    }
  };

  const containerStyle = {
    maxWidth: '400px',
    margin: '2rem auto',
    padding: '2rem',
    backgroundColor: 'white',
    borderRadius: '8px',
    boxShadow: '0 2px 15px rgba(0,0,0,0.1)',
  };

  const inputStyle = {
    width: '100%',
    padding: '0.8rem',
    marginBottom: '1.5rem',
    border: '1px solid #ddd',
    borderRadius: '4px',
    fontSize: '1rem',
  };

  const buttonStyle = {
    width: '100%',
    padding: '1rem',
    backgroundColor: '#3498db',
    color: 'white',
    border: 'none',
    borderRadius: '4px',
    fontSize: '1rem',
    cursor: 'pointer',
  };

  return (
    <div style={containerStyle}>
      <h2 style={{ textAlign: 'center', color: '#2c3e50', marginBottom: '2rem' }}>Login</h2>
      {error && (
        <p style={{ color: '#e74c3c', textAlign: 'center', marginBottom: '1rem' }}>
          {error}
        </p>
      )}
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          placeholder="Username"
          style={inputStyle}
          value={credentials.username}
          onChange={(e) =>
            setCredentials({ ...credentials, username: e.target.value })
          }
        />
        <input
          type="password"
          placeholder="Password"
          style={inputStyle}
          value={credentials.password}
          onChange={(e) =>
            setCredentials({ ...credentials, password: e.target.value })
          }
        />
        <button type="submit" style={buttonStyle}>
          Login
        </button>
      </form>
    </div>
  );
};

export default Login;