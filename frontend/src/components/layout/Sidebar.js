import React, { useState, useMemo } from 'react';
import { Link, useLocation } from 'react-router-dom';
import useAuth from '../../hooks/useAuth';
import { USER_ROLES } from '../../constants';
import { filterNavigation } from '../../utils/rbac';
import './Sidebar.css';

const Sidebar = () => {
  const { user, hasUserPermission, canUserAccess } = useAuth();
  const location = useLocation();
  const [isCollapsed, setIsCollapsed] = useState(false);

  const navigation = useMemo(() => {
    if (!user) return [];
    
    return [
      {
        name: 'Dashboard',
        href: '/dashboard',
        icon: (
          <svg className="nav-icon" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" d="M2.25 12l8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25" />
          </svg>
        ),
        roles: Object.values(USER_ROLES),
        description: 'Overview and analytics',
      },
    {
      name: 'Properties',
      href: user.role === 'admin' ? '/admin/properties' : '/landlord/properties',
      icon: (
        <svg className="nav-icon" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" d="M8.25 21v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21m0 0h4.5V9.375c0-.621.504-1.125 1.125-1.125h.75c.621 0 1.125.504 1.125 1.125v11.25c-1.5 0-3-1.5-3-3.375S15.75 14.25 17.25 14.25" />
        </svg>
      ),
      roles: [USER_ROLES.ADMIN, USER_ROLES.LANDLORD],
      description: user.role === 'admin' ? 'All properties management' : 'Your properties',
    },
    {
      name: 'Leases',
      href: user.role === 'admin' ? '/admin/leases' : 
            user.role === 'landlord' ? '/landlord/leases' : '/tenant/leases',
      icon: (
        <svg className="nav-icon" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
        </svg>
      ),
      roles: [USER_ROLES.ADMIN, USER_ROLES.LANDLORD, USER_ROLES.TENANT],
      description: user.role === 'admin' ? 'All lease agreements' : 
                   user.role === 'landlord' ? 'Your lease agreements' : 'Your lease information',
    },
    {
      name: 'Maintenance',
      href: user.role === 'admin' ? '/admin/maintenances' : 
            user.role === 'tenant' ? '/tenant/maintenance' : '/maintenanceTeam/maintenances',
      icon: (
        <svg className="nav-icon" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" d="M11.42 15.17L17.25 21A2.652 2.652 0 0021 17.25l-5.877-5.877M11.42 15.17l2.496-3.03c.317-.384.74-.626 1.208-.766M11.42 15.17l-4.655 5.653a2.548 2.548 0 11-3.586-3.586l6.837-5.63m5.108-.233c.55-.164 1.163-.188 1.743-.14a4.5 4.5 0 004.486-6.336l-3.276 3.277a3.004 3.004 0 01-2.25-2.25l3.276-3.276a4.5 4.5 0 00-6.336 4.486c.091 1.076-.071 2.264-.904 2.95l-.102.085m-1.745 1.437L5.909 7.5H4.5L2.25 3.75l1.5-1.5L7.5 4.5v1.409l4.26 4.26m-1.745 1.437l1.745-1.437m6.615 8.206L15.75 15.75M4.867 19.125h.008v.008h-.008v-.008z" />
        </svg>
      ),
      roles: Object.values(USER_ROLES),
      description: user.role === 'admin' ? 'All maintenance requests' : 
                   user.role === 'tenant' ? 'Submit & view requests' : 
                   user.role === 'maintenanceTeam' ? 'Assigned requests' : 'Property maintenance',
    },
    {
      name: 'Invoices',
      href: user.role === 'admin' ? '/admin/accounting/invoices' : '/accounting/invoices',
      icon: (
        <svg className="nav-icon" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25M9 16.5v.75m3-3v3M15 12v5.25m-4.5-15H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
        </svg>
      ),
      roles: [USER_ROLES.ADMIN, USER_ROLES.LANDLORD, USER_ROLES.TENANT],
      description: user.role === 'admin' ? 'All invoices' : 
                   user.role === 'landlord' ? 'Property invoices' : 'Your invoices',
    },
    {
      name: 'Expenses',
      href: user.role === 'admin' ? '/admin/accounting/expenses' : '/accounting/expenses',
      icon: (
        <svg className="nav-icon" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" d="M2.25 18.75a60.07 60.07 0 0115.797 2.101c.727.198 1.453-.342 1.453-1.096V18.75M3.75 4.5v.75A.75.75 0 013 6h-.75m0 0v-.375c0-.621.504-1.125 1.125-1.125H15.75c.621 0 1.125.504 1.125 1.125v.375m-13.5 0h12m-12 0v5.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V6h-12z" />
        </svg>
      ),
      roles: [USER_ROLES.ADMIN, USER_ROLES.LANDLORD],
      description: user.role === 'admin' ? 'All expenses' : 'Property expenses',
    },
    {
      name: 'Users',
      href: '/admin/users',
      icon: (
        <svg className="nav-icon" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z" />
        </svg>
      ),
      roles: [USER_ROLES.ADMIN],
      description: 'User management',
    },
      {
        name: 'Reports',
        href: '/admin/reports',
        icon: (
          <svg className="nav-icon" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z" />
          </svg>
        ),
        roles: [USER_ROLES.ADMIN],
        description: 'System-wide reports',
      },
    ];
  }, [user]);

  const filteredNavigation = useMemo(() => 
    filterNavigation(navigation, user?.role), [navigation, user?.role]
  );

  const isActive = (href) => {
    return location.pathname === href || location.pathname.startsWith(href + '/');
  };

  if (!user) return null;

  return (
    <div className={`sidebar ${isCollapsed ? 'sidebar--collapsed' : ''}`}>
      <div className="sidebar__header">
        <div className="sidebar__brand">
          <div className="brand-icon">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 2L2 7v10c0 5.55 3.84 9.74 9 11 5.16-1.26 9-5.45 9-11V7l-10-5z"/>
            </svg>
          </div>
          {!isCollapsed && (
            <div className="brand-text">
              <h1 className="brand-title">PropertyHub</h1>
              <p className="brand-subtitle">Management System</p>
            </div>
          )}
        </div>
        <button
          className="sidebar__toggle"
          onClick={() => setIsCollapsed(!isCollapsed)}
          aria-label={isCollapsed ? 'Expand sidebar' : 'Collapse sidebar'}
        >
          <svg className="toggle-icon" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" d={isCollapsed ? "M8.25 4.5l7.5 7.5-7.5 7.5" : "M15.75 19.5L8.25 12l7.5-7.5"} />
          </svg>
        </button>
      </div>

      <nav className="sidebar__nav">
        <ul className="nav-list">
          {filteredNavigation.map((item) => (
            <li key={item.name} className="nav-item">
              <Link
                to={item.href}
                className={`nav-link ${isActive(item.href) ? 'nav-link--active' : ''}`}
              >
                <span className="nav-link__icon">
                  {item.icon}
                </span>
                {!isCollapsed && (
                  <span className="nav-link__text">{item.name}</span>
                )}
              </Link>
            </li>
          ))}
        </ul>
      </nav>

      <div className="sidebar__footer">
        <div className="user-info">
          <div className="user-avatar">
            <span className="avatar-text">
              {user.firstName?.charAt(0)?.toUpperCase() || user.username?.charAt(0)?.toUpperCase() || 'U'}
            </span>
          </div>
          {!isCollapsed && (
            <div className="user-details">
              <p className="user-name">{user.firstName || user.username}</p>
              <p className="user-role">{user.role}</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default Sidebar;
