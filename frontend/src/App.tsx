import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { Dashboard } from './features/torrents/Dashboard';

export const App: React.FC = () => {
  return (
      <BrowserRouter>
        <Routes>
          <Route path="/" element={renderRedirect()} />
          <Route path="/dashboard" element={renderDashboard()} />
        </Routes>
      </BrowserRouter>
  );
};

const renderRedirect = (): React.ReactElement => {
  return <Navigate to="/dashboard" replace />;
};

const renderDashboard = (): React.ReactElement => {
  return <Dashboard />;
};