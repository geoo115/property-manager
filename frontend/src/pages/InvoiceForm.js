import React, { useState, useEffect } from 'react';
import { toast } from 'react-toastify';
import axiosInstance from '../api/axiosInstance';

const InvoiceForm = ({ initialValues = {}, onSubmit, onCancel }) => {
  const [formData, setFormData] = useState({
    tenant_id: initialValues.tenant_id || '',
    property_id: initialValues.property_id || '',
    amount: initialValues.amount || '',
    paid_amount: initialValues.paid_amount || '0',
    invoice_date: initialValues.invoice_date
      ? initialValues.invoice_date.split('T')[0]
      : new Date().toISOString().split('T')[0],
    due_date: initialValues.due_date
      ? initialValues.due_date.split('T')[0]
      : new Date().toISOString().split('T')[0],
    category: initialValues.category || 'rent',
    payment_status: initialValues.payment_status || 'pending',
    refunded_amount: initialValues.refunded_amount || '0',
    recurring_interval: initialValues.recurring_interval || '',
    recurring: initialValues.recurring || false,
  });

  const [tenants, setTenants] = useState([]);
  const [properties, setProperties] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [tenantsResponse, propertiesResponse] = await Promise.all([
          axiosInstance.get('admin/users/active'),
          axiosInstance.get('admin/properties'),
        ]);
        setTenants(tenantsResponse.data.users || []);
        setProperties(propertiesResponse.data.properties || []);
        setLoading(false);
      } catch (err) {
        console.error('Error fetching data:', err);
        setError('Failed to fetch data');
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value,
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (isSubmitting) return;
    setIsSubmitting(true);

    const payload = {
      tenant_id: parseInt(formData.tenant_id, 10),
      property_id: parseInt(formData.property_id, 10),
      amount: parseFloat(formData.amount),
      paid_amount: parseFloat(formData.paid_amount),
      invoice_date: formData.invoice_date,
      due_date: formData.due_date,
      category: formData.category,
      payment_status: formData.payment_status,
      refunded_amount: parseFloat(formData.refunded_amount) || 0,
      recurring_interval: formData.recurring_interval,
      recurring: formData.recurring,
    };

    console.log('Submitting payload:', payload);

    if (isNaN(payload.tenant_id) || isNaN(payload.property_id)) {
      toast.error('Please select a tenant and property');
      setIsSubmitting(false);
      return;
    }
    if (isNaN(payload.amount) || payload.amount <= 0) {
      toast.error('Amount must be a positive number');
      setIsSubmitting(false);
      return;
    }
    if (isNaN(payload.paid_amount) || payload.paid_amount < 0) {
      toast.error('Paid amount must be non-negative');
      setIsSubmitting(false);
      return;
    }
    if (new Date(payload.due_date) < new Date(payload.invoice_date)) {
      toast.error('Due date cannot be before invoice date');
      setIsSubmitting(false);
      return;
    }

    try {
      let response;
      if (initialValues.id) {
        response = await axiosInstance.put(`/admin/accounting/invoices/${initialValues.id}`, payload);
        toast.success('Invoice updated successfully');
      } else {
        response = await axiosInstance.post('/admin/accounting/invoices', payload);
        toast.success('Invoice created successfully');
      }
      if (onSubmit) onSubmit(payload); // Pass original payload, not response.data
    } catch (error) {
      console.error('Error submitting invoice:', error);
      if (error.response) {
        console.error('Server response:', error.response.data);
        toast.error(error.response.data.details || error.response.data.error || 'Failed to submit invoice');
      } else {
        toast.error('Failed to submit invoice');
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  if (loading) return <p>Loading data...</p>;
  if (error) return <p>{error}</p>;

  return (
    <div className="modal-overlay" onClick={onCancel}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <button className="modal-close" onClick={onCancel}>×</button>
        <h2>{initialValues.id ? 'Edit Invoice' : 'New Invoice'}</h2>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="tenant_id">Tenant</label>
            <select
              id="tenant_id"
              name="tenant_id"
              value={formData.tenant_id}
              onChange={handleChange}
              required
            >
              <option value="">Select a tenant</option>
              {tenants.map((tenant) => (
                <option key={tenant.id} value={tenant.id}>
                  {tenant.first_name} {tenant.last_name}
                </option>
              ))}
            </select>
          </div>

          <div className="form-group">
            <label htmlFor="property_id">Property</label>
            <select
              id="property_id"
              name="property_id"
              value={formData.property_id}
              onChange={handleChange}
              required
            >
              <option value="">Select a property</option>
              {properties.map((property) => (
                <option key={property.id} value={property.id}>
                  {property.address || property.name || `Property ${property.id}`}
                </option>
              ))}
            </select>
          </div>

          <div className="form-group">
            <label htmlFor="amount">Amount</label>
            <input
              id="amount"
              type="number"
              name="amount"
              value={formData.amount}
              onChange={handleChange}
              min="0.01"
              step="0.01"
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="paid_amount">Paid Amount</label>
            <input
              id="paid_amount"
              type="number"
              name="paid_amount"
              value={formData.paid_amount}
              onChange={handleChange}
              min="0"
              step="0.01"
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="invoice_date">Invoice Date</label>
            <input
              id="invoice_date"
              type="date"
              name="invoice_date"
              value={formData.invoice_date}
              onChange={handleChange}
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="due_date">Due Date</label>
            <input
              id="due_date"
              type="date"
              name="due_date"
              value={formData.due_date}
              onChange={handleChange}
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="category">Category</label>
            <select
              id="category"
              name="category"
              value={formData.category}
              onChange={handleChange}
              required
            >
              <option value="rent">Rent</option>
              <option value="deposit">Deposit</option>
            </select>
          </div>

          <div className="form-group">
            <label htmlFor="payment_status">Payment Status</label> {/* Changed from 'status' */}
            <select
  id="payment_status"
  name="payment_status"
  value={formData.payment_status}
  onChange={handleChange}
  required
>
              <option value="pending">Pending</option>
              <option value="paid">Paid</option>
              <option value="overdue">Overdue</option>
            </select>
          </div>

          {/* Additional fields for UpdateInvoice */}
          <div className="form-group">
            <label htmlFor="refunded_amount">Refunded Amount</label>
            <input
              id="refunded_amount"
              type="number"
              name="refunded_amount"
              value={formData.refunded_amount}
              onChange={handleChange}
              min="0"
              step="0.01"
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="recurring_interval">Recurring Interval</label>
            <select
              id="recurring_interval"
              name="recurring_interval"
              value={formData.recurring_interval}
              onChange={handleChange}
            >
              <option value="">None</option>
              <option value="monthly">Monthly</option>
              <option value="yearly">Yearly</option>
            </select>
          </div>

          <div className="form-group">
            <label htmlFor="recurring">Recurring</label>
            <input
              id="recurring"
              type="checkbox"
              name="recurring"
              checked={formData.recurring}
              onChange={handleChange}
            />
          </div>

          <div className="form-actions">
            <button type="button" onClick={onCancel}>Cancel</button>
            <button type="submit" disabled={isSubmitting}>
              {initialValues.id ? 'Update' : 'Create'} Invoice
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default InvoiceForm;