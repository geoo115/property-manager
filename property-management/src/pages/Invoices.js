import React, { useState, useEffect, useContext, useCallback } from 'react';
import { AuthContext } from '../context/AuthContext';
import axiosInstance from '../api/axiosInstance';
import { TailSpin } from 'react-loader-spinner';
import { toast } from 'react-toastify';
import InvoiceForm from './InvoiceForm';
import './accounting.css';

const Invoices = () => {
  const { user } = useContext(AuthContext);
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

  return (
    <div className="accounting-container">
      <div className="accounting-header">
        <h2>Invoices</h2>
        <button
          className="button button-primary"
          onClick={() => {
            setCurrentInvoice(null);
            setIsFormVisible(true);
          }}
        >
          Add Invoice
        </button>
      </div>

      {isLoading ? (
        <div className="loading-container">
          <TailSpin color="#3B82F6" height={80} width={80} />
        </div>
      ) : error ? (
        <div className="error-message">{error}</div>
      ) : (
        <div className="data-container">
          <table className="accounting-table">
            <thead>
              <tr>
                <th>Amount</th>
                <th>Date</th>
                <th>Status</th>
                <th>Tenant Name</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {invoices.map((invoice) => (
                <tr key={invoice.id}>
                  <td>${invoice.amount}</td>
                  <td>{new Date(invoice.invoice_date).toLocaleDateString()}</td>
                  <td>
                    <span className={`status-badge status-${(invoice.payment_status || 'pending').toLowerCase()}`}>
                      {invoice.payment_status || 'Pending'}
                    </span>
                  </td>
                  <td>{invoice.tenant?.first_name} {invoice.tenant?.last_name}</td>
                  <td>
                    <button
                      className="button button-secondary"
                      onClick={() => {
                        setCurrentInvoice(invoice);
                        setIsFormVisible(true);
                      }}
                    >
                      Edit
                    </button>
                    <button
                      className="button button-danger"
                      onClick={() => handleDelete(invoice.id)}
                      style={{ marginLeft: '10px' }}
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {isFormVisible && (
        <InvoiceForm
          initialValues={currentInvoice || {}}
          onSubmit={handleFormSubmit}
          onCancel={() => {
            setIsFormVisible(false);
            setCurrentInvoice(null);
          }}
        />
      )}
    </div>
  );
};

export default Invoices;