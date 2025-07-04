import React, { useEffect, useState, useContext } from 'react';
import { AuthContext } from '../context/AuthContext';
import { getLeases, createLease, updateLease, deleteLease } from '../api/leases';
import LeaseForm from './LeaseForm';
import DataTable from '../components/common/DataTable';
import Button from '../components/common/Button';
import RoleBasedActions from '../components/common/RoleBasedActions';
import { AdminOnly, LandlordOnly } from '../components/common/RoleBasedContent';
import { toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

const Leases = () => {
  const { hasUserPermission } = useContext(AuthContext);
  const [leases, setLeases] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [currentLease, setCurrentLease] = useState(null);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchLeases();
  }, []);

  const fetchLeases = async () => {
    try {
      const data = await getLeases();
      setLeases(data.leases || []);
      setIsLoading(false);
    } catch (error) {
      setError('Failed to fetch leases');
      setIsLoading(false);
    }
  };

  const handleDelete = async (id) => {
    if (!hasUserPermission('leases', 'delete')) {
      toast.error('You do not have permission to delete leases');
      return;
    }
    
    if (window.confirm('Are you sure you want to delete this lease?')) {
      try {
        await deleteLease(id);
        toast.success('Lease deleted successfully');
        fetchLeases();
      } catch (error) {
        toast.error('Failed to delete lease');
      }
    }
  };

  const handleSubmit = async (leaseData) => {
    const action = currentLease ? 'update' : 'create';
    if (!hasUserPermission('leases', action)) {
      toast.error(`You do not have permission to ${action} leases`);
      return;
    }
    
    try {
      if (currentLease) {
        await updateLease(currentLease.id, leaseData);
        toast.success('Lease updated successfully');
      } else {
        await createLease(leaseData);
        toast.success('Lease created successfully');
      }
      setIsModalOpen(false);
      fetchLeases();
    } catch (error) {
      toast.error(error.response?.data?.message || 'Operation failed');
    }
  };

  // Define table columns
  const columns = [
    { 
      key: 'tenant', 
      label: 'Tenant', 
      sortable: true,
      render: (value, row) => `${row.tenant?.first_name || ''} ${row.tenant?.last_name || ''}`.trim() || 'N/A'
    },
    { 
      key: 'property', 
      label: 'Property', 
      sortable: true,
      render: (value, row) => `${row.property?.name || 'N/A'} (${row.property?.post_code || ''})`
    },
    { 
      key: 'start_date', 
      label: 'Start Date', 
      sortable: true,
      render: (value) => new Date(value).toLocaleDateString()
    },
    { 
      key: 'end_date', 
      label: 'End Date', 
      sortable: true,
      render: (value) => new Date(value).toLocaleDateString()
    },
    { 
      key: 'monthly_rent', 
      label: 'Monthly Rent', 
      sortable: true,
      render: (value) => `£${value?.toFixed(2) || '0.00'}`
    },
    {
      key: 'actions',
      label: 'Actions',
      render: (value, row) => (
        <RoleBasedActions
          resource="leases"
          item={row}
          onView={() => {
            setCurrentLease(row);
            setIsModalOpen(true);
          }}
          onEdit={() => {
            setCurrentLease(row);
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
          <h3>Error Loading Leases</h3>
          <p>{error}</p>
          <Button onClick={fetchLeases}>Try Again</Button>
        </div>
      </div>
    );
  }

  return (
    <div className="page-container">
      <div className="page-header">
        <div>
          <h1>Leases</h1>
          <p className="page-description">Manage property lease agreements</p>
        </div>
        <AdminOnly>
          <Button
            variant="primary"
            onClick={() => {
              setCurrentLease(null);
              setIsModalOpen(true);
            }}
          >
            Add Lease
          </Button>
        </AdminOnly>
        <LandlordOnly>
          <Button
            variant="primary"
            onClick={() => {
              setCurrentLease(null);
              setIsModalOpen(true);
            }}
          >
            Add Lease
          </Button>
        </LandlordOnly>
      </div>

      <DataTable
        data={leases}
        columns={columns}
        loading={isLoading}
        searchable={true}
        searchPlaceholder="Search leases..."
        emptyMessage="No leases found"
      />

      {isModalOpen && (
        <div className="modal-backdrop">
          <div className="modal-container">
            <div className="modal-header">
              <h3>{currentLease ? 'Edit Lease' : 'Create New Lease'}</h3>
              <button 
                className="modal-close"
                onClick={() => setIsModalOpen(false)}
                aria-label="Close modal"
              >
                ×
              </button>
            </div>
            <div className="modal-content">
              <LeaseForm
                initialValues={currentLease || {}}
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

export default Leases;