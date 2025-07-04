import React, { useContext } from 'react';
import { AuthContext } from '../context/AuthContext';
import { AdminOnly, LandlordOnly, TenantOnly, MaintenanceOnly } from '../components/common/RoleBasedContent';
import RoleBasedActions from '../components/common/RoleBasedActions';
import { USER_ROLES } from '../constants';
import { ROLE_PERMISSIONS } from '../utils/rbac';

const RBACTest = () => {
  const { user, hasUserPermission, hasRole } = useContext(AuthContext);

  if (!user) {
    return (
      <div className="page-container">
        <div className="error-state">
          <h3>Access Denied</h3>
          <p>You must be logged in to view this page.</p>
        </div>
      </div>
    );
  }

  // Sample test data
  const testItem = { id: 1, name: 'Test Item' };

  return (
    <div className="page-container">
      <div className="page-header">
        <div>
          <h1>RBAC Test Page</h1>
          <p className="page-description">Test Role-Based Access Control implementation</p>
        </div>
      </div>

      <div className="rbac-test-content">
        <div className="test-section">
          <h2>Current User Information</h2>
          <div className="user-info">
            <p><strong>Name:</strong> {user.name}</p>
            <p><strong>Role:</strong> <span className={`role-badge role-${user.role?.toLowerCase()}`}>{user.role}</span></p>
            <p><strong>Email:</strong> {user.email}</p>
          </div>
        </div>

        <div className="test-section">
          <h2>Role-Based Content Components</h2>
          
          <AdminOnly>
            <div className="role-content admin-content">
              <h3>🛡️ Admin Only Content</h3>
              <p>This content is only visible to Administrators.</p>
              <p>Admin permissions: Full system access, user management, system-wide reports</p>
            </div>
          </AdminOnly>

          <LandlordOnly>
            <div className="role-content landlord-content">
              <h3>🏠 Landlord Only Content</h3>
              <p>This content is only visible to Landlords.</p>
              <p>Landlord permissions: Property management, tenant relations, financial tracking</p>
            </div>
          </LandlordOnly>

          <TenantOnly>
            <div className="role-content tenant-content">
              <h3>🏘️ Tenant Only Content</h3>
              <p>This content is only visible to Tenants.</p>
              <p>Tenant permissions: View personal information, submit maintenance requests</p>
            </div>
          </TenantOnly>

          <MaintenanceOnly>
            <div className="role-content maintenance-content">
              <h3>🔧 Maintenance Team Only Content</h3>
              <p>This content is only visible to Maintenance Team members.</p>
              <p>Maintenance permissions: View and update maintenance requests</p>
            </div>
          </MaintenanceOnly>
        </div>

        <div className="test-section">
          <h2>Permission Testing</h2>
          <div className="permission-grid">
            {['users', 'properties', 'leases', 'maintenance', 'invoices', 'expenses'].map(resource => (
              <div key={resource} className="permission-test">
                <h4>{resource.charAt(0).toUpperCase() + resource.slice(1)}</h4>
                <div className="permission-actions">
                  {['create', 'read', 'update', 'delete'].map(action => (
                    <div key={action} className="permission-check">
                      <span className={`permission-indicator ${hasUserPermission(resource, action) ? 'allowed' : 'denied'}`}>
                        {hasUserPermission(resource, action) ? '✅' : '❌'}
                      </span>
                      <span className="permission-text">{action}</span>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="test-section">
          <h2>Role-Based Actions Component</h2>
          <div className="actions-test">
            <h4>Properties Actions</h4>
            <RoleBasedActions
              resource="properties"
              item={testItem}
              onView={() => alert('View property')}
              onEdit={() => alert('Edit property')}
              onDelete={() => alert('Delete property')}
            />
            
            <h4>Users Actions</h4>
            <RoleBasedActions
              resource="users"
              item={testItem}
              onView={() => alert('View user')}
              onEdit={() => alert('Edit user')}
              onDelete={() => alert('Delete user')}
            />
            
            <h4>Maintenance Actions</h4>
            <RoleBasedActions
              resource="maintenance"
              item={testItem}
              onView={() => alert('View maintenance')}
              onEdit={() => alert('Edit maintenance')}
              onDelete={() => alert('Delete maintenance')}
              customActions={[
                {
                  name: 'Assign',
                  icon: '👷',
                  onClick: () => alert('Assign maintenance'),
                  variant: 'secondary',
                  size: 'small',
                  permission: 'update'
                },
                {
                  name: 'Complete',
                  icon: '✅',
                  onClick: () => alert('Complete maintenance'),
                  variant: 'success',
                  size: 'small',
                  permission: 'update'
                }
              ]}
            />
          </div>
        </div>

        <div className="test-section">
          <h2>Role Permissions Matrix</h2>
          <div className="permissions-matrix">
            <table className="matrix-table">
              <thead>
                <tr>
                  <th>Resource</th>
                  <th>Admin</th>
                  <th>Landlord</th>
                  <th>Tenant</th>
                  <th>Maintenance</th>
                </tr>
              </thead>
              <tbody>
                {Object.keys(ROLE_PERMISSIONS[USER_ROLES.ADMIN] || {}).map(resource => (
                  <tr key={resource}>
                    <td className="resource-name">{resource}</td>
                    {Object.values(USER_ROLES).map(role => (
                      <td key={role} className="role-permissions">
                        {ROLE_PERMISSIONS[role]?.[resource]?.join(', ') || 'None'}
                      </td>
                    ))}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        <div className="test-section">
          <h2>Navigation Test</h2>
          <p>Check if navigation items are properly filtered based on your role:</p>
          <div className="nav-test">
            <p><strong>Current Role:</strong> {user.role}</p>
            <p><strong>Should see:</strong></p>
            <ul>
              {hasRole(USER_ROLES.ADMIN) && (
                <>
                  <li>All navigation items (Dashboard, Properties, Leases, Users, Maintenance, Invoices, Expenses)</li>
                  <li>User management section</li>
                  <li>System-wide reports</li>
                </>
              )}
              {hasRole(USER_ROLES.LANDLORD) && (
                <>
                  <li>Dashboard, Properties, Leases, Maintenance, Invoices, Expenses</li>
                  <li>Tenant communication features</li>
                  <li>Financial tracking</li>
                </>
              )}
              {hasRole(USER_ROLES.TENANT) && (
                <>
                  <li>Dashboard, Profile, Leases, Maintenance (submit only), Invoices (view only)</li>
                  <li>Personal information management</li>
                  <li>Maintenance request submission</li>
                </>
              )}
              {hasRole(USER_ROLES.MAINTENANCE_TEAM) && (
                <>
                  <li>Dashboard, Maintenance, Reports</li>
                  <li>Maintenance request management</li>
                  <li>Task assignment and completion</li>
                </>
              )}
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
};

export default RBACTest;
