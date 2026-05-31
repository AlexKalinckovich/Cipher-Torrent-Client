import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { Dashboard } from './features/torrents/Dashboard';
import {Authentication} from "./features/auth/Authentication.tsx";
import {Profile} from "./features/profile/Profile.tsx";
import {PacketInspector} from "@/features/packetInspector/PacketInspector.tsx";

export const App: React.FC = () => {
  return (
      <BrowserRouter>
        <Routes>
          <Route path="/" element={renderAuthentication()} />
          <Route path="/dashboard" element={renderDashboard()} />
          <Route path="/profile" element={renderUserProfile()}/>
          <Route path="/inspector" element={renderPackageInspector()} />
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

const renderPackageInspector = () : React.ReactElement => {
    return <PacketInspector />;
}