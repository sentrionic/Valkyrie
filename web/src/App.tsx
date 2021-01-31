import * as React from 'react';
import { QueryClient, QueryClientProvider } from 'react-query';
import { Routes } from './routes/Routes';
import { GlobalState } from './components/sections/GlobalState';

export const App = () => (
  <QueryClientProvider client={new QueryClient()}>
    <GlobalState>
      <Routes />
    </ GlobalState>
  </QueryClientProvider>
);
