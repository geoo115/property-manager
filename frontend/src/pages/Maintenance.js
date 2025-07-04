import React, { useEffect, useState, useContext, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { AuthContext } from '../context/AuthContext';
import { getMaintenances, createMaintenance, updateMaintenance, deleteMaintenance } from '../api/maintenance';
import MaintenanceForm from './MaintenanceForm';
import DataTable from '../components/common/DataTable';
import Button from '../components/common/Button';
import RoleBasedActions from '../components/common/RoleBasedActions';
import { TenantOnly, AdminOnly, LandlordOnly, MaintenanceOnly } from '../components/common/RoleBasedContent';
import { toast } from 'react-toastify';

const Maintenance = () => {
  const { leaseId } = useParams();
  const { user, hasUserPermission } = useContext(AuthContext);
  const navigate = useNavigate();
  const [maintenances, setMaintenances] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [currentMaintenance, setCurrentMaintenance] = useState(null);

  const fetchMaintenances = useCallback(async () => {
    setIsLoading(true);
    try {
      const data = await getMaintenances(
        user.role === 'tenant' ? leaseId : null,
        user.role === 'landlord' ? user.property_id : null
      );
      setMaintenances(data.maintenances || []);
    } catch (error) {
      toast.error('Failed to load maintenance requests');
    } finally {
      setIsLoading(false);
    }
  }, [leaseId, user]);

  useEffect(() => {
    if (user.role === 'tenant' && !leaseId) {
      navigate('/tenant/leases');
      return;
    }
    fetchMaintenances();
  }, [fetchMaintenances, user.role, leaseId, navigate]);

  const handleSubmit = async (formData) => {
    const action = currentMaintenance ? 'update' : 'create';
    if (!hasUserPermission('maintenance', action)) {
      toast.error(`You do not have permission to ${action} maintenance requests`);
      return;
    }
    
    try {
      if (currentMaintenance) {
        await updateMaintenance(currentMaintenance.id, formData);
        toast.success('Request updated successfully');
      } else {
        switch(user.role) {
          case 'tenant':
            await createMaintenance(leaseId, formData);
            break;
          case 'admin':
            await createMaintenance(null, formData, formData.property_id);
            break;
          case 'landlord':
            await createMaintenance(null, formData, user.property_id);
            break;
          default:
            throw new Error('Unauthorized operation');
        }
        toast.success('Request created successfully');
      }
      setIsModalOpen(false);
      setCurrentMaintenance(null);
      await fetchMaintenances();
    } catch (error) {
      toast.error(error.message || 'Operation failed');
    }
  };

  const handleDelete = async (id) => {
    if (!hasUserPermission('maintenance', 'delete')) {
      toast.error('You do not have permission to delete maintenance requests');
      return;
    }
    
    if (!window.confirm('Delete this request?')) return;
    try {
      await deleteMaintenance(id);
      await fetchMaintenances();
      toast.success('Request deleted');
    } catch (error) {
      toast.error('Deletion failed');
    }
  };

  // Define table columns
  const columns = [
    { 
      key: 'property', 
      label: 'Property', 
      sortable: true,
      render: (value, row) => row.property?.name || row.property?.address || 'N/A'
    },
    { 
      key: 'reporter', 
      label: 'Requested By', 
      sortable: true,
      render: (value, row) => `${row.reporter?.first_name || ''} ${row.reporter?.last_name || ''}`.trim() || 'N/A'
    },
    { key: 'description', label: 'Description', sortable: true },
    { 
      key: 'requested_at', 
      label: 'Date', 
      sortable: true,
      render: (value) => new Date(value).toLocaleDateString()
    },
    { 
      key: 'status', 
      label: 'Status', 
      sortable: true,
      render: (value) => (
        <span className={`status-badge status-${value?.toLowerCase() || 'unknown'}`}>
          {value || 'Unknown'}
        </span>
      )
    },
    {
      key: 'actions',
      label: 'Actions',
      render: (value, row) => (
        <RoleBasedActions
          resource="maintenance"
          item={row}
          onView={() => {
            setCurrentMaintenance(row);
            setIsModalOpen(true);
          }}
          onEdit={() => {
            setCurrentMaintenance(row);
            setIsModalOpen(true);
          }}
          onDelete={() => handleDelete(row.id)}
          customActions={[
            {
              name: 'Assign',
              icon: '👷',
              onClick: () => handleAssign(row.id),
              variant: 'secondary',
              size: 'small',
              permission: 'update'
            },
            {
              name: 'Complete',
              icon: '✅',
              onClick: () => handleComplete(row.id),
              variant: 'success',
              size: 'small',
              permission: 'update'
            }
          ]}
        />
      )
    }
  ];

  // Handler functions for custom actions
  const handleAssign = async (id) => {
    // Implementation for assigning maintenance request
    toast.info('Assign functionality to be implemented');
  };

  const handleComplete = async (id) => {
    // Implementation for completing maintenance request
    toast.info('Complete functionality to be implemented');
  };

  return (
    <div className="page-container">
      <div className="page-header">
        <div>
          <h1>Maintenance Requests</h1>
          <p className="page-description">
            {user.role === 'tenant' ? 'Track your property maintenance requests' : 'Manage maintenance requests'}
          </p>
        </div>
        <TenantOnly>
          <Button
            variant="primary"
            onClick={() => setIsModalOpen(true)}
          >
            Submit Request
          </Button>
        </TenantOnly>
        <AdminOnly>
          <Button
            variant="primary"
            onClick={() => setIsModalOpen(true)}
          >
            New Request
          </Button>
        </AdminOnly>
        <LandlordOnly>
          <Button
            variant="primary"
            onClick={() => setIsModalOpen(true)}
          >
            New Request
          </Button>
        </LandlordOnly>
      </div>

      <DataTable
        data={maintenances}
        columns={columns}
        loading={isLoading}
        searchable={true}
        searchPlaceholder="Search maintenance requests..."
        emptyMessage="No maintenance requests found"
      />

      {isModalOpen && (
        <div className="modal-backdrop">
          <div className="modal-container">
            <div className="modal-header">
              <h3>{currentMaintenance ? 'Edit Request' : 'New Request'}</h3>
              <button 
                className="modal-close"
                onClick={() => {
                  setCurrentMaintenance(null);
                  setIsModalOpen(false);
                }}
                aria-label="Close modal"
              >
                ×
              </button>
            </div>
            <div className="modal-content">
              <MaintenanceForm
                initialValues={currentMaintenance || {}}
                onSubmit={handleSubmit}
                onCancel={() => {
                  setCurrentMaintenance(null);
                  setIsModalOpen(false);
                }}
                leaseId={leaseId}
                userRole={user.role}
              />
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Maintenance;