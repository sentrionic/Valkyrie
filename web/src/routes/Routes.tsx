import React from 'react';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { Login } from './Login';
import { Register } from './Register';
import { ForgotPassword } from './ForgotPassword';
import { ResetPassword } from './ResetPassword';
import { Home } from './Home';
import { ViewGuild } from './ViewGuild';
import { AuthRoute } from './AuthRoute';
import { Settings } from './Settings';
import { Landing } from './Landing';
import { Invite } from './Invite';

export const AppRoutes: React.FC = () => (
  <BrowserRouter>
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />
      <Route path="/forgot-password" element={<ForgotPassword />} />
      <Route path="/reset-password/:token" element={<ResetPassword />} />
      <Route path="/" element={<Landing />} />
      <Route
        path="/channels/me"
        element={
          <AuthRoute>
            <Home />
          </AuthRoute>
        }
      />
      <Route
        path="/channels/me/:channelId"
        element={
          <AuthRoute>
            <Home />
          </AuthRoute>
        }
      />
      <Route
        path="/channels/:guildId/:channelId"
        element={
          <AuthRoute>
            <ViewGuild />
          </AuthRoute>
        }
      />
      <Route
        path="/account"
        element={
          <AuthRoute>
            <Settings />
          </AuthRoute>
        }
      />
      <Route
        path="/:link"
        element={
          <AuthRoute>
            <Invite />
          </AuthRoute>
        }
      />
    </Routes>
  </BrowserRouter>
);
