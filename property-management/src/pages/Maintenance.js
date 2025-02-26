import React, { useEffect, useState, useContext, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { AuthContext } from '../context/AuthContext';
import { getMaintenances, createMaintenance, updateMaintenance, deleteMaintenance } from '../api/maintenance';
import MaintenanceForm from './MaintenanceForm';
import { TailSpin } from 'react-loader-spinner';
import { toast } from 'react-toastify';
import './maintenance.css';

const Maintenance = () => {
  const { leaseId } = useParams();
  const { user } = useContext(AuthContext);
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
    if (!window.confirm('Delete this request?')) return;
    try {
      await deleteMaintenance(id);
      await fetchMaintenances();
      toast.success('Request deleted');
    } catch (error) {
      toast.error('Deletion failed');
    }
  };

  if (isLoading) {
    return <div className="loading-container"><TailSpin color="#3B82F6" height={80} width={80} /></div>;
  }

  return (
    <div className="maintenance-container">
      <div className="maintenance-header">
        <h2>{user.role === 'tenant' ? 'Property Maintenance Requests' : 'Maintenance Requests'}</h2>
        {['admin', 'tenant', 'landlord'].includes(user.role) && (
          <button className="btn-primary" onClick={() => setIsModalOpen(true)}>
            New Request
          </button>
        )}
      </div>

      <div className="maintenance-table-container">
        <table className="maintenance-table">
          <thead>
            <tr>
              <th>Property</th>
              <th>Requested By</th>
              <th>Description</th>
              <th>Date</th>
              <th>Status</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {maintenances.map((req) => (
              <tr key={req.id}>
                <td>{req.property?.name || req.property?.address || 'N/A'}</td>
                <td>{req.reporter?.first_name} {req.reporter?.last_name}</td>
                <td>{req.description}</td>
                <td>{new Date(req.requested_at).toLocaleDateString()}</td>
                <td>
                  <span className={`status-badge status-${req.status.toLowerCase()}`}>
                    {req.status}
                  </span>
                </td>
                <td>
                  {['admin', 'maintenanceTeam'].includes(user.role) && (
                    <>
                      <button className="btn-edit" onClick={() => {
                        setCurrentMaintenance(req);
                        setIsModalOpen(true);
                      }}>
                        Edit
                      </button>
                      <button className="btn-delete" onClick={() => handleDelete(req.id)}>
                        Delete
                      </button>
                    </>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {isModalOpen && (
        <div className="modal-backdrop">
          <div className="modal-container">
            <h3>{currentMaintenance ? 'Edit Request' : 'New Request'}</h3>
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
      )}
    </div>
  );
};

export default Maintenance;