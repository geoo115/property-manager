import React, { useState, useEffect, useContext, useCallback } from 'react';
import { AuthContext } from '../context/AuthContext';
import axiosInstance from '../api/axiosInstance';
import { TailSpin } from 'react-loader-spinner';
import { toast } from 'react-toastify';
import ExpenseForm from './ExpenseForm';

const Expenses = () => {
  const { user } = useContext(AuthContext);
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

  return (
    <div className="accounting-container">
      <div className="accounting-header">
        <h2>Expenses</h2>
        <button
          className="button button-primary"
          onClick={() => {
            setCurrentExpense(null);
            setIsFormVisible(true);
          }}
        >
          Add Expense
        </button>
      </div>

      {isLoading ? (
        <div className="loading-container">
          <TailSpin color="#3B82F6" height={80} width={80} />
        </div>
      ) : error ? (
        <div className="error-message">{error}</div>
      ) : (
        <div className="accounting-table-container">
          <table className="accounting-table">
            <thead>
              <tr>
                <th>Description</th>
                <th>Amount</th>
                <th>Date</th>
                <th>Property</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {expenses.map((expense) => (
                <tr key={expense.id}>
                  <td>{expense.description}</td>
                  <td>{expense.amount}</td>
                  <td>{new Date(expense.expense_date).toLocaleDateString()}</td>
                  <td>{expense.property?.name || 'N/A'}</td>
                  <td>
                    <button
                      className="button button-edit"
                      onClick={() => {
                        setCurrentExpense(expense);
                        setIsFormVisible(true);
                      }}
                    >
                      Edit
                    </button>
                    <button
                      className="button button-delete"
                      onClick={() => handleDelete(expense.id)}
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
        <ExpenseForm
          initialValues={currentExpense || {}}
          properties={properties}
          onSubmit={handleFormSubmit}
          onCancel={() => setIsFormVisible(false)}
        />
      )}
    </div>
  );
};

export default Expenses;