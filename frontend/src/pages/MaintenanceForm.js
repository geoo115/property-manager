import React, { useState, useEffect, useContext } from 'react';
import { AuthContext } from '../context/AuthContext';
import axiosInstance from '../api/axiosInstance';
import { toast } from 'react-toastify';
import './maintenance.css';

const MaintenanceForm = ({ initialValues, onSubmit, onCancel, leaseId, userRole }) => {
  const { user } = useContext(AuthContext);
  const [formData, setFormData] = useState({
    description: initialValues.description || '',
    property_id: initialValues.property_id || '',
    status: initialValues.status || 'pending',
  });

  const [properties, setProperties] = useState([]);
  const [loading, setLoading] = useState(true);

  const statusOptions = [
    { value: 'pending', label: 'Pending' },
    { value: 'in_progress', label: 'In Progress' },
    { value: 'completed', label: 'Completed' },
    { value: 'cancelled', label: 'Cancelled' },
  ];

  useEffect(() => {
    const fetchData = async () => {
      try {
        if (userRole === 'tenant' && leaseId) {
          const { data } = await axiosInstance.get(`/tenant/leases/${leaseId}`);
          setFormData(prev => ({
            ...prev,
            property_id: data.lease?.property_id || '',
            description: initialValues.description || prev.description
          }));
        }

        if (userRole === 'admin') {
          const { data } = await axiosInstance.get('/admin/properties');
          setProperties(data.properties || []);
        }
      } catch (error) {
        toast.error('Failed to load required data');
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [leaseId, userRole, initialValues.description]);

  const handleSubmit = (e) => {
    e.preventDefault();
    
    if (!formData.description || (userRole === 'admin' && !formData.property_id)) {
      toast.error('Please fill all required fields');
      return;
    }

    // Only send updatable fields
    const payload = {
      description: formData.description,
      status: formData.status,
      ...(userRole === 'admin' && { property_id: formData.property_id })
    };

    onSubmit(payload);
  };

  return (
    <form className="maintenance-form" onSubmit={handleSubmit}>
      <div className="form-grid">
        {(userRole === 'admin' || userRole === 'maintenanceTeam') && (
          <div className="form-group">
            <label>Status *</label>
            <select
              value={formData.status}
              onChange={(e) => setFormData({ ...formData, status: e.target.value })}
              className="form-select"
              required
              disabled={loading}
            >
              {statusOptions.map(option => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </div>
        )}

        {userRole === 'tenant' && (
          <>
            <div className="form-group">
              <label>Associated Lease</label>
              <input type="text" value={leaseId} readOnly className="form-input" />
            </div>
            <div className="form-group">
              <label>Property ID</label>
              <input type="text" value={formData.property_id} readOnly className="form-input" />
            </div>
          </>
        )}

        {userRole === 'admin' && (
          <div className="form-group">
            <label>Select Property *</label>
            <select
              value={formData.property_id}
              onChange={(e) => setFormData({ ...formData, property_id: e.target.value })}
              className="form-select"
              required
              disabled={loading}
            >
              <option value="">Select Property</option>
              {properties.map(property => (
                <option key={property.id} value={property.id}>
                  {property.name} ({property.address})
                </option>
              ))}
            </select>
          </div>
        )}

        <div className="form-group full-width">
          <label>Description *</label>
          <textarea
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            className="form-input"
            required
            rows="4"
            placeholder="Describe the maintenance issue..."
            disabled={loading}
          />
        </div>

        <div className="form-actions">
          <button type="button" onClick={onCancel} className="btn-secondary">
            Cancel
          </button>
          <button type="submit" className="btn-primary" disabled={loading}>
            {loading ? 'Loading...' : (initialValues.id ? 'Update' : 'Create')}
          </button>
        </div>
      </div>
    </form>
  );
};

export default MaintenanceForm;