import * as React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { AppRoutes } from './routes/Routes';
import { GlobalState } from './components/sections/GlobalState';

const client = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      staleTime: Infinity,
      cacheTime: 0,
    },
  },
});

export const App: React.FC = () => (
  <QueryClientProvider client={client}>
    <GlobalState>
      <AppRoutes />
    </GlobalState>
  </QueryClientProvider>
);
