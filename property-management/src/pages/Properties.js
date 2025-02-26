import React, { useEffect, useState, useContext } from 'react';
import { AuthContext } from '../context/AuthContext';
import { getProperties, createProperty, updateProperty, deleteProperty } from '../api/properties';
import PropertyForm from './PropertyForm';
import { TailSpin } from 'react-loader-spinner';
import { toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import './Properties.css';

const Properties = () => {
  const { user } = useContext(AuthContext);
  const [properties, setProperties] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [currentProperty, setCurrentProperty] = useState(null);
  const [error, setError] = useState('');

  // Fetch properties based on role
  useEffect(() => {
    fetchProperties();
  }, []);

  const fetchProperties = async () => {
    try {
      const data = await getProperties();
      setProperties(data.properties || []);
    } catch (error) {
      setError('Failed to fetch properties');
    } finally {
      setIsLoading(false);
    }
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Are you sure you want to delete this property?')) return;

    try {
      await deleteProperty(id);
      toast.success('Property deleted successfully');
      fetchProperties();
    } catch (error) {
      toast.error(error.response?.data?.message || 'Failed to delete property');
    }
  };

  const handleSubmit = async (propertyData) => {
    try {
      if (currentProperty) {
        await updateProperty(currentProperty.id, propertyData);
        toast.success('Property updated successfully');
      } else {
        await createProperty(propertyData);
        toast.success('Property created successfully');
      }
      setIsModalOpen(false);
      fetchProperties();
    } catch (error) {
      toast.error(error.response?.data?.message || 'Operation failed');
    }
  };

  if (isLoading) {
    return (
      <div className="spinner-container">
        <TailSpin color="#3B82F6" height={80} width={80} />
      </div>
    );
  }

  if (error) {
    return <div className="error-message">{error}</div>;
  }

  return (
    <div className="properties-container">
      <div className="properties-header">
        <h2>Property Management</h2>
        {user.role === 'admin' && (
          <button
            className="btn-primary"
            onClick={() => {
              setCurrentProperty(null);
              setIsModalOpen(true);
            }}
          >
            Add New Property
          </button>
        )}
      </div>

      <div className="properties-table-container">
        <table className="properties-table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Address</th>
              <th>City</th>
              <th>Post Code</th>
              <th>Bedrooms</th>
              <th>Price</th>
              {user.role === 'admin' && <th>Actions</th>}
            </tr>
          </thead>
          <tbody>
            {properties.map((property) => (
              <tr key={property.id}>
                <td>{property.name}</td>
                <td>{property.address}</td>
                <td>{property.city}</td>
                <td>{property.post_code}</td>
                <td>{property.bedrooms}</td>
                <td>£{property.price}</td>
                {user.role === 'admin' && (
                  <td className="actions-cell">
                    <button
                      className="btn-edit"
                      onClick={() => {
                        setCurrentProperty(property);
                        setIsModalOpen(true);
                      }}
                    >
                      Edit
                    </button>
                    <button className="btn-delete" onClick={() => handleDelete(property.id)}>
                      Delete
                    </button>
                  </td>
                )}
              </tr>
            ))}
          </tbody>
        </table>
        {properties.length === 0 && <div className="no-properties">No properties found</div>}
      </div>

      {isModalOpen && (
        <div className="modal-backdrop">
          <div className="modal-container">
            <h3>{currentProperty ? 'Edit Property' : 'Create New Property'}</h3>
            <PropertyForm
              initialValues={currentProperty || {}}
              onSubmit={handleSubmit}
              onCancel={() => setIsModalOpen(false)}
            />
          </div>
        </div>
      )}
    </div>
  );
};

export default Properties;
