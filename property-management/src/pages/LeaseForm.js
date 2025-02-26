import React, { useState, useEffect } from 'react';
import { toast } from 'react-toastify';
import axiosInstance from '../api/axiosInstance';
import './Leases.css';

const LeaseForm = ({ initialValues, onSubmit, onCancel }) => {
  // Utility function to convert ISO dates to YYYY-MM-DD format
  const parseISOToDateInput = (isoString) => {
    if (!isoString) return '';
    try {
      return new Date(isoString).toISOString().split('T')[0];
    } catch {
      return '';
    }
  };

  // State for form data
  const [formData, setFormData] = useState({
    tenant_id: initialValues.tenant_id || '',
    property_id: initialValues.property_id || '',
    start_date: parseISOToDateInput(initialValues.start_date),
    end_date: parseISOToDateInput(initialValues.end_date),
    monthly_rent: initialValues.monthly_rent || 0,
    security_deposit: initialValues.security_deposit || 0,
  });

  const [tenants, setTenants] = useState([]);
  const [properties, setProperties] = useState([]);
  const [loading, setLoading] = useState(true);

  // Fetch tenants and properties on component mount
  useEffect(() => {
    const fetchData = async () => {
      try {
        const [tenantsRes, propertiesRes] = await Promise.all([
          axiosInstance.get('/admin/users?role=tenant'),
          axiosInstance.get('/admin/properties'),
        ]);

        setTenants(tenantsRes.data.users);
        setProperties(propertiesRes.data.properties);
        setLoading(false);
      } catch (error) {
        toast.error('Failed to load required data');
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  // Update form data when initialValues change
  useEffect(() => {
    setFormData({
      ...formData,
      start_date: parseISOToDateInput(initialValues.start_date),
      end_date: parseISOToDateInput(initialValues.end_date),
    });
  }, [initialValues]);

  // Handle form input changes
  const handleChange = (e) => {
    let value;

    if (e.target.type === 'number') {
      value = parseFloat(e.target.value);
    } else if (e.target.tagName === 'SELECT') {
      value = parseInt(e.target.value, 10);
    } else {
      value = e.target.value;
    }

    setFormData({
      ...formData,
      [e.target.name]: value,
    });
  };

  // Handle form submission
  const handleSubmit = (e) => {
    e.preventDefault();

    // Validate form data
    if (!formData.tenant_id) {
      toast.error('Please select a tenant');
      return;
    }

    if (!formData.property_id) {
      toast.error('Please select a property');
      return;
    }

    if (new Date(formData.end_date) <= new Date(formData.start_date)) {
      toast.error('End date must be after start date');
      return;
    }

    if (formData.monthly_rent <= 0) {
      toast.error('Monthly rent must be greater than 0');
      return;
    }

    // Convert dates to ISO format before submission
    const processedData = {
      ...formData,
      start_date: new Date(formData.start_date).toISOString(),
      end_date: new Date(formData.end_date).toISOString(),
    };

    onSubmit(processedData);
  };

  if (loading) {
    return <div className="loading-container">Loading form data...</div>;
  }

  return (
    <form className="lease-form" onSubmit={handleSubmit}>
      <div className="form-grid">
        {/* Tenant Select */}
        <div className="form-group">
          <label>Tenant</label>
          <select
            name="tenant_id"
            value={formData.tenant_id}
            onChange={handleChange}
            className="form-select"
            required
          >
            <option value="">Select Tenant</option>
            {tenants.map((tenant) => (
              <option key={tenant.id} value={tenant.id}>
                {tenant.first_name} {tenant.last_name}
              </option>
            ))}
          </select>
        </div>

        {/* Property Select */}
        <div className="form-group">
          <label>Property</label>
          <select
            name="property_id"
            value={formData.property_id}
            onChange={handleChange}
            className="form-select"
            required
          >
            <option value="">Select Property</option>
            {properties.map((property) => (
              <option key={property.id} value={property.id}>
                {property.name} ({property.address})
              </option>
            ))}
          </select>
        </div>

        {/* Start Date Input */}
        <div className="form-group">
          <label>Start Date</label>
          <input
            type="date"
            name="start_date"
            value={formData.start_date}
            onChange={handleChange}
            className="form-control"
            required
          />
        </div>

        {/* End Date Input */}
        <div className="form-group">
          <label>End Date</label>
          <input
            type="date"
            name="end_date"
            value={formData.end_date}
            onChange={handleChange}
            className="form-control"
            required
          />
        </div>

        {/* Monthly Rent Input */}
        <div className="form-group">
          <label>Monthly Rent (£)</label>
          <input
            type="number"
            name="monthly_rent"
            value={formData.monthly_rent}
            onChange={handleChange}
            min="0.01"
            step="0.01"
            className="form-control"
            required
          />
        </div>

        {/* Security Deposit Input */}
        <div className="form-group">
          <label>Security Deposit (£)</label>
          <input
            type="number"
            name="security_deposit"
            value={formData.security_deposit}
            onChange={handleChange}
            min="0.01"
            step="0.01"
            className="form-control"
            required
          />
        </div>

        {/* Form Actions */}
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

export default LeaseForm;