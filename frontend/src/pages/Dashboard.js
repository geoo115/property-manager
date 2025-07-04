import React, { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';

import useAuth from '../hooks/useAuth';
import useApi from '../hooks/useApi';
import { USER_ROLES, ROUTES } from '../constants';
import axiosInstance from '../api/axiosInstance';

import LoadingSpinner from '../components/common/LoadingSpinner';
import Alert from '../components/common/Alert';
import Button from '../components/common/Button';
import { RoleSwitch, AdminOnly, LandlordOnly, TenantOnly, MaintenanceOnly } from '../components/common/RoleBasedContent';
import './Dashboard.css';

const Dashboard = () => {
  const navigate = useNavigate();
  const { user, hasRole, getUserDashboardConfig } = useAuth();
  const { execute, loading } = useApi();
  const [stats, setStats] = useState({});
  const [recentActivities, setRecentActivities] = useState([]);
  const [error, setError] = useState(null);

  const dashboardConfig = getUserDashboardConfig();

  const fetchDashboardData = useCallback(async () => {
    try {
      setError(null);
      
      // Use role-specific endpoints based on user role
      let statsEndpoint = '/admin/dashboard/stats';
      let activitiesEndpoint = '/admin/dashboard/activities';
      
      switch (user?.role) {
        case USER_ROLES.ADMIN:
          statsEndpoint = '/admin/dashboard/stats';
          activitiesEndpoint = '/admin/dashboard/activities';
          break;
        case USER_ROLES.LANDLORD:
          statsEndpoint = '/landlord/dashboard/stats';
          activitiesEndpoint = '/landlord/dashboard/activities';
          break;
        case USER_ROLES.TENANT:
          statsEndpoint = '/tenant/dashboard/stats';
          activitiesEndpoint = '/tenant/dashboard/activities';
          break;
        case USER_ROLES.MAINTENANCE_TEAM:
          statsEndpoint = '/maintenanceTeam/dashboard/stats';
          activitiesEndpoint = '/maintenanceTeam/dashboard/activities';
          break;
        default:
          statsEndpoint = '/admin/dashboard/stats';
          activitiesEndpoint = '/admin/dashboard/activities';
      }

      await execute(
        () => axiosInstance.get(statsEndpoint),
        {
          onSuccess: (response) => {
            setStats(response.data);
          },
          onError: (err) => {
            console.error('Error fetching stats:', err);
            setError('Failed to load dashboard statistics');
          },
          showLoading: false,
        }
      );

      // Fetch recent activities
      await execute(
        () => axiosInstance.get(activitiesEndpoint),
        {
          onSuccess: (response) => {
            setRecentActivities(response.data.activities || []);
          },
          onError: (err) => {
            console.error('Error fetching activities:', err);
          },
          showLoading: false,
        }
      );
    } catch (err) {
      setError('Failed to load dashboard data');
      toast.error('Failed to load dashboard data');
    }
  }, [execute]);

  useEffect(() => {
    if (user) {
      fetchDashboardData();
    }
  }, [user, fetchDashboardData]);

  const renderStatsCards = () => {
    const cards = [];

    if (hasRole(USER_ROLES.ADMIN)) {
      // Primary admin statistics
      cards.push(
        { title: 'Total Properties', value: stats.totalProperties || 0, color: 'primary', icon: 'home', trend: '+12%' },
        { title: 'Active Leases', value: stats.activeLeases || 0, color: 'success', icon: 'file-contract', trend: '+8%' },
        { title: 'Total Users', value: stats.totalUsers || 0, color: 'info', icon: 'users', trend: '+5%' },
        { title: 'Total Revenue', value: `$${(stats.totalRevenue || 0).toLocaleString()}`, color: 'success', icon: 'dollar-sign', trend: '+15%' },
      );

      // Secondary admin statistics
      if (stats.totalTenants !== undefined || stats.totalLandlords !== undefined) {
        cards.push(
          { title: 'Tenants', value: stats.totalTenants || 0, color: 'info', icon: 'user-friends' },
          { title: 'Landlords', value: stats.totalLandlords || 0, color: 'purple', icon: 'user-tie' },
        );
      }

      // Operational metrics
      cards.push(
        { title: 'Pending Maintenance', value: stats.pendingMaintenance || 0, color: 'warning', icon: 'tools', alert: stats.pendingMaintenance > 5 },
        { title: 'Completed Maintenance', value: stats.completedMaintenance || 0, color: 'success', icon: 'check-circle' },
      );

      // Financial metrics
      if (stats.totalExpenses !== undefined) {
        cards.push(
          { title: 'Total Expenses', value: `$${(stats.totalExpenses || 0).toLocaleString()}`, color: 'danger', icon: 'credit-card', trend: '-3%' },
          { title: 'Monthly Expenses', value: `$${(stats.monthlyExpenses || 0).toLocaleString()}`, color: 'warning', icon: 'calendar-alt' },
        );
      }

      // Performance metrics
      if (stats.occupancyRate !== undefined) {
        cards.push(
          { title: 'Occupancy Rate', value: `${(stats.occupancyRate || 0).toFixed(1)}%`, color: 'info', icon: 'chart-pie', trend: '+2%' },
        );
      }
    } else if (hasRole(USER_ROLES.LANDLORD)) {
      cards.push(
        { title: 'My Properties', value: stats.myProperties || 0, color: 'primary', icon: 'home' },
        { title: 'Active Leases', value: stats.activeLeases || 0, color: 'success', icon: 'file-contract' },
        { title: 'Monthly Revenue', value: `$${(stats.monthlyRevenue || 0).toLocaleString()}`, color: 'success', icon: 'dollar-sign' },
        { title: 'Maintenance Requests', value: stats.maintenanceRequests || 0, color: 'warning', icon: 'tools' }
      );
    } else if (hasRole(USER_ROLES.TENANT)) {
      cards.push(
        { title: 'My Leases', value: stats.myLeases || 0, color: 'primary', icon: 'file-contract' },
        { title: 'Outstanding Invoices', value: stats.outstandingInvoices || 0, color: 'danger', icon: 'file-invoice' },
        { title: 'Maintenance Requests', value: stats.maintenanceRequests || 0, color: 'warning', icon: 'tools' },
        { title: 'Total Paid', value: `$${(stats.totalPaid || 0).toLocaleString()}`, color: 'success', icon: 'check-circle' }
      );
    }

    return cards.map((card, index) => (
      <div key={index} className={`stats-card stats-card--${card.color} ${card.alert ? 'stats-card--alert' : ''}`}>
        <div className="stats-card__header">
          <div className="stats-card__icon">
            <svg className="icon" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
              {getIconPath(card.icon)}
            </svg>
          </div>
          {card.trend && (
            <div className={`stats-card__trend ${card.trend.startsWith('+') ? 'trend--positive' : 'trend--negative'}`}>
              <svg className="trend-icon" fill="none" viewBox="0 0 24 24" strokeWidth={2} stroke="currentColor">
                {card.trend.startsWith('+') ? (
                  <path strokeLinecap="round" strokeLinejoin="round" d="M2.25 18L9 11.25l4.306 4.307a11.95 11.95 0 015.814-5.519l2.74-1.22m0 0l-5.94-2.28m5.94 2.28l-2.28 5.94" />
                ) : (
                  <path strokeLinecap="round" strokeLinejoin="round" d="M2.25 6L9 12.75l4.286-4.286a11.948 11.948 0 014.306 6.43l.776 2.898m0 0l3.182-5.511m-3.182 5.511l-5.511-3.182" />
                )}
              </svg>
              <span>{card.trend}</span>
            </div>
          )}
        </div>
        
        <div className="stats-card__content">
          <h3 className="stats-card__title">{card.title}</h3>
          <div className="stats-card__value">{card.value}</div>
        </div>

        <div className="stats-card__footer">
          <div className="stats-card__label">
            {card.alert ? 'Needs attention' : 'Last updated'}
          </div>
        </div>
      </div>
    ));
  };

  const getIconPath = (iconName) => {
    const icons = {
      'home': <path strokeLinecap="round" strokeLinejoin="round" d="m2.25 12 8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25" />,
      'file-contract': <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />,
      'users': <path strokeLinecap="round" strokeLinejoin="round" d="M15 19.128a9.38 9.38 0 002.625.372 9.337 9.337 0 004.121-.952 4.125 4.125 0 00-7.533-2.493M15 19.128v-.003c0-1.113-.285-2.16-.786-3.07M15 19.128v.106A12.318 12.318 0 018.624 21c-2.331 0-4.512-.645-6.374-1.766l-.001-.109a6.375 6.375 0 0111.964-3.07M12 6.375a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zm8.25 2.25a2.625 2.625 0 11-5.25 0 2.625 2.625 0 015.25 0z" />,
      'dollar-sign': <path strokeLinecap="round" strokeLinejoin="round" d="M12 6v12m-3-2.818l.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />,
      'user-friends': <path strokeLinecap="round" strokeLinejoin="round" d="M18 18.72a9.094 9.094 0 003.741-.479 3 3 0 00-4.682-2.72m.94 3.198l.001.031c0 .225-.012.447-.037.666A11.944 11.944 0 0112 21c-2.17 0-4.207-.576-5.963-1.584A6.062 6.062 0 016 18.719m12 0a5.971 5.971 0 00-.941-3.197m0 0A5.995 5.995 0 0012 12.75a5.995 5.995 0 00-5.058 2.772m0 0a3 3 0 00-4.681 2.72 8.986 8.986 0 003.74.477m.94-3.197a5.971 5.971 0 00-.94 3.197M15 6.75a3 3 0 11-6 0 3 3 0 016 0zm6 3a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0zm-13.5 0a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z" />,
      'user-tie': <path strokeLinecap="round" strokeLinejoin="round" d="M20.25 14.15v4.25c0 1.094-.787 2.036-1.872 2.18-2.087.277-4.216.42-6.378.42s-4.291-.143-6.378-.42c-1.085-.144-1.872-1.086-1.872-2.18v-4.25m16.5 0a2.18 2.18 0 00.75-1.661V8.706c0-1.081-.768-2.015-1.837-2.175a48.114 48.114 0 00-3.413-.387m4.5 8.006c-.194.165-.42.295-.673.38A23.978 23.978 0 0112 15.75c-2.648 0-5.195-.429-7.577-1.22a2.016 2.016 0 01-.673-.38m0 0A2.18 2.18 0 013 12.489V8.706c0-1.081.768-2.015 1.837-2.175a48.111 48.111 0 013.413-.387m7.5 0V5.25A2.25 2.25 0 0013.5 3h-3a2.25 2.25 0 00-2.25 2.25v.894m7.5 0a48.667 48.667 0 00-7.5 0M12 12.75h.008v.008H12v-.008z" />,
      'tools': <path strokeLinecap="round" strokeLinejoin="round" d="M11.42 15.17L17.25 21A2.652 2.652 0 0021 17.25l-5.877-5.877M11.42 15.17l2.496-3.03c.317-.384.74-.626 1.208-.766M11.42 15.17l-4.655 5.653a2.548 2.548 0 11-3.586-3.586l6.837-5.63m5.108-.233c.55-.164 1.163-.188 1.743-.14a4.5 4.5 0 004.486-6.336l-3.276 3.277a3.004 3.004 0 01-2.25-2.25l3.276-3.276a4.5 4.5 0 00-6.336 4.486c.091 1.076-.071 2.264-.904 2.95l-.102.085m-1.745 1.437L5.909 7.5H4.5L2.25 3.75l1.5-1.5L7.5 4.5v1.409l4.26 4.26m-1.745 1.437l1.745-1.437m6.615 8.206L15.75 15.75M4.867 19.125h.008v.008h-.008v-.008z" />,
      'check-circle': <path strokeLinecap="round" strokeLinejoin="round" d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />,
      'credit-card': <path strokeLinecap="round" strokeLinejoin="round" d="M2.25 8.25h19.5M2.25 9h19.5m-16.5 5.25h6m-6 2.25h3m-3.75 3h15a2.25 2.25 0 002.25-2.25V6.75A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25v10.5A2.25 2.25 0 004.5 19.5z" />,
      'calendar-alt': <path strokeLinecap="round" strokeLinejoin="round" d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 012.25-2.25h13.5A2.25 2.25 0 0121 7.5v11.25m-18 0A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75m-18 0v-7.5A2.25 2.25 0 015.25 9h13.5A2.25 2.25 0 0121 11.25v7.5" />,
      'chart-pie': <path strokeLinecap="round" strokeLinejoin="round" d="M10.5 6a7.5 7.5 0 107.5 7.5h-7.5V6z M13.5 10.5H21A7.5 7.5 0 0013.5 3v7.5z" />,
      'file-invoice': <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
    };
    return icons[iconName] || icons['home'];
  };

  const formatTimestamp = (timestamp) => {
    try {
      const date = new Date(timestamp);
      const now = new Date();
      const diffMs = now - date;
      const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
      const diffDays = Math.floor(diffHours / 24);

      if (diffHours < 1) {
        return 'Just now';
      } else if (diffHours < 24) {
        return `${diffHours} hour${diffHours > 1 ? 's' : ''} ago`;
      } else if (diffDays < 7) {
        return `${diffDays} day${diffDays > 1 ? 's' : ''} ago`;
      } else {
        return date.toLocaleDateString();
      }
    } catch (error) {
      return 'Recently';
    }
  };

  const getActivityColor = (type) => {
    const colors = {
      'property': 'blue',
      'lease': 'green',
      'maintenance': 'orange',
      'user': 'purple',
      'invoice': 'blue',
      'expense': 'red'
    };
    return colors[type] || 'gray';
  };

  const renderRecentActivities = () => {
    if (!recentActivities.length) {
      return (
        <div className="activities-section">
          <div className="section-header">
            <h3 className="section-title">Recent Activities</h3>
          </div>
          <div className="empty-state">
            <svg className="empty-icon" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" d="M9.568 3H5.25A2.25 2.25 0 003 5.25v4.318c0 .597.237 1.17.659 1.591l9.581 9.581c.699.699 1.78.872 2.607.33a18.095 18.095 0 005.223-5.223c.542-.827.369-1.908-.33-2.607L11.16 3.66A2.25 2.25 0 009.568 3z" />
              <path strokeLinecap="round" strokeLinejoin="round" d="M6 6h.008v.008H6V6z" />
            </svg>
            <p className="empty-text">No recent activities to display.</p>
            <p className="empty-subtext">Activity will appear here as actions are performed.</p>
          </div>
        </div>
      );
    }

    return (
      <div className="activities-section">
        <div className="section-header">
          <h3 className="section-title">Recent Activities</h3>
          <Button variant="secondary" size="small">View All</Button>
        </div>
        <div className="activities-list">
          {recentActivities.slice(0, 6).map((activity, index) => (
            <div key={index} className={`activity-item activity-item--${getActivityColor(activity.type)}`}>
              <div className="activity-item__icon">
                <svg className="activity-icon" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                  {getActivityIcon(activity.type)}
                </svg>
              </div>
              <div className="activity-item__content">
                <div className="activity-item__header">
                  <div className="activity-item__title">{activity.title}</div>
                  <div className="activity-item__time">{formatTimestamp(activity.timestamp)}</div>
                </div>
                <div className="activity-item__description">{activity.description}</div>
                <div className="activity-item__type">
                  <span className={`activity-badge activity-badge--${getActivityColor(activity.type)}`}>
                    {activity.type}
                  </span>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    );
  };

  const getActivityIcon = (type) => {
    const icons = {
      'property': <path strokeLinecap="round" strokeLinejoin="round" d="m2.25 12 8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25" />,
      'lease': <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />,
      'maintenance': <path strokeLinecap="round" strokeLinejoin="round" d="M11.42 15.17L17.25 21A2.652 2.652 0 0021 17.25l-5.877-5.877M11.42 15.17l2.496-3.03c.317-.384.74-.626 1.208-.766M11.42 15.17l-4.655 5.653a2.548 2.548 0 11-3.586-3.586l6.837-5.63m5.108-.233c.55-.164 1.163-.188 1.743-.14a4.5 4.5 0 004.486-6.336l-3.276 3.277a3.004 3.004 0 01-2.25-2.25l3.276-3.276a4.5 4.5 0 00-6.336 4.486c.091 1.076-.071 2.264-.904 2.95l-.102.085m-1.745 1.437L5.909 7.5H4.5L2.25 3.75l1.5-1.5L7.5 4.5v1.409l4.26 4.26m-1.745 1.437l1.745-1.437m6.615 8.206L15.75 15.75M4.867 19.125h.008v.008h-.008v-.008z" />,
      'user': <path strokeLinecap="round" strokeLinejoin="round" d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z" />,
      'invoice': <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />,
      'expense': <path strokeLinecap="round" strokeLinejoin="round" d="M2.25 8.25h19.5M2.25 9h19.5m-16.5 5.25h6m-6 2.25h3m-3.75 3h15a2.25 2.25 0 002.25-2.25V6.75A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25v10.5A2.25 2.25 0 004.5 19.5z" />
    };
    return icons[type] || icons['property'];
  };

  const renderQuickActions = () => {
    return (
      <div className="quick-actions-section">
        <div className="section-header">
          <h3 className="section-title">Quick Actions</h3>
          <div className="section-subtitle">
            <RoleSwitch
              admin="Full system management"
              landlord="Property and tenant management"
              tenant="Personal account management"
              maintenanceTeam="Maintenance request management"
            />
          </div>
        </div>
        
        <div className="quick-actions-grid">
          {/* Admin Actions - Full system access */}
          <AdminOnly>
            <div className="action-card" onClick={() => navigate(ROUTES.ADMIN_PROPERTIES)}>
              <div className="action-icon action-icon--primary">
                <svg fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
                </svg>
              </div>
              <div className="action-content">
                <h4>Add Property</h4>
                <p>Create new property</p>
              </div>
            </div>
            
            <div className="action-card" onClick={() => navigate(ROUTES.ADMIN_USERS)}>
              <div className="action-icon action-icon--info">
                <svg fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M18 7.5v3m0 0v3m0-3h3m-3 0h-3m-2.25-4.125a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zM3 19.235v-.11a6.375 6.375 0 0112.75 0v.109A12.318 12.318 0 009.374 21c-2.331 0-4.512-.645-6.374-1.764z" />
                </svg>
              </div>
              <div className="action-content">
                <h4>User Management</h4>
                <p>Manage system users</p>
              </div>
            </div>
            
            <div className="action-card" onClick={() => navigate('/admin/reports')}>
              <div className="action-icon action-icon--warning">
                <svg fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z" />
                </svg>
              </div>
              <div className="action-content">
                <h4>System Reports</h4>
                <p>View system-wide reports</p>
              </div>
            </div>
          </AdminOnly>

          {/* Landlord Actions - Property management, tenant relations, financial tracking */}
          <LandlordOnly>
            <div className="action-card" onClick={() => navigate(ROUTES.LANDLORD_PROPERTIES)}>
              <div className="action-icon action-icon--primary">
                <svg fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
                </svg>
              </div>
              <div className="action-content">
                <h4>Add Property</h4>
                <p>Expand your portfolio</p>
              </div>
            </div>
            
            <div className="action-card" onClick={() => navigate(ROUTES.LANDLORD_LEASES)}>
              <div className="action-icon action-icon--success">
                <svg fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
                </svg>
              </div>
              <div className="action-content">
                <h4>Tenant Relations</h4>
                <p>Manage leases & tenants</p>
              </div>
            </div>
            
            <div className="action-card" onClick={() => navigate(ROUTES.LANDLORD_EXPENSES)}>
              <div className="action-icon action-icon--danger">
                <svg fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M2.25 8.25h19.5M2.25 9h19.5m-16.5 5.25h6m-6 2.25h3m-3.75 3h15a2.25 2.25 0 002.25-2.25V6.75A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25v10.5A2.25 2.25 0 004.5 19.5z" />
                </svg>
              </div>
              <div className="action-content">
                <h4>Financial Tracking</h4>
                <p>Track expenses & income</p>
              </div>
            </div>
          </LandlordOnly>

          {/* Tenant Actions - View personal information, submit maintenance requests */}
          <TenantOnly>
            <div className="action-card" onClick={() => navigate(ROUTES.TENANT_MAINTENANCE)}>
              <div className="action-icon action-icon--warning">
                <svg fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M11.42 15.17L17.25 21A2.652 2.652 0 0021 17.25l-5.877-5.877M11.42 15.17l2.496-3.03c.317-.384.74-.626 1.208-.766M11.42 15.17l-4.655 5.653a2.548 2.548 0 11-3.586-3.586l6.837-5.63m5.108-.233c.55-.164 1.163-.188 1.743-.14a4.5 4.5 0 004.486-6.336l-3.276 3.277a3.004 3.004 0 01-2.25-2.25l3.276-3.276a4.5 4.5 0 00-6.336 4.486c.091 1.076-.071 2.264-.904 2.95l-.102.085m-1.745 1.437L5.909 7.5H4.5L2.25 3.75l1.5-1.5L7.5 4.5v1.409l4.26 4.26m-1.745 1.437l1.745-1.437m6.615 8.206L15.75 15.75M4.867 19.125h.008v.008h-.008v-.008z" />
                </svg>
              </div>
              <div className="action-content">
                <h4>Submit Maintenance</h4>
                <p>Report issues & requests</p>
              </div>
            </div>
            
            <div className="action-card" onClick={() => navigate('/tenant/profile')}>
              <div className="action-icon action-icon--info">
                <svg fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z" />
                </svg>
              </div>
              <div className="action-content">
                <h4>Personal Information</h4>
                <p>View & update profile</p>
              </div>
            </div>
            
            <div className="action-card" onClick={() => navigate('/tenant/invoices')}>
              <div className="action-icon action-icon--primary">
                <svg fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />
                </svg>
              </div>
              <div className="action-content">
                <h4>View Invoices</h4>
                <p>Check payments & bills</p>
              </div>
            </div>
          </TenantOnly>

          {/* Maintenance Team Actions - View and update maintenance requests */}
          <MaintenanceOnly>
            <div className="action-card" onClick={() => navigate(ROUTES.MAINTENANCE_TEAM_MAINTENANCE)}>
              <div className="action-icon action-icon--warning">
                <svg fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M8.25 6.75h12M8.25 12h12m-12 5.25h12M3.75 6.75h.007v.008H3.75V6.75zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zM3.75 12h.007v.008H3.75V12zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zM3.75 17.25h.007v.008H3.75v-.008zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z" />
                </svg>
              </div>
              <div className="action-content">
                <h4>Work Queue</h4>
                <p>View assigned requests</p>
              </div>
            </div>
            
            <div className="action-card" onClick={() => navigate('/maintenanceTeam/reports')}>
              <div className="action-icon action-icon--info">
                <svg fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z" />
                </svg>
              </div>
              <div className="action-content">
                <h4>Reports</h4>
                <p>View performance metrics</p>
              </div>
            </div>
          </MaintenanceOnly>
        </div>
      </div>
    );
  };

  const getQuickActionIcon = (iconName) => {
    const icons = {
      'plus': <path strokeLinecap="round" strokeLinejoin="round" d="M12 4.5v15m7.5-7.5h-15" />,
      'file-plus': <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m3.75 9v6m3-3H9m1.5-12H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />,
      'user-plus': <path strokeLinecap="round" strokeLinejoin="round" d="M19 7.5v3m0 0v3m0-3h3m-3 0h-3m-2.25-4.125a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0zM4 19.235v-.11a6.375 6.375 0 0112.674-1.334c.343.061.672.133.994.216" />,
      'chart-bar': <path strokeLinecap="round" strokeLinejoin="round" d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z" />,
      'receipt': <path strokeLinecap="round" strokeLinejoin="round" d="M9 14.25l6-6m4.5-3.493V21.75l-3.75-1.5-3.75 1.5-3.75-1.5-3.75 1.5V4.757c0-1.108.806-2.057 1.907-2.185a48.507 48.507 0 0111.186 0c1.1.128 1.907 1.077 1.907 2.185zM9.75 9h.008v.008H9.75V9zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0zm4.125 4.5h.008v.008h-.008V13.5zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z" />,
      'chart-line': <path strokeLinecap="round" strokeLinejoin="round" d="M2.25 18L9 11.25l4.306 4.307a11.95 11.95 0 015.814-5.519l2.74-1.22m0 0l-5.94-2.28m5.94 2.28l-2.28 5.94" />,
      'tools': <path strokeLinecap="round" strokeLinejoin="round" d="M11.42 15.17L17.25 21A2.652 2.652 0 0021 17.25l-5.877-5.877M11.42 15.17l2.496-3.03c.317-.384.74-.626 1.208-.766M11.42 15.17l-4.655 5.653a2.548 2.548 0 11-3.586-3.586l6.837-5.63m5.108-.233c.55-.164 1.163-.188 1.743-.14a4.5 4.5 0 004.486-6.336l-3.276 3.277a3.004 3.004 0 01-2.25-2.25l3.276-3.276a4.5 4.5 0 00-6.336 4.486c.091 1.076-.071 2.264-.904 2.95l-.102.085m-1.745 1.437L5.909 7.5H4.5L2.25 3.75l1.5-1.5L7.5 4.5v1.409l4.26 4.26m-1.745 1.437l1.745-1.437m6.615 8.206L15.75 15.75M4.867 19.125h.008v.008h-.008v-.008z" />,
      'file-invoice': <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />,
      'file-contract': <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z" />,
      'headset': <path strokeLinecap="round" strokeLinejoin="round" d="M19.114 5.636a9 9 0 010 12.728M16.463 8.288a5.25 5.25 0 010 7.424M6.75 8.25l4.72-4.72a.75.75 0 011.28.53v15.88a.75.75 0 01-1.28.53l-4.72-4.72H4.51c-.88 0-1.59-.79-1.59-1.75v-4.5c0-.96.71-1.75 1.59-1.75h2.24z" />
    };
    return icons[iconName] || icons['plus'];
  };

  if (loading) {
    return (
      <div className="page-content">
        <div className="loading-container">
          <LoadingSpinner />
        </div>
      </div>
    );
  }

  return (
    <div className="page-content">
      <div className="dashboard">
        <div className="dashboard__header">
          <div className="dashboard__title-section">
            <h1 className="dashboard__title">
              Welcome back, {user?.firstName || user?.username}! 👋
            </h1>
            <p className="dashboard__subtitle">
              Here's what's happening with your properties today.
            </p>
          </div>
          <div className="dashboard__actions">
            <Button 
              variant="secondary" 
              size="small"
              onClick={fetchDashboardData}
              disabled={loading}
            >
              <svg className="refresh-icon" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99" />
              </svg>
              Refresh
            </Button>
          </div>
        </div>

        {error && (
          <Alert type="error">{error}</Alert>
        )}

        <div className="dashboard__content">
          <div className="dashboard__stats">
            {renderStatsCards()}
          </div>

          <div className="dashboard__sections">
            <div className="dashboard__primary">
              {renderRecentActivities()}
            </div>
            
            <div className="dashboard__secondary">
              {renderQuickActions()}
              
              {hasRole(USER_ROLES.ADMIN) && (
                <div className="system-summary">
                  <div className="section-header">
                    <h3 className="section-title">System Summary</h3>
                  </div>
                  <div className="summary-grid">
                    <div className="summary-item">
                      <div className="summary-item__icon summary-item__icon--info">
                        <svg fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                          <path strokeLinecap="round" strokeLinejoin="round" d="M10.5 6a7.5 7.5 0 107.5 7.5h-7.5V6z M13.5 10.5H21A7.5 7.5 0 0013.5 3v7.5z" />
                        </svg>
                      </div>
                      <div className="summary-item__content">
                        <span className="summary-item__label">Occupancy Rate</span>
                        <span className="summary-item__value">{(stats.occupancyRate || 0).toFixed(1)}%</span>
                      </div>
                    </div>
                    
                    <div className="summary-item">
                      <div className="summary-item__icon summary-item__icon--success">
                        <svg fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                          <path strokeLinecap="round" strokeLinejoin="round" d="M12 6v12m-3-2.818l.879.659c1.171.879 3.07.879 4.242 0 1.172-.879 1.172-2.303 0-3.182C13.536 12.219 12.768 12 12 12c-.725 0-1.45-.22-2.003-.659-1.106-.879-1.106-2.303 0-3.182s2.9-.879 4.006 0l.415.33M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                      </div>
                      <div className="summary-item__content">
                        <span className="summary-item__label">Revenue/Expense Ratio</span>
                        <span className="summary-item__value">
                          {stats.totalExpenses ? ((stats.totalRevenue || 0) / stats.totalExpenses).toFixed(2) : '0.00'}
                        </span>
                      </div>
                    </div>
                    
                    <div className="summary-item">
                      <div className="summary-item__icon summary-item__icon--warning">
                        <svg fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
                          <path strokeLinecap="round" strokeLinejoin="round" d="M2.25 18L9 11.25l4.306 4.307a11.95 11.95 0 015.814-5.519l2.74-1.22m0 0l-5.94-2.28m5.94 2.28l-2.28 5.94" />
                        </svg>
                      </div>
                      <div className="summary-item__content">
                        <span className="summary-item__label">Avg Revenue per Property</span>
                        <span className="summary-item__value">
                          ${stats.totalProperties ? ((stats.totalRevenue || 0) / stats.totalProperties).toLocaleString() : '0'}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
