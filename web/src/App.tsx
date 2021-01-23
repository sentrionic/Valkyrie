import * as React from 'react';
import { QueryClient, QueryClientProvider } from 'react-query';
import { Routes } from './routes/Routes';

export const App = () => (
  <QueryClientProvider client={new QueryClient()}>
    <Routes />
  </QueryClientProvider>
);
