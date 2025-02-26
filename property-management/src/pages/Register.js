// src/pages/Register.js
import React, { useContext, useState } from 'react';
import { AuthContext } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';

const Register = () => {
  const { register } = useContext(AuthContext);
  const navigate = useNavigate();
  const [userData, setUserData] = useState({
    username: '',
    first_name: '',
    last_name: '',
    email: '',
    password: '',
    role: 'tenant',
    phone: ''
  });
  const [error, setError] = useState(null);

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await register(userData);
      navigate('/login');
    } catch (err) {
      setError("Registration failed. Please check your information.");
    }
  };

  const containerStyle = {
    maxWidth: '500px',
    margin: '2rem auto',
    padding: '2rem',
    backgroundColor: 'white',
    borderRadius: '8px',
    boxShadow: '0 2px 15px rgba(0,0,0,0.1)',
  };

  const inputStyle = {
    width: '100%',
    padding: '0.8rem',
    marginBottom: '1rem',
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
    marginTop: '1rem',
  };

  return (
    <div style={containerStyle}>
      <h2 style={{ textAlign: 'center', color: '#2c3e50', marginBottom: '2rem' }}>Create Account</h2>
      {error && <p style={{ color: '#e74c3c', textAlign: 'center', marginBottom: '1rem' }}>{error}</p>}
      <form onSubmit={handleSubmit}>
        {Object.keys(userData).map((key) => (
          key !== 'role' && (
            <input
              key={key}
              type={key === 'password' ? 'password' : key === 'email' ? 'email' : 'text'}
              placeholder={key.replace('_', ' ').toUpperCase()}
              style={inputStyle}
              value={userData[key]}
              onChange={(e) => setUserData({...userData, [key]: e.target.value})}
            />
          )
        ))}
        <button type="submit" style={buttonStyle}>Register</button>
      </form>
    </div>
  );
};

export default Register;