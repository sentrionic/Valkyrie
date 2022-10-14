import { renderHook } from '@testing-library/react-hooks';
import { QueryClientProvider, useQuery } from '@tanstack/react-query';
import { rest } from 'msw';
import * as React from 'react';
import { createQueryClientWrapper, createTestQueryClientWithData, IQueryWrapperProps } from '../testUtils';
import { cKey } from '../../lib/utils/querykeys';
import { server } from '../../setupTests';
import { mockGuild } from '../fixture/guildFixtures';
import { getChannels } from '../../lib/api/handler/channel';
import { mockChannel, mockChannelList } from '../fixture/channelFixtures';
import { useGetCurrentChannel } from '../../lib/utils/hooks/useGetCurrentChannel';
import { Channel } from '../../lib/models/channel';

describe('useQuery - getChannels', () => {
  it("successfully fetches the current guild's channel list", async () => {
    const guildId = '12312456127277383';

    const { result, waitForNextUpdate } = renderHook(
      () =>
        useQuery([cKey, guildId], async () => {
          const { data } = await getChannels(guildId);
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

    const channel = result.current.data?.[0];
    expect(channel).toEqual(mockChannel);
    expect(channel?.id).toBeDefined();
    expect(channel?.name).toBeDefined();
    expect(channel?.isPublic).toBeDefined();
    expect(channel?.hasNotification).toBeFalsy();
  });

  it('returns an error when the server returns status 500', async () => {
    const guildId = '12312456127277383';
    server.use(rest.get('*', (req, res, ctx) => res(ctx.status(500))));

    const { result, waitFor } = renderHook(
      () =>
        useQuery([cKey, guildId], async () => {
          const { data } = await getChannels(guildId);
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

describe('useGetCurrentChannel', () => {
  it('successfully fetches the channel for the given ID and key', async () => {
    const channelId = mockChannel.id;
    const key = [cKey, mockGuild.id];

    const wrapper: React.FC<IQueryWrapperProps> = ({ children }) => (
      <QueryClientProvider client={createTestQueryClientWithData(key, mockChannelList)}>{children}</QueryClientProvider>
    );

    const { result, unmount } = renderHook(() => useGetCurrentChannel(channelId, mockGuild.id), {
      wrapper,
    });

    try {
      expect(result.current?.id).toEqual(channelId);
      expect(result.current).toEqual(mockChannel);
    } finally {
      unmount();
    }
  });

  it('returns undefined if it cannot find a channel with the given id', async () => {
    const channelId = mockChannel.id;
    const key = [cKey, mockGuild.id];
    const channel: Channel = {
      id: '12345676890345345',
      name: 'Guild Name',
      hasNotification: false,
      isPublic: true,
    };

    const wrapper: React.FC<IQueryWrapperProps> = ({ children }) => (
      <QueryClientProvider client={createTestQueryClientWithData(key, [channel])}>{children}</QueryClientProvider>
    );

    const { result, unmount } = renderHook(() => useGetCurrentChannel(channelId, mockGuild.id), {
      wrapper,
    });

    try {
      expect(result.current).toBeUndefined();
    } finally {
      unmount();
    }
  });

  it('returns undefined if there is not initial data', async () => {
    const channelId = mockChannel.id;
    const key = [cKey, mockGuild.id];

    const wrapper: React.FC<IQueryWrapperProps> = ({ children }) => (
      <QueryClientProvider client={createTestQueryClientWithData(key, [])}>{children}</QueryClientProvider>
    );

    const { result, unmount } = renderHook(() => useGetCurrentChannel(channelId, mockGuild.id), {
      wrapper,
    });

    try {
      expect(result.current).toBeUndefined();
    } finally {
      unmount();
    }
  });
});
