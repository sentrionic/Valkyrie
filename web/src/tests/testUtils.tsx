import { rest } from 'msw';
import * as React from 'react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { mockGuildList } from './fixture/guildFixtures';
import { mockDMChannelList } from './fixture/dmFixtures';
import { mockAccount } from './fixture/accountFixture';
import { mockFriendList } from './fixture/friendFixture';
import { mockRequestList } from './fixture/requestFixtures';
import { mockChannelList } from './fixture/channelFixtures';
import { mockMemberList } from './fixture/memberFixtures';
import { mockMessageList } from './fixture/messageFixtures';

const createTestQueryClient = (): QueryClient =>
  new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
    logger: {
      // eslint-disable-next-line no-console
      log: console.log,
      // eslint-disable-next-line no-console
      warn: console.warn,
      error: () => {},
    },
  });

export const createTestQueryClientWithData = (key: string[], data: any[]): QueryClient => {
  const client = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
    logger: {
      // eslint-disable-next-line no-console
      log: console.log,
      // eslint-disable-next-line no-console
      warn: console.warn,
      error: () => {},
    },
  });
  client.setQueryData(key, () => [...data]);

  return client;
};

export interface IQueryWrapperProps {
  children: React.ReactNode;
}

export const createQueryClientWrapper: React.FC<IQueryWrapperProps> = ({ children }) => (
  <QueryClientProvider client={createTestQueryClient()}>{children}</QueryClientProvider>
);

export const handlers = [
  rest.get('*/guilds', (req, res, ctx) => res(ctx.status(200), ctx.json(mockGuildList))),
  rest.get('*/channels/me/dm', (req, res, ctx) => res(ctx.status(200), ctx.json(mockDMChannelList))),
  rest.get('*/account', (req, res, ctx) => res(ctx.status(200), ctx.json(mockAccount))),
  rest.get('*/account/me/friends', (req, res, ctx) => res(ctx.status(200), ctx.json(mockFriendList))),
  rest.get('*/account/me/pending', (req, res, ctx) => res(ctx.status(200), ctx.json(mockRequestList))),
  rest.get('*/channels/*', (req, res, ctx) => res(ctx.status(200), ctx.json(mockChannelList))),
  rest.get('*/guilds/*/members', (req, res, ctx) => res(ctx.status(200), ctx.json(mockMemberList))),
  rest.get('*/messages/*', (req, res, ctx) => res(ctx.status(200), ctx.json(mockMessageList))),
];
