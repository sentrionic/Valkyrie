/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { Navigate } from 'react-router-dom';
import { userStore } from '../lib/stores/userStore';

interface IProps {
  children: React.ReactNode;
}

export const AuthRoute: React.FC<IProps> = ({ children }) => {
  const storage = JSON.parse(localStorage.getItem('user-storage')!!);
  const current = userStore((state) => state.current);

  if (current || storage?.state?.current) {
    return <>{children}</>;
  }

  return <Navigate to="/login" />;
};
