import React, { Suspense } from "react";
import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";

import { AuthProvider } from "./context/AuthContext";
import { USER_ROLES, ROUTES } from "./constants";

// Layout Components
import AppLayout from "./components/layout/AppLayout";
import ProtectedRoute from "./components/ProtectedRoute";
import LoadingSpinner from "./components/common/LoadingSpinner";

// Import global styles
import "./App.css";
import "./styles.css";

// Lazy load pages for better performance
const Login = React.lazy(() => import("./pages/Login"));
const Register = React.lazy(() => import("./pages/Register"));
const Dashboard = React.lazy(() => import("./pages/Dashboard"));
const Properties = React.lazy(() => import("./pages/Properties"));
const Leases = React.lazy(() => import("./pages/Leases"));
const Users = React.lazy(() => import("./pages/Users"));
const Maintenance = React.lazy(() => import("./pages/Maintenance"));
const Invoices = React.lazy(() => import("./pages/Invoices"));
const Expenses = React.lazy(() => import("./pages/Expenses"));
const TenantProfile = React.lazy(() => import("./pages/TenantProfile"));
const RBACTest = React.lazy(() => import("./pages/RBACTest"));

// Enhanced loading fallback component
const PageLoader = () => (
  <div className="page-loader">
    <LoadingSpinner size="large" />
    <span className="page-loader__text">Loading page...</span>
  </div>
);

const App = () => {
  return (
    <AuthProvider>
      <Router>
        <div className="app">
          <Suspense fallback={<PageLoader />}>
            <Routes>
              {/* Public routes - No layout needed */}
              <Route path={ROUTES.LOGIN} element={<Login />} />
              <Route path={ROUTES.REGISTER} element={<Register />} />

              {/* Protected routes with app layout */}
              <Route element={<AppLayout />}>
                {/* Redirect root to dashboard */}
                <Route path={ROUTES.HOME} element={<Navigate to="/dashboard" replace />} />
                
                {/* Dashboard - accessible to all authenticated users */}
                <Route element={<ProtectedRoute allowedRoles={Object.values(USER_ROLES)} />}>
                  <Route path="/dashboard" element={<Dashboard />} />
                </Route>

                {/* Admin Routes - Full Access */}
                <Route element={<ProtectedRoute allowedRoles={[USER_ROLES.ADMIN]} />}>
                  <Route path={ROUTES.ADMIN_PROPERTIES} element={<Properties />} />
                  <Route path={ROUTES.ADMIN_LEASES} element={<Leases />} />
                  <Route path={ROUTES.ADMIN_USERS} element={<Users />} />
                  <Route path={ROUTES.ADMIN_MAINTENANCE} element={<Maintenance />} />
                  <Route path={ROUTES.ADMIN_INVOICES} element={<Invoices />} />
                  <Route path={ROUTES.ADMIN_EXPENSES} element={<Expenses />} />
                  <Route path="/admin/reports" element={<div>System Reports</div>} />
                </Route>

                {/* Landlord Routes - Property management, tenant relations, financial tracking */}
                <Route element={<ProtectedRoute allowedRoles={[USER_ROLES.LANDLORD]} />}>
                  <Route path={ROUTES.LANDLORD_PROPERTIES} element={<Properties />} />
                  <Route path={ROUTES.LANDLORD_LEASES} element={<Leases />} />
                  <Route path={ROUTES.LANDLORD_INVOICES} element={<Invoices />} />
                  <Route path={ROUTES.LANDLORD_EXPENSES} element={<Expenses />} />
                  <Route path="/landlord/tenants" element={<div>Tenant Relations</div>} />
                </Route>

                {/* Tenant Routes - View personal information, submit maintenance requests */}
                <Route element={<ProtectedRoute allowedRoles={[USER_ROLES.TENANT]} />}>
                  <Route path={ROUTES.TENANT_LEASES} element={<Leases />} />
                  <Route path={ROUTES.TENANT_MAINTENANCE} element={<Maintenance />} />
                  <Route path="/tenant/profile" element={<TenantProfile />} />
                  <Route path="/tenant/invoices" element={<Invoices />} />
                  <Route path="/tenant/payments" element={<div>Payment History</div>} />
                </Route>
                
                {/* Maintenance Team Routes - View and update maintenance requests */}
                <Route element={<ProtectedRoute allowedRoles={[USER_ROLES.MAINTENANCE_TEAM]} />}>
                  <Route path={ROUTES.MAINTENANCE_TEAM_MAINTENANCE} element={<Maintenance />} />
                  <Route path="/maintenanceTeam/reports" element={<div>Maintenance Reports</div>} />
                </Route>
                
                {/* RBAC Test Route - For development/testing purposes */}
                <Route element={<ProtectedRoute allowedRoles={Object.values(USER_ROLES)} />}>
                  <Route path="/rbac-test" element={<RBACTest />} />
                </Route>
                
                {/* Catch all route */}
                <Route path="*" element={<Navigate to="/dashboard" replace />} />
              </Route>
            </Routes>
          </Suspense>
          
          <ToastContainer 
            position="top-right"
            autoClose={5000}
            hideProgressBar={false}
            newestOnTop
            closeOnClick
            rtl={false}
            pauseOnFocusLoss
            draggable
            pauseOnHover
            theme="light"
            toastClassName="modern-toast"
            bodyClassName="modern-toast-body"
            limit={3}
            style={{ zIndex: 99999 }}
          />
        </div>
      </Router>
    </AuthProvider>
  );
};

export default App;