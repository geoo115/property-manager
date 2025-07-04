import React, { useContext } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { AuthContext } from '../context/AuthContext';
import Button from './common/Button';

const Navbar = () => {
  const { user, logout } = useContext(AuthContext);
  const location = useLocation();

  const isActive = (path) => {
    return location.pathname === path;
  };

  const getNavLinks = () => {
    if (!user) return [];

    const baseLinks = [
      { to: '/dashboard', label: 'Dashboard', icon: '📊' },
      { to: '/properties', label: 'Properties', icon: '🏠' },
      { to: '/leases', label: 'Leases', icon: '📋' },
      { to: '/maintenance', label: 'Maintenance', icon: '🔧' },
      { to: '/accounting', label: 'Accounting', icon: '💰' },
    ];

    // Add role-specific links
    if (user.role === 'admin') {
      baseLinks.push({ to: '/users', label: 'Users', icon: '👥' });
    }

    return baseLinks;
  };

  const handleLogout = () => {
    logout();
  };

  if (!user) {
    return null;
  }

  const navLinks = getNavLinks();

  return (
    <nav className="navbar">
      <Link to="/dashboard" className="nav-brand">
        <span role="img" aria-label="Property Manager">🏢</span>
        <span>Property Manager</span>
      </Link>

      <div className="nav-links">
        {navLinks.map((link) => (
          <Link
            key={link.to}
            to={link.to}
            className={`nav-link ${isActive(link.to) ? 'active' : ''}`}
            aria-current={isActive(link.to) ? 'page' : undefined}
          >
            <span role="img" aria-label={link.label} className="nav-icon">
              {link.icon}
            </span>
            <span className="nav-text">{link.label}</span>
          </Link>
        ))}
      </div>

      <div className="nav-user">
        <span className="welcome-message">
          Welcome, {user.name} ({user.role})
        </span>
        <Button
          variant="danger"
          size="small"
          onClick={handleLogout}
          className="logout-button"
        >
          Logout
        </Button>
      </div>
    </nav>
  );
};

export default Navbar;