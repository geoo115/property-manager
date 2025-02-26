import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";

import Navbar from "./components/Navbar";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Dashboard from "./pages/Dashboard";
import Properties from "./pages/Properties";
import Leases from "./pages/Leases";
import Users from "./pages/Users";
import Maintenance from "./pages/Maintenance";
import Invoices from "./pages/Invoices";
import Expenses from "./pages/Expenses";
import ProtectedRoute from "./components/ProtectedRoute";

function App() {
  return (
    <Router>
      <Navbar />
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />

        {/* Admin Routes - Full Access */}
        <Route element={<ProtectedRoute allowedRoles={["admin"]} />}>
          <Route path="/admin/dashboard" element={<Dashboard />} />
          <Route path="/admin/properties" element={<Properties />} />
          <Route path="/admin/leases" element={<Leases />} />
          <Route path="/admin/users" element={<Users />} />
          <Route path="/admin/maintenances" element={<Maintenance />} />
          {/* Admin Accounting routes */}
          <Route path="/admin/accounting/invoices" element={<Invoices />} />
          <Route path="/admin/accounting/expenses" element={<Expenses />} />
        </Route>

        {/* Landlord Routes - View Only */}
        <Route element={<ProtectedRoute allowedRoles={["landlord"]} />}>
          <Route path="/landlord/properties" element={<Properties />} />
          <Route path="/landlord/leases" element={<Leases />} />
          {/* Landlord Accounting: you may choose a common route */}
          <Route path="/accounting/invoices" element={<Invoices />} />
          <Route path="/accounting/expenses" element={<Expenses />} />
        </Route>

        {/* Tenant & Maintenance Team Routes */}
        <Route element={<ProtectedRoute allowedRoles={["tenant"]} />}>
          <Route path="/tenant/leases" element={<Leases />} />
          <Route path="/tenant/maintenance/:leaseId" element={<Maintenance />} />
        </Route>
        <Route element={<ProtectedRoute allowedRoles={["maintenanceTeam"]} />}>
          <Route path="/maintenanceTeam/maintenances" element={<Maintenance />} />
        </Route>
      </Routes>
      <ToastContainer />
    </Router>
  );
}

export default App;