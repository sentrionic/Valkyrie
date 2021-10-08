import { renderHook } from '@testing-library/react-hooks';
import { QueryClientProvider, useQuery } from 'react-query';
import { rest } from 'msw';
import * as React from 'react';
import { dmKey } from '../../lib/utils/querykeys';
import { createQueryClientWrapper, createTestQueryClientWithData } from '../testUtils';
import { server } from '../../setupTests';
import { getUserDMs } from '../../lib/api/handler/dm';
import { mockDMChannel, mockDMChannelList } from '../fixture/dmFixtures';
import { useGetCurrentDM } from '../../lib/utils/hooks/useGetCurrentDM';
import { DMChannel } from '../../lib/models/dm';

describe('useQuery - getUserDMs', () => {
  it("successfully fetches the user's dm list", async () => {
    const { result, waitForNextUpdate } = renderHook(
      () =>
        useQuery(dmKey, async () => {
          const { data } = await getUserDMs();
          return data;
        }),
      {
        wrapper: createQueryClientWrapper,
      }
    );

    expect(result.current.isFetching).toBe(true);
    await waitForNextUpdate();

    expect(result.current.isFetching).toBe(false);
    expect(result.current.isSuccess).toBe(true);

    expect(result.current.data).not.toBeNull();
    expect(result.current.data).not.toBeUndefined();
    expect(result.current.data?.length).toEqual(1);

    const dm = result.current.data?.[0];
    expect(dm).toEqual(mockDMChannel);
    expect(dm?.id).toBeDefined();
    expect(dm?.user).toBeDefined();

    const user = dm?.user;
    expect(user?.id).toBeDefined();
    expect(user?.image).toBeDefined();
    expect(user?.username).toBeDefined();
    expect(user?.isOnline).toBeDefined();
  });

  it('returns an error when the server returns status 500', async () => {
    server.use(rest.get('*', (req, res, ctx) => res(ctx.status(500))));

    const { result, waitFor } = renderHook(
      () =>
        useQuery(dmKey, async () => {
          const { data } = await getUserDMs();
          return data;
        }),
      {
        wrapper: createQueryClientWrapper,
      }
    );

    await waitFor(() => result.current.isError);

    expect(result.current.error).toBeDefined();
    expect(result.current.data).toBeUndefined();
  });
});

describe('useGetCurrentDM', () => {
  it('successfully fetches the dm for the given ID', async () => {
    const channelId = mockDMChannel.id;

    const wrapper: React.FC = ({ children }) => (
      <QueryClientProvider client={createTestQueryClientWithData(dmKey, mockDMChannelList)}>
        {children}
      </QueryClientProvider>
    );

    const { result, unmount } = renderHook(() => useGetCurrentDM(channelId), {
      wrapper,
    });

    try {
      expect(result.current?.id).toEqual(channelId);
      expect(result.current).toEqual(mockDMChannel);
    } finally {
      unmount();
    }
  });

  it('returns undefined if it cannot find a dm with the given id', async () => {
    const channelId = mockDMChannel.id;
    const dmChannel: DMChannel = {
      id: '12345676890345345',
      user: {
        image: '',
        isOnline: true,
        isFriend: false,
        username: 'Test User',
        id: '123941059157915',
      },
    };

    const wrapper: React.FC = ({ children }) => (
      <QueryClientProvider client={createTestQueryClientWithData(dmKey, [dmChannel])}>{children}</QueryClientProvider>
    );

    const { result, unmount } = renderHook(() => useGetCurrentDM(channelId), {
      wrapper,
    });

    try {
      expect(result.current).toBeUndefined();
    } finally {
      unmount();
    }
  });

  it('returns undefined if there is not initial data', async () => {
    const channelId = mockDMChannel.id;

    const wrapper: React.FC = ({ children }) => (
      <QueryClientProvider client={createTestQueryClientWithData(dmKey, [])}>{children}</QueryClientProvider>
    );

    const { result, unmount } = renderHook(() => useGetCurrentDM(channelId), {
      wrapper,
    });

    try {
      expect(result.current).toBeUndefined();
    } finally {
      unmount();
    }
  });
});
