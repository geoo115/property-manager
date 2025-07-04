import React, { useState, useEffect, useContext, useCallback } from 'react';
import { AuthContext } from '../context/AuthContext';
import axiosInstance from '../api/axiosInstance';
import DataTable from '../components/common/DataTable';
import Button from '../components/common/Button';
import RoleBasedActions from '../components/common/RoleBasedActions';
import { AdminOnly, LandlordOnly } from '../components/common/RoleBasedContent';
import { toast } from 'react-toastify';
import ExpenseForm from './ExpenseForm';

const Expenses = () => {
  const { user, hasUserPermission } = useContext(AuthContext);
  const [expenses, setExpenses] = useState([]);
  const [properties, setProperties] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [isFormVisible, setIsFormVisible] = useState(false);
  const [currentExpense, setCurrentExpense] = useState(null);

  const fetchExpenses = useCallback(async () => {
    setIsLoading(true);
    try {
      let res;
      if (user.role === 'landlord') {
        res = await axiosInstance.get('/landlord/expenses');
      } else if (user.role === 'admin') {
        res = await axiosInstance.get('/admin/accounting/expenses');
      } else {
        res = await axiosInstance.get('/admin/accounting/expenses');
      }
      setExpenses(res.data.expenses || []);
    } catch (err) {
      setError(err.response?.data?.message || 'Failed to fetch expenses');
      toast.error(err.response?.data?.message || 'Failed to fetch expenses');
    } finally {
      setIsLoading(false);
    }
  }, [user.role]);

  const fetchProperties = useCallback(async () => {
    try {
      const endpoint = user.role === 'landlord' ? '/landlord/properties' : '/admin/properties';
      const res = await axiosInstance.get(endpoint);
      setProperties(res.data.properties || []);
    } catch (err) {
      toast.error(err.response?.data?.message || 'Failed to fetch properties');
    }
  }, [user.role]);

  const handleFormSubmit = async (formData) => {
    const action = currentExpense?.id ? 'update' : 'create';
    if (!hasUserPermission('expenses', action)) {
      toast.error(`You do not have permission to ${action} expenses`);
      return;
    }
    
    try {
      console.log('Submitting form data:', formData);
      const baseEndpoint = user.role === 'landlord' ? '/landlord' : '/admin/accounting';
      if (currentExpense?.id) {
        await axiosInstance.put(`${baseEndpoint}/expense/${currentExpense.id}`, formData);
        toast.success('Expense updated successfully');
      } else {
        await axiosInstance.post(`${baseEndpoint}/expense`, formData);
        toast.success('Expense created successfully');
      }
      fetchExpenses();
      setIsFormVisible(false);
    } catch (err) {
      console.error('Submission error:', err.response?.data);
      toast.error(err.response?.data?.message || 'Failed to save expense');
    }
  };

  const handleDelete = async (id) => {
    if (!hasUserPermission('expenses', 'delete')) {
      toast.error('You do not have permission to delete expenses');
      return;
    }
    
    try {
      const endpoint = user.role === 'landlord' ? '/landlord' : '/admin/accounting';
      await axiosInstance.delete(`${endpoint}/expense/${id}`);
      toast.success('Expense deleted successfully');
      fetchExpenses();
    } catch (err) {
      toast.error(err.response?.data?.message || 'Failed to delete expense');
    }
  };

  useEffect(() => {
    fetchExpenses();
    fetchProperties();
  }, [fetchExpenses, fetchProperties]);

  // Define table columns
  const columns = [
    { key: 'description', label: 'Description', sortable: true },
    { 
      key: 'amount', 
      label: 'Amount', 
      sortable: true,
      render: (value) => `$${value || '0'}`
    },
    { 
      key: 'expense_date', 
      label: 'Date', 
      sortable: true,
      render: (value) => new Date(value).toLocaleDateString()
    },
    { 
      key: 'property', 
      label: 'Property', 
      sortable: true,
      render: (value, row) => row.property?.name || 'N/A'
    },
    {
      key: 'actions',
      label: 'Actions',
      render: (value, row) => (
        <RoleBasedActions
          resource="expenses"
          item={row}
          onView={() => {
            setCurrentExpense(row);
            setIsFormVisible(true);
          }}
          onEdit={() => {
            setCurrentExpense(row);
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
          <h3>Error Loading Expenses</h3>
          <p>{error}</p>
          <Button onClick={fetchExpenses}>Try Again</Button>
        </div>
      </div>
    );
  }

  return (
    <div className="page-container">
      <div className="page-header">
        <div>
          <h1>Expenses</h1>
          <p className="page-description">Track and manage property expenses</p>
        </div>
        <AdminOnly>
          <Button
            variant="primary"
            onClick={() => {
              setCurrentExpense(null);
              setIsFormVisible(true);
            }}
          >
            Add Expense
          </Button>
        </AdminOnly>
        <LandlordOnly>
          <Button
            variant="primary"
            onClick={() => {
              setCurrentExpense(null);
              setIsFormVisible(true);
            }}
          >
            Add Expense
          </Button>
        </LandlordOnly>
      </div>

      <DataTable
        data={expenses}
        columns={columns}
        loading={isLoading}
        searchable={true}
        searchPlaceholder="Search expenses..."
        emptyMessage="No expenses found"
      />

      {isFormVisible && (
        <div className="modal-backdrop">
          <div className="modal-container">
            <div className="modal-header">
              <h3>{currentExpense ? 'Edit Expense' : 'Create Expense'}</h3>
              <button 
                className="modal-close"
                onClick={() => setIsFormVisible(false)}
                aria-label="Close modal"
              >
                ×
              </button>
            </div>
            <div className="modal-content">
              <ExpenseForm
                initialValues={currentExpense || {}}
                properties={properties}
                onSubmit={handleFormSubmit}
                onCancel={() => setIsFormVisible(false)}
              />
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Expenses;