import React, { useState, useEffect } from 'react';
import useAuth from '../../hooks/useAuth';
import Button from '../common/Button';
import './Header.css';

const Header = () => {
  const { user, logout } = useAuth();
  const [currentDateTime, setCurrentDateTime] = useState(new Date());

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentDateTime(new Date());
    }, 1000);

    return () => clearInterval(timer);
  }, []);

  const handleLogout = () => {
    logout();
  };

  const formatDateTime = (date) => {
    const options = { 
      weekday: 'long', 
      year: 'numeric', 
      month: 'long', 
      day: 'numeric' 
    };
    const dateStr = date.toLocaleDateString('en-US', options);
    const timeStr = date.toLocaleTimeString('en-US', { 
      hour: '2-digit', 
      minute: '2-digit' 
    });
    return `${dateStr}\n${timeStr}`;
  };

  if (!user) return null;

  return (
    <header className="header">
      <div className="header__datetime">
        <span className="header__datetime-text">
          {formatDateTime(currentDateTime)}
        </span>
      </div>

      <div className="header__actions">
        <Button
          variant="danger"
          size="small"
          onClick={handleLogout}
          icon={
            <svg fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15M12 9l-3 3m0 0l3 3m-3-3h12.75" />
            </svg>
          }
        >
          Logout
        </Button>
      </div>
    </header>
  );
};

export default Header;
