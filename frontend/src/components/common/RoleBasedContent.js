import React from 'react';
import PropTypes from 'prop-types';
import useAuth from '../../hooks/useAuth';
import { USER_ROLES } from '../../constants';

/**
 * RoleBasedContent component - Renders content based on user role
 */
const RoleBasedContent = ({ children, ...props }) => {
  const { user, hasRole, hasAnyRole, hasUserPermission } = useAuth();

  if (!user) {
    return null;
  }

  // If children is a function, call it with auth utilities
  if (typeof children === 'function') {
    return children({ 
      user, 
      hasRole, 
      hasAnyRole, 
      hasUserPermission,
      isAdmin: hasRole(USER_ROLES.ADMIN),
      isLandlord: hasRole(USER_ROLES.LANDLORD),
      isTenant: hasRole(USER_ROLES.TENANT),
      isMaintenanceTeam: hasRole(USER_ROLES.MAINTENANCE_TEAM),
    });
  }

  return children;
};

/**
 * AdminOnly component - Only renders for admin users
 */
export const AdminOnly = ({ children, fallback = null }) => {
  const { hasRole } = useAuth();
  
  if (!hasRole(USER_ROLES.ADMIN)) {
    return fallback;
  }
  
  return children;
};

/**
 * LandlordOnly component - Only renders for landlord users
 */
export const LandlordOnly = ({ children, fallback = null }) => {
  const { hasRole } = useAuth();
  
  if (!hasRole(USER_ROLES.LANDLORD)) {
    return fallback;
  }
  
  return children;
};

/**
 * TenantOnly component - Only renders for tenant users
 */
export const TenantOnly = ({ children, fallback = null }) => {
  const { hasRole } = useAuth();
  
  if (!hasRole(USER_ROLES.TENANT)) {
    return fallback;
  }
  
  return children;
};

/**
 * MaintenanceOnly component - Only renders for maintenance team users
 */
export const MaintenanceOnly = ({ children, fallback = null }) => {
  const { hasRole } = useAuth();
  
  if (!hasRole(USER_ROLES.MAINTENANCE_TEAM)) {
    return fallback;
  }
  
  return children;
};

/**
 * RequireRole component - Only renders if user has specific role(s)
 */
export const RequireRole = ({ roles, children, fallback = null }) => {
  const { hasAnyRole } = useAuth();
  
  const roleArray = Array.isArray(roles) ? roles : [roles];
  
  if (!hasAnyRole(roleArray)) {
    return fallback;
  }
  
  return children;
};

/**
 * RequirePermission component - Only renders if user has specific permission
 */
export const RequirePermission = ({ resource, action, children, fallback = null }) => {
  const { hasUserPermission } = useAuth();
  
  if (!hasUserPermission(resource, action)) {
    return fallback;
  }
  
  return children;
};

/**
 * RoleSwitch component - Renders different content based on user role
 */
export const RoleSwitch = ({ 
  admin, 
  landlord, 
  tenant, 
  maintenanceTeam, 
  fallback = null 
}) => {
  const { user } = useAuth();
  
  if (!user) return fallback;
  
  switch (user.role) {
    case USER_ROLES.ADMIN:
      return admin || fallback;
    case USER_ROLES.LANDLORD:
      return landlord || fallback;
    case USER_ROLES.TENANT:
      return tenant || fallback;
    case USER_ROLES.MAINTENANCE_TEAM:
      return maintenanceTeam || fallback;
    default:
      return fallback;
  }
};

RoleBasedContent.propTypes = {
  children: PropTypes.oneOfType([PropTypes.node, PropTypes.func]).isRequired,
};

AdminOnly.propTypes = {
  children: PropTypes.node.isRequired,
  fallback: PropTypes.node,
};

LandlordOnly.propTypes = {
  children: PropTypes.node.isRequired,
  fallback: PropTypes.node,
};

TenantOnly.propTypes = {
  children: PropTypes.node.isRequired,
  fallback: PropTypes.node,
};

MaintenanceOnly.propTypes = {
  children: PropTypes.node.isRequired,
  fallback: PropTypes.node,
};

RequireRole.propTypes = {
  roles: PropTypes.oneOfType([PropTypes.string, PropTypes.arrayOf(PropTypes.string)]).isRequired,
  children: PropTypes.node.isRequired,
  fallback: PropTypes.node,
};

RequirePermission.propTypes = {
  resource: PropTypes.string.isRequired,
  action: PropTypes.string.isRequired,
  children: PropTypes.node.isRequired,
  fallback: PropTypes.node,
};

RoleSwitch.propTypes = {
  admin: PropTypes.node,
  landlord: PropTypes.node,
  tenant: PropTypes.node,
  maintenanceTeam: PropTypes.node,
  fallback: PropTypes.node,
};

export default RoleBasedContent;
