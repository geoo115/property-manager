import React from 'react';
import { Outlet } from 'react-router-dom';
import Sidebar from './Sidebar';
import Header from './Header';
import useAuth from '../../hooks/useAuth';
import './AppLayout.css';

const AppLayout = () => {
  const { user } = useAuth();

  // Don't render layout if user is not authenticated
  if (!user) {
    return <Outlet />;
  }

  return (
    <div className="app-layout">
      <Sidebar />
      <div className="app-layout__main">
        <Header />
        <main className="app-layout__content">
          <div className="page-content">
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  );
};

export default AppLayout;
