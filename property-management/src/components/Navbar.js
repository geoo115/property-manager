import React, { useContext, useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { AuthContext } from '../context/AuthContext';
import axiosInstance from '../api/axiosInstance';

const Navbar = () => {
  const { user, logout } = useContext(AuthContext);
  const navigate = useNavigate();
  const [leaseId, setLeaseId] = useState(null);

  const [lease, setLease] = useState(null);

  useEffect(() => {
    if (user?.role === 'tenant') {
      axiosInstance.get('/tenant/leases')
        .then(res => {
          console.log("Tenant leases:", res.data);
          if (res.data.leases.length > 0) {
            const firstLease = res.data.leases[0]; // Assuming one lease per tenant
            setLease(firstLease);
            setLeaseId(firstLease.id);  // Update leaseId here
          }
        })
        .catch(err => console.error("Error fetching leases", err));
    }
  }, [user]);
  
  
  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <nav className="navbar">
      <Link to="/" className="nav-link">Dashboard</Link>

      {user ? (
        <>
          <span className="welcome-message">
            Welcome {user.role} {user.username}
          </span>

          {/* Admin - Full Access */}
          {user.role === "admin" && (
            <>
              <Link to="/admin/properties" className="nav-link">Properties</Link>
              <Link to="/admin/leases" className="nav-link">Leases</Link>
              <Link to="/admin/users" className="nav-link">Users</Link>
              <Link to="/admin/maintenances" className="nav-link">Maintenance</Link>
              <Link to="/admin/accounting/invoices" className="nav-link">Invoices</Link>
              <Link to="/admin/accounting/expenses" className="nav-link">Expenses</Link>
            </>
          )}

          {/* Landlord - View Only */}
          {user.role === "landlord" && (
            <>
              <Link to="/landlord/properties" className="nav-link">My Properties</Link>
              <Link to="/landlord/leases" className="nav-link">My Leases</Link>
              <Link to="/accounting/invoices" className="nav-link">Invoices</Link>
              <Link to="/accounting/expenses" className="nav-link">Expenses</Link>
            </>
          )}

          {/* Tenant - Limited Access */}
          {user.role === "tenant" && leaseId && (
            <>
              <Link to="/tenant/leases" className="nav-link">My Leases</Link>
              <Link to={`/tenant/maintenance/${leaseId}`} className="nav-link">Maintenance</Link>
            </>
          )}

          {/* Maintenance Team */}
          {user.role === "maintenanceTeam" && (
            <Link to="/maintenanceTeam/maintenances" className="nav-link">Maintenance</Link>
          )}

          <button onClick={handleLogout} className="nav-link logout-button">Logout</button>
        </>
      ) : (
        <>
          <Link to="/login" className="nav-link">Login</Link>
          <Link to="/register" className="nav-link">Register</Link>
        </>
      )}
    </nav>
  );
};

export default Navbar;