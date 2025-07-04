import React, { useState, useEffect, useContext, useCallback } from 'react';
import { AuthContext } from '../context/AuthContext';
import axiosInstance from '../api/axiosInstance';
import DataTable from '../components/common/DataTable';
import Button from '../components/common/Button';
import RoleBasedActions from '../components/common/RoleBasedActions';
import { AdminOnly, LandlordOnly } from '../components/common/RoleBasedContent';
import { toast } from 'react-toastify';
import InvoiceForm from './InvoiceForm';

const Invoices = () => {
  const { user, hasUserPermission } = useContext(AuthContext);
  const [invoices, setInvoices] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [isFormVisible, setIsFormVisible] = useState(false);
  const [currentInvoice, setCurrentInvoice] = useState(null);

  const fetchInvoices = useCallback(async () => {
    setIsLoading(true);
    try {
      let res;
      if (user.role === 'tenant') {
        res = await axiosInstance.get('/tenant/invoices');
      } else if (user.role === 'landlord') {
        res = await axiosInstance.get('/landlord/invoices');
      } else if (user.role === 'admin') {
        res = await axiosInstance.get('/admin/accounting/invoices');
      } else {
        res = await axiosInstance.get('/invoices');
      }
      setInvoices(res.data.invoices || []);
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to fetch invoices');
      toast.error(err.response?.data?.message || 'Failed to fetch invoices');
    } finally {
      setIsLoading(false);
    }
  }, [user.role]);;

  const handleFormSubmit = async (formData) => {
    const action = currentInvoice?.id ? 'update' : 'create';
    if (!hasUserPermission('invoices', action)) {
      toast.error(`You do not have permission to ${action} invoices`);
      return;
    }
    
    console.log('handleFormSubmit called with:', formData);
    try {
      if (currentInvoice?.id) {
        const response = await axiosInstance.put(`/admin/accounting/invoices/${currentInvoice.id}`, formData);
        console.log('Update response:', response.data);
        toast.success('Invoice updated successfully');
      } else {
        const response = await axiosInstance.post('/admin/accounting/invoices', formData);
        console.log('Create response:', response.data);
        toast.success('Invoice created successfully');
      }
      await fetchInvoices();
      setIsFormVisible(false);
    } catch (err) {
      console.error('Submission error:', err.response?.data);
      toast.error(err.response?.data?.details || err.response?.data?.error || 'Failed to save invoice');
    }
  };

  const handleDelete = async (id) => {
    if (!hasUserPermission('invoices', 'delete')) {
      toast.error('You do not have permission to delete invoices');
      return;
    }
    
    if (!window.confirm('Are you sure you want to delete this invoice?')) return;
    try {
      await axiosInstance.delete(`/admin/accounting/invoices/${id}`);
      toast.success('Invoice deleted successfully');
      await fetchInvoices(); // Refresh the list
    } catch (err) {
      console.error('Delete error:', err.response?.data);
      toast.error(err.response?.data?.details || err.response?.data?.error || 'Failed to delete invoice');
    }
  };

  useEffect(() => {
    fetchInvoices();
  }, [fetchInvoices]);

  // Define table columns
  const columns = [
    { 
      key: 'amount', 
      label: 'Amount', 
      sortable: true,
      render: (value) => `$${value || '0'}`
    },
    { 
      key: 'invoice_date', 
      label: 'Date', 
      sortable: true,
      render: (value) => new Date(value).toLocaleDateString()
    },
    { 
      key: 'payment_status', 
      label: 'Status', 
      sortable: true,
      render: (value) => (
        <span className={`status-badge status-${(value || 'pending').toLowerCase()}`}>
          {value || 'Pending'}
        </span>
      )
    },
    { 
      key: 'tenant', 
      label: 'Tenant', 
      sortable: true,
      render: (value, row) => `${row.tenant?.first_name || ''} ${row.tenant?.last_name || ''}`.trim() || 'N/A'
    },
    {
      key: 'actions',
      label: 'Actions',
      render: (value, row) => (
        <RoleBasedActions
          resource="invoices"
          item={row}
          onView={() => {
            setCurrentInvoice(row);
            setIsFormVisible(true);
          }}
          onEdit={() => {
            setCurrentInvoice(row);
            setIsFormVisible(true);
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
          <h3>Error Loading Invoices</h3>
          <p>{error}</p>
          <Button onClick={fetchInvoices}>Try Again</Button>
        </div>
      </div>
    );
  }

  return (
    <div className="page-container">
      <div className="page-header">
        <div>
          <h1>Invoices</h1>
          <p className="page-description">Manage invoices and billing</p>
        </div>
        <AdminOnly>
          <Button
            variant="primary"
            onClick={() => {
              setCurrentInvoice(null);
              setIsFormVisible(true);
            }}
          >
            Add Invoice
          </Button>
        </AdminOnly>
        <LandlordOnly>
          <Button
            variant="primary"
            onClick={() => {
              setCurrentInvoice(null);
              setIsFormVisible(true);
            }}
          >
            Add Invoice
          </Button>
        </LandlordOnly>
      </div>

      <DataTable
        data={invoices}
        columns={columns}
        loading={isLoading}
        searchable={true}
        searchPlaceholder="Search invoices..."
        emptyMessage="No invoices found"
      />

      {isFormVisible && (
        <div className="modal-backdrop">
          <div className="modal-container">
            <div className="modal-header">
              <h3>{currentInvoice ? 'Edit Invoice' : 'Create Invoice'}</h3>
              <button 
                className="modal-close"
                onClick={() => {
                  setIsFormVisible(false);
                  setCurrentInvoice(null);
                }}
                aria-label="Close modal"
              >
                ×
              </button>
            </div>
            <div className="modal-content">
              <InvoiceForm
                initialValues={currentInvoice || {}}
                onSubmit={handleFormSubmit}
                onCancel={() => {
                  setIsFormVisible(false);
                  setCurrentInvoice(null);
                }}
              />
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Invoices;