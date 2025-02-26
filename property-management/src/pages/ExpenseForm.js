// ExpenseForm.js
import React, { useState } from 'react';

const ExpenseForm = ({ initialValues = {}, onSubmit, onCancel, properties }) => {
  const initialDate = initialValues.expense_date
    ? new Date(initialValues.expense_date)
    : new Date();

  const [formData, setFormData] = useState({
    description: initialValues.description || '',
    amount: initialValues.amount || '',
    expense_date: `${initialDate.getFullYear()}-${String(initialDate.getMonth() + 1).padStart(2, '0')}-${String(initialDate.getDate()).padStart(2, '0')}`,
    property_id: initialValues.property_id || (properties.length > 0 ? properties[0].id : ''),
    category: initialValues.category || '',
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    const convertedFormData = {
      ...formData,
      amount: parseFloat(formData.amount), // Convert amount to float
      property_id: parseInt(formData.property_id, 10), // Convert property_id to integer
    };
    onSubmit(convertedFormData);
  };
  
  

  return (
    <div className="modal-overlay" onClick={onCancel}>
      <div className="modal-content" onClick={(e) => e.stopPropagation()}>
        <button className="modal-close" onClick={onCancel}>&times;</button>
        <h2>{initialValues.id ? 'Edit Expense' : 'New Expense'}</h2>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label className="form-label">Description</label>
            <input
              type="text"
              name="description"
              className="form-input"
              value={formData.description}
              onChange={handleChange}
              required
            />
          </div>
          <div className="form-group">
            <label className="form-label">Amount</label>
            <input
              type="number"
              name="amount"
              className="form-input"
              value={formData.amount}
              onChange={handleChange}
              required
            />
          </div>
          <div className="form-group">
            <label className="form-label">Date</label>
            <input
              type="date"
              name="expense_date"
              className="form-input"
              value={formData.expense_date}
              onChange={handleChange}
              required
            />
          </div>
          <div className="form-group">
            <label className="form-label">Property</label>
            <select
              name="property_id"
              className="form-input"
              value={formData.property_id}
              onChange={handleChange}
              required
            >
              {properties.map(property => (
                <option key={property.id} value={property.id}>
                  {property.name}
                </option>
              ))}
            </select>
          </div>
          <div className="form-group">
              <label className="form-label">Category</label>
              <input
                type="text"
                name="category"
                className="form-input"
                value={formData.category}
                onChange={handleChange}
                required
              />
          </div>
          <div className="form-actions">
            <button type="button" className="button button-secondary" onClick={onCancel}>
              Cancel
            </button>
            <button type="submit" className="button button-primary">
              {initialValues.id ? 'Update' : 'Create'} Expense
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default ExpenseForm;