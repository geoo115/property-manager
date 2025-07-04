import React, { useEffect, useState, useContext } from 'react';
import { AuthContext } from '../context/AuthContext';
import { getProperties, createProperty, updateProperty, deleteProperty } from '../api/properties';
import PropertyForm from './PropertyForm';
import DataTable from '../components/common/DataTable';
import Button from '../components/common/Button';
import RoleBasedActions from '../components/common/RoleBasedActions';
import { AdminOnly, LandlordOnly } from '../components/common/RoleBasedContent';
import { toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

const Properties = () => {
  const { user, hasUserPermission } = useContext(AuthContext);
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
    if (!hasUserPermission('properties', 'delete')) {
      toast.error('You do not have permission to delete properties');
      return;
    }
    
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
    const action = currentProperty ? 'update' : 'create';
    if (!hasUserPermission('properties', action)) {
      toast.error(`You do not have permission to ${action} properties`);
      return;
    }
    
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

  // Define table columns
  const columns = [
    { key: 'name', title: 'Name', sortable: true },
    { key: 'address', title: 'Address', sortable: true },
    { key: 'city', title: 'City', sortable: true },
    { key: 'post_code', title: 'Post Code', sortable: true },
    { key: 'bedrooms', title: 'Bedrooms', sortable: true },
    { 
      key: 'price', 
      title: 'Price', 
      sortable: true,
      render: (value) => `£${value}`
    },
    {
      key: 'actions',
      title: 'Actions',
      render: (value, row) => (
        <RoleBasedActions
          resource="properties"
          item={row}
          onView={() => {
            setCurrentProperty(row);
            setIsModalOpen(true);
          }}
          onEdit={() => {
            setCurrentProperty(row);
            setIsModalOpen(true);
          }}
          onDelete={() => handleDelete(row.id)}
        />
      )
    }
  ];

  if (error) {
    return (
      <div className="page-container">
        <div className="error-state">
          <h3>Error Loading Properties</h3>
          <p>{error}</p>
          <Button onClick={fetchProperties}>Try Again</Button>
        </div>
      </div>
    );
  }

  return (
    <div className="page-container">
      <div className="page-header">
        <div>
          <h1>Properties</h1>
          <p className="page-description">Manage your property portfolio</p>
        </div>
        <AdminOnly>
          <Button
            variant="primary"
            onClick={() => {
              setCurrentProperty(null);
              setIsModalOpen(true);
            }}
          >
            Add Property
          </Button>
        </AdminOnly>
      </div>

      <DataTable
        data={properties}
        columns={columns}
        loading={isLoading}
        searchable={true}
        searchPlaceholder="Search properties..."
        emptyMessage="No properties found"
      />

      {isModalOpen && (
        <div className="modal-backdrop">
          <div className="modal-container">
            <div className="modal-header">
              <h3>{currentProperty ? 'Edit Property' : 'Create New Property'}</h3>
              <button 
                className="modal-close"
                onClick={() => setIsModalOpen(false)}
                aria-label="Close modal"
              >
                ×
              </button>
            </div>
            <div className="modal-content">
              <PropertyForm
                initialValues={currentProperty || {}}
                onSubmit={handleSubmit}
                onCancel={() => setIsModalOpen(false)}
              />
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Properties;
