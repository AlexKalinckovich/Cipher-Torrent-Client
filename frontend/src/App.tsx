import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { Dashboard } from './features/torrents/Dashboard';
import {Authentication} from "./features/auth/Authentication.tsx";

export const App: React.FC = () => {
  return (
      <BrowserRouter>
        <Routes>
          <Route path="/" element={renderAuthentication()} />
          <Route path="/dashboard" element={renderDashboard()} />
        </Routes>
      </BrowserRouter>
  );
};

const renderAuthentication = (): React.ReactElement => {
  return <Authentication/>;
};

const renderDashboard = (): React.ReactElement => {
  return <Dashboard />;
};