import React, { useState, useCallback } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { toast } from 'react-toastify';

import useAuth from '../hooks/useAuth';
import useFormValidation from '../hooks/useFormValidation';
import { validationRules } from '../utils/validation';
import { ROLE_ROUTES, ROLE_OPTIONS } from '../constants';

import Card from '../components/common/Card';
import Button from '../components/common/Button';
import FormInput from '../components/forms/FormInput';
import FormSelect from '../components/forms/FormSelect';
import Alert from '../components/common/Alert';
import LoadingSpinner from '../components/common/LoadingSpinner';

const Register = () => {
  const { register, loading, error, clearError } = useAuth();
  const navigate = useNavigate();
  const [isSubmitting, setIsSubmitting] = useState(false);

  const initialValues = {
    username: '',
    firstName: '',
    lastName: '',
    email: '',
    password: '',
    confirmPassword: '',
    role: 'tenant',
    phone: '',
  };

  const validationConfig = {
    username: [validationRules.username],
    firstName: [validationRules.required('First Name')],
    lastName: [validationRules.required('Last Name')],
    email: [validationRules.required('Email'), validationRules.email],
    password: [validationRules.password],
    confirmPassword: [
      validationRules.required('Confirm Password'),
      (value) => validationRules.confirmPassword(values.password)(value),
    ],
    phone: [validationRules.phone],
  };

  const {
    values,
    errors,
    handleChange,
    handleBlur,
    handleSubmit,
  } = useFormValidation(initialValues, validationConfig);

  const onSubmit = useCallback(async (formData) => {
    try {
      setIsSubmitting(true);
      clearError();
      
      // Remove confirmPassword from the data sent to the server
      const { confirmPassword, ...registrationData } = formData;
      
      // Map form fields to API fields
      const apiData = {
        username: registrationData.username,
        first_name: registrationData.firstName,
        last_name: registrationData.lastName,
        email: registrationData.email,
        password: registrationData.password,
        role: registrationData.role,
        phone: registrationData.phone,
      };
      
      const response = await register(apiData);
      
      // Get user role and redirect appropriately
      const userRole = response.user?.role || registrationData.role;
      const redirectPath = ROLE_ROUTES[userRole] || '/';
      
      toast.success('Registration successful! Welcome to Property Manager.');
      navigate(redirectPath);
    } catch (err) {
      console.error('Registration failed:', err);
      toast.error(err.message || 'Registration failed. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  }, [register, navigate, clearError]);

  if (loading) {
    return (
      <div className="form-container">
        <LoadingSpinner />
      </div>
    );
  }

  return (
    <div className="form-container">
      <Card>
        <form onSubmit={handleSubmit(onSubmit)}>
          <h2 className="form-title">Create Account</h2>
          
          {error && (
            <Alert
              type="error"
              message={error}
              onClose={clearError}
              className="mb-3"
            />
          )}

          <FormInput
            name="firstName"
            label="First Name"
            type="text"
            placeholder="Enter your first name"
            value={values.firstName}
            onChange={handleChange}
            onBlur={handleBlur}
            error={errors.firstName}
            required
          />

          <FormInput
            name="lastName"
            label="Last Name"
            type="text"
            placeholder="Enter your last name"
            value={values.lastName}
            onChange={handleChange}
            onBlur={handleBlur}
            error={errors.lastName}
            required
          />

          <FormInput
            name="username"
            label="Username"
            type="text"
            placeholder="Choose a username"
            value={values.username}
            onChange={handleChange}
            onBlur={handleBlur}
            error={errors.username}
            required
          />

          <FormInput
            name="email"
            label="Email Address"
            type="email"
            placeholder="Enter your email address"
            value={values.email}
            onChange={handleChange}
            onBlur={handleBlur}
            error={errors.email}
            required
          />

          <FormInput
            name="phone"
            label="Phone Number"
            type="tel"
            placeholder="Enter your phone number"
            value={values.phone}
            onChange={handleChange}
            onBlur={handleBlur}
            error={errors.phone}
          />

          <FormSelect
            name="role"
            label="Role"
            options={ROLE_OPTIONS}
            value={values.role}
            onChange={handleChange}
            onBlur={handleBlur}
            error={errors.role}
            required
          />

          <FormInput
            name="password"
            label="Password"
            type="password"
            placeholder="Create a strong password"
            value={values.password}
            onChange={handleChange}
            onBlur={handleBlur}
            error={errors.password}
            required
          />

          <FormInput
            name="confirmPassword"
            label="Confirm Password"
            type="password"
            placeholder="Confirm your password"
            value={values.confirmPassword}
            onChange={handleChange}
            onBlur={handleBlur}
            error={errors.confirmPassword}
            required
          />

          <Button
            type="submit"
            variant="primary"
            size="large"
            disabled={isSubmitting}
            loading={isSubmitting}
            className="w-100"
          >
            {isSubmitting ? 'Creating Account...' : 'Create Account'}
          </Button>

          <div className="text-center mt-3">
            <p>
              Already have an account?{' '}
              <Link to="/login" className="text-primary">
                Sign in here
              </Link>
            </p>
          </div>
        </form>
      </Card>
    </div>
  );
};

export default Register;