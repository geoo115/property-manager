import React, { useState, useEffect } from 'react';
import { toast } from 'react-toastify';
import axiosInstance from '../api/axiosInstance';
import './Properties.css';

const PropertyForm = ({ initialValues, onSubmit, onCancel }) => {
  const [formData, setFormData] = useState({
    name: initialValues.name || '',
    description: initialValues.description || '',
    bedrooms: initialValues.bedrooms || 1,
    bathrooms: initialValues.bathrooms || 1,
    price: initialValues.price || 0,
    square_feet: initialValues.square_feet || 0,
    address: initialValues.address || '',
    city: initialValues.city || '',
    post_code: initialValues.post_code || '',
    owner_id: initialValues.owner_id || '',
    available: initialValues.available || true,
  });

  const [owners, setOwners] = useState([]);
  const [loadingOwners, setLoadingOwners] = useState(true);
  const [ownersError, setOwnersError] = useState('');

  useEffect(() => {
    const fetchOwners = async () => {
      try {
        const response = await axiosInstance.get('/admin/users', {
          params: { role: 'landlord' }
        });
        setOwners(response.data.users);
        setLoadingOwners(false);
      } catch (error) {
        setOwnersError('Failed to fetch owners');
        setLoadingOwners(false);
      }
    };
    fetchOwners();
  }, []);

  const handleChange = (e) => {
    const value = e.target.type === 'checkbox' 
      ? e.target.checked
      : e.target.type === 'number'
      ? parseFloat(e.target.value)
      : e.target.value;

    setFormData({ ...formData, [e.target.name]: value });
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    
    if (formData.bedrooms < 1) {
      toast.error("Bedrooms must be at least 1");
      return;
    }
  
    if (formData.price <= 0) {
      toast.error("Price must be greater than 0");
      return;
    }
  
    if (!formData.owner_id) {
      toast.error("Please select an owner");
      return;
    }
  
    // Ensure owner_id is a number
    const dataToSend = { ...formData, owner_id: Number(formData.owner_id) };
  
    onSubmit(dataToSend);
  };
  

  return (
    <form className="property-form" onSubmit={handleSubmit}>
      <div className="form-grid">
        <div className="form-group">
          <label>Property Name</label>
          <input
            type="text"
            name="name"
            value={formData.name}
            onChange={handleChange}
            required
          />
        </div>

        <div className="form-group floating-label">
          <label>Description</label>
          <textarea
            name="description"
            value={formData.description}
            onChange={handleChange}
            required
          />
        </div>

        <div className="form-group">
          <label>Bedrooms</label>
          <input
            type="number"
            name="bedrooms"
            value={formData.bedrooms}
            onChange={handleChange}
            min="1"
            required
          />
        </div>

        <div className="form-group">
          <label>Bathrooms</label>
          <input
            type="number"
            name="bathrooms"
            value={formData.bathrooms}
            onChange={handleChange}
            min="1"
            required
          />
        </div>

        <div className="form-group">
          <label>Price (£)</label>
          <input
            type="number"
            name="price"
            value={formData.price}
            onChange={handleChange}
            min="0.01"
            step="0.01"
            required
          />
        </div>

        <div className="form-group">
          <label>Square Feet</label>
          <input
            type="number"
            name="square_feet"
            value={formData.square_feet}
            onChange={handleChange}
            min="1"
            required
          />
        </div>

        <div className="form-group">
          <label>Address</label>
          <input
            type="text"
            name="address"
            value={formData.address}
            onChange={handleChange}
            required
          />
        </div>

        <div className="form-group">
          <label>City</label>
          <input
            type="text"
            name="city"
            value={formData.city}
            onChange={handleChange}
            required
          />
        </div>

        <div className="form-group">
          <label>Post Code</label>
          <input
            type="text"
            name="post_code"
            value={formData.post_code}
            onChange={handleChange}
            required
          />
        </div>

        <div className="form-group">
          <label>Owner</label>
          {loadingOwners ? (
            <div className="loading-text">Loading owners...</div>
          ) : ownersError ? (
            <div className="error-text">{ownersError}</div>
          ) : (
            <select
              name="owner_id"
              value={formData.owner_id}
              onChange={handleChange}
              required
            >
              <option value="">Select Owner</option>
              {owners.map((owner) => (
                <option key={owner.id} value={owner.id}>
                  {owner.first_name} {owner.last_name} ({owner.email})
                </option>
              ))}
            </select>
          )}
        </div>

        <div className="form-group checkbox-group">
          <label>Available</label>
          <input
            type="checkbox"
            name="available"
            checked={formData.available}
            onChange={handleChange}
          />
        </div>

        <div className="form-actions">
          <button type="button" className="btn-secondary" onClick={onCancel}>
            Cancel
          </button>
          <button type="submit" className="btn-primary">
            {initialValues.id ? 'Update' : 'Create'}
          </button>
        </div>
      </div>
    </form>
  );
};

export default PropertyForm;