import React, { useState } from 'react';
import { toast } from 'react-toastify';
import './User.css';

const UserForm = ({ initialValues, onSubmit, onCancel }) => {
  const [formData, setFormData] = useState({
    username: initialValues.username || '',
    firstName: initialValues.first_name || '',
    lastName: initialValues.last_name || '',
    password: '',
    confirmPassword: '',
    email: initialValues.email || '',
    role: initialValues.role || 'tenant',
    phone: initialValues.phone || '',
  });

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = (e) => {
    e.preventDefault();

    if (!initialValues.id && formData.password !== formData.confirmPassword) {
      toast.error("Passwords don't match");
      return;
    }

    const apiData = {
      ...formData,
      first_name: formData.firstName,
      last_name: formData.lastName,
    };
    
    delete apiData.firstName;
    delete apiData.lastName;
    delete apiData.confirmPassword;

    onSubmit(apiData);
  };

  return (
    <form className="user-form" onSubmit={handleSubmit}>
      <div className="form-grid">
        <div className="form-group">
          <label>Username</label>
          <input
            type="text"
            name="username"
            value={formData.username}
            onChange={handleChange}
            className="form-input"
            required
          />
        </div>

        <div className="form-group">
          <label>First Name</label>
          <input
            type="text"
            name="firstName"
            value={formData.firstName}
            onChange={handleChange}
            className="form-input"
            required
          />
        </div>

        <div className="form-group">
          <label>Last Name</label>
          <input
            type="text"
            name="lastName"
            value={formData.lastName}
            onChange={handleChange}
            className="form-input"
            required
          />
        </div>

        <div className="form-group">
          <label>Email</label>
          <input
            type="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
            className="form-input"
            required
          />
        </div>

        <div className="form-group">
          <label>Password</label>
          <input
            type="password"
            name="password"
            value={formData.password}
            onChange={handleChange}
            className="form-input"
            required={!initialValues.id}
          />
        </div>

        {!initialValues.id && (
          <div className="form-group">
            <label>Confirm Password</label>
            <input
              type="password"
              name="confirmPassword"
              value={formData.confirmPassword}
              onChange={handleChange}
              className="form-input"
              required
            />
          </div>
        )}

        <div className="form-group">
          <label>Role</label>
          <select
            name="role"
            value={formData.role}
            onChange={handleChange}
            className="form-select form-input"
            required
          >
            <option value="admin">Admin</option>
            <option value="tenant">Tenant</option>
            <option value="landlord">Landlord</option>
            <option value="maintenanceTeam">Maintenance Team</option>
          </select>
        </div>

        <div className="form-group">
          <label>Phone</label>
          <input
            type="tel"
            name="phone"
            value={formData.phone}
            onChange={handleChange}
            className="form-input"
            pattern="[0-9]{10}"
            title="10-digit phone number"
            required
          />
        </div>

        <div className="form-actions">
          <button
            type="button"
            onClick={onCancel}
            className="btn-secondary"
          >
            Cancel
          </button>
          <button
            type="submit"
            className="btn-primary"
          >
            {initialValues.id ? 'Update' : 'Create'}
          </button>
        </div>
      </div>
    </form>
  );
};

export default UserForm;