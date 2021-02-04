import * as React from 'react';
import { QueryClient, QueryClientProvider } from 'react-query';
import { Routes } from './routes/Routes';
import { GlobalState } from './components/sections/GlobalState';

const client = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      staleTime: Infinity,
      cacheTime: 0
    }
  }
});

export const App = () => (
  <QueryClientProvider client={client}>
    <GlobalState>
      <Routes />
    </ GlobalState>
  </QueryClientProvider>
);
