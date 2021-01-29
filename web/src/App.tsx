import * as React from 'react';
import { QueryClient, QueryClientProvider } from 'react-query';
import { SWRConfig } from 'swr';
import { Routes } from './routes/Routes';
import { request } from './lib/api/setupAxios';

export const App = () => (
  <QueryClientProvider client={new QueryClient()}>
    <SWRConfig value={{ fetcher: request, revalidateOnFocus: false }}>
      <Routes />
    </SWRConfig>
  </QueryClientProvider>
);
