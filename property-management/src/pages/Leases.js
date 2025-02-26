import React, { useEffect, useState } from 'react';
import { getLeases, createLease, updateLease, deleteLease } from '../api/leases';
import LeaseForm from './LeaseForm';
import { TailSpin } from 'react-loader-spinner';
import { toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import './Leases.css';

const Leases = () => {
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

  if (isLoading) {
    return (
      <div className="loading-container">
        <TailSpin color="#3B82F6" height={80} width={80} />
      </div>
    );
  }

  if (error) {
    return (
      <div className="error-message">
        {error}
      </div>
    );
  }

  return (
    <div className="leases-container">
      <div className="leases-header">
        <h2>Lease Management</h2>
        <button
          className="btn-primary"
          onClick={() => {
            setCurrentLease(null);
            setIsModalOpen(true);
          }}
        >
          Add New Lease
        </button>
      </div>

      <div className="leases-table-container">
        <table className="leases-table">
          <thead>
            <tr>
              <th>Tenant</th>
              <th>Property</th>
              <th>Start Date</th>
              <th>End Date</th>
              <th>Monthly Rent</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {leases.map((lease) => (
              <tr key={lease.id}>
                <td>{lease.tenant?.first_name} {lease.tenant?.last_name}</td>
                <td>{lease.property?.name} ({lease.property?.post_code})</td>
                <td>{new Date(lease.start_date).toLocaleDateString()}</td>
                <td>{new Date(lease.end_date).toLocaleDateString()}</td>
                <td>£{lease.monthly_rent.toFixed(2)}</td>
                <td>
                  <div className="table-actions">
                    <button
                      className="btn-edit"
                      onClick={() => {
                        setCurrentLease(lease);
                        setIsModalOpen(true);
                      }}
                    >
                      Edit
                    </button>
                    <button
                      className="btn-delete"
                      onClick={() => handleDelete(lease.id)}
                    >
                      Delete
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {leases.length === 0 && (
          <div className="no-leases">No leases found</div>
        )}
      </div>

      {isModalOpen && (
        <div className="modal-backdrop">
          <div className="modal-container">
            <h3>{currentLease ? 'Edit Lease' : 'Create New Lease'}</h3>
            <LeaseForm
              initialValues={currentLease || {}}
              onSubmit={handleSubmit}
              onCancel={() => setIsModalOpen(false)}
            />
          </div>
        </div>
      )}
    </div>
  );
};

export default Leases;