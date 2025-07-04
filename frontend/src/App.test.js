import React from 'react';
import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import App from './App';

// Mock the AuthContext
jest.mock('./context/AuthContext', () => ({
  AuthProvider: ({ children }) => children,
}));

// Mock the lazy-loaded components
jest.mock('./pages/Login', () => {
  return function Login() {
    return <div>Login Page</div>;
  };
});

jest.mock('./pages/Register', () => {
  return function Register() {
    return <div>Register Page</div>;
  };
});

const renderApp = () => {
  return render(
    <BrowserRouter>
      <App />
    </BrowserRouter>
  );
};

test('renders without crashing', () => {
  renderApp();
  expect(document.body).toBeInTheDocument();
});

test('renders application structure', () => {
  renderApp();
  // Check if the main app div is rendered
  const appDiv = document.querySelector('.app');
  expect(appDiv).toBeInTheDocument();
});
