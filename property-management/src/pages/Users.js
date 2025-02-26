import React, { useEffect, useState } from 'react';
import { getUsers, createUser, updateUser, deleteUser } from '../api/users';
import UserForm from './UserForm';
import { TailSpin } from 'react-loader-spinner';
import { toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import './User.css';

const Users = () => {
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
    <div className="user-container">
      <div className="user-header">
        <h2>User Management</h2>
        <button
          className="btn-primary"
          onClick={() => {
            setCurrentUser(null);
            setIsModalOpen(true);
          }}
        >
          Add New User
        </button>
      </div>

      <div className="user-table-container">
        <table className="user-table">
          <thead>
            <tr>
              <th>Username</th>
              <th>Email</th>
              <th>Role</th>
              <th>Phone</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {users.map((user) => (
              <tr key={user.id}>
                <td>{user.username}</td>
                <td>{user.email}</td>
                <td>
                  <span className={`role-badge role-${user.role.toLowerCase().replace(' ', '')}`}>
                    {user.role}
                  </span>
                </td>
                <td>{user.phone}</td>
                <td>
                  <div className="table-actions">
                    <button
                      className="btn-edit"
                      onClick={() => {
                        setCurrentUser(user);
                        setIsModalOpen(true);
                      }}
                    >
                      Edit
                    </button>
                    <button
                      className="btn-delete"
                      onClick={() => handleDelete(user.id)}
                    >
                      Delete
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
        {users.length === 0 && (
          <div className="no-users">No users found</div>
        )}
      </div>

      {isModalOpen && (
        <div className="modal-backdrop">
          <div className="modal-container">
            <h3>{currentUser ? 'Edit User' : 'Create New User'}</h3>
            <UserForm
              initialValues={currentUser || {}}
              onSubmit={handleSubmit}
              onCancel={() => setIsModalOpen(false)}
            />
          </div>
        </div>
      )}
    </div>
  );
};

export default Users;