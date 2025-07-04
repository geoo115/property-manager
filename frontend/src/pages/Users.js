import React, { useEffect, useState, useContext } from 'react';
import { AuthContext } from '../context/AuthContext';
import { getUsers, createUser, updateUser, deleteUser } from '../api/users';
import UserForm from './UserForm';
import DataTable from '../components/common/DataTable';
import Button from '../components/common/Button';
import RoleBasedActions from '../components/common/RoleBasedActions';
import { AdminOnly } from '../components/common/RoleBasedContent';
import { toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

const Users = () => {
  const { hasUserPermission } = useContext(AuthContext);
  const [users, setUsers] = useState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [currentUser, setCurrentUser] = useState(null);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    try {
      const data = await getUsers();
      console.log("Fetched users:", data); 
      setUsers(data.users || []);
      setIsLoading(false);
    } catch (error) {
      console.error("Fetch users error:", error); 
      setError('Failed to fetch users');
      setIsLoading(false);
    }
  };
  
  const handleDelete = async (id) => {
    if (!hasUserPermission('users', 'delete')) {
      toast.error('You do not have permission to delete users');
      return;
    }
    
    if (window.confirm('Are you sure you want to delete this user?')) {
      try {
        await deleteUser(id);
        toast.success('User deleted successfully');
        fetchUsers();
      } catch (error) {
        toast.error('Failed to delete user');
      }
    }
  };

  const handleSubmit = async (userData) => {
    const action = currentUser ? 'update' : 'create';
    if (!hasUserPermission('users', action)) {
      toast.error(`You do not have permission to ${action} users`);
      return;
    }
    
    try {
      if (currentUser) {
        await updateUser(currentUser.id, userData);
        toast.success('User updated successfully');
      } else {
        await createUser(userData);
        toast.success('User created successfully');
      }
      setIsModalOpen(false);
      fetchUsers();
    } catch (error) {
      toast.error(error.response?.data?.message || 'Operation failed');
    }
  };

  // Define table columns
  const columns = [
    { key: 'username', label: 'Username', sortable: true },
    { key: 'email', label: 'Email', sortable: true },
    { 
      key: 'role', 
      label: 'Role', 
      sortable: true,
      render: (value) => (
        <span className={`role-badge role-${value?.toLowerCase().replace(' ', '') || 'unknown'}`}>
          {value || 'N/A'}
        </span>
      )
    },
    { key: 'phone', label: 'Phone', sortable: true },
    {
      key: 'actions',
      label: 'Actions',
      render: (value, row) => (
        <RoleBasedActions
          resource="users"
          item={row}
          onView={() => {
            setCurrentUser(row);
            setIsModalOpen(true);
          }}
          onEdit={() => {
            setCurrentUser(row);
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
          <h3>Error Loading Users</h3>
          <p>{error}</p>
          <Button onClick={fetchUsers}>Try Again</Button>
        </div>
      </div>
    );
  }

  return (
    <div className="page-container">
      <div className="page-header">
        <div>
          <h1>Users</h1>
          <p className="page-description">Manage system users and permissions</p>
        </div>
        <AdminOnly>
          <Button
            variant="primary"
            onClick={() => {
              setCurrentUser(null);
              setIsModalOpen(true);
            }}
          >
            Add User
          </Button>
        </AdminOnly>
      </div>

      <DataTable
        data={users}
        columns={columns}
        loading={isLoading}
        searchable={true}
        searchPlaceholder="Search users..."
        emptyMessage="No users found"
      />

      {isModalOpen && (
        <div className="modal-backdrop">
          <div className="modal-container">
            <div className="modal-header">
              <h3>{currentUser ? 'Edit User' : 'Create New User'}</h3>
              <button 
                className="modal-close"
                onClick={() => setIsModalOpen(false)}
                aria-label="Close modal"
              >
                ×
              </button>
            </div>
            <div className="modal-content">
              <UserForm
                initialValues={currentUser || {}}
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

export default Users;