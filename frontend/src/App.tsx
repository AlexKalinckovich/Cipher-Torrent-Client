import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { Dashboard } from './features/torrents/Dashboard';
import {Authentication} from "./features/auth/Authentication.tsx";
import {Profile} from "./features/profile/Profile.tsx";

export const App: React.FC = () => {
  return (
      <BrowserRouter>
        <Routes>
          <Route path="/" element={renderAuthentication()} />
          <Route path="/dashboard" element={renderDashboard()} />
          <Route path="/profile" element={renderUserProfile()}/>
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

const renderUserProfile = () : React.ReactElement => {
    return <Profile />;
}