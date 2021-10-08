import { renderHook } from '@testing-library/react-hooks';
import { QueryClientProvider, useQuery } from 'react-query';
import { rest } from 'msw';
import * as React from 'react';
import { createQueryClientWrapper, createTestQueryClientWithData } from '../testUtils';
import { gKey } from '../../lib/utils/querykeys';
import { getUserGuilds } from '../../lib/api/handler/guilds';
import { server } from '../../setupTests';
import { useGetCurrentGuild } from '../../lib/utils/hooks/useGetCurrentGuild';
import { mockGuild, mockGuildList } from '../fixture/guildFixtures';
import { Guild } from '../../lib/models/guild';

describe('useQuery - getUserGuilds', () => {
  it("successfully fetches the user's guild list", async () => {
    const { result, waitForNextUpdate } = renderHook(
      () =>
        useQuery(gKey, async () => {
          const { data } = await getUserGuilds();
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

    const guild = result.current.data?.[0];
    expect(guild).toEqual(mockGuild);
    expect(guild?.id).toBeDefined();
    expect(guild?.name).toBeDefined();
    expect(guild?.ownerId).toBeDefined();
    expect(guild?.hasNotification).toBeFalsy();
    expect(guild?.default_channel_id).toBeDefined();
  });

  it('returns an error when the server returns status 500', async () => {
    server.use(rest.get('*', (req, res, ctx) => res(ctx.status(500))));

    const { result, waitFor } = renderHook(
      () =>
        useQuery(gKey, async () => {
          const { data } = await getUserGuilds();
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

describe('useGetCurrentGuild', () => {
  it('successfully fetches the guild for the given ID', async () => {
    const guildId = mockGuild.id;

    const wrapper: React.FC = ({ children }) => (
      <QueryClientProvider client={createTestQueryClientWithData(gKey, mockGuildList)}>{children}</QueryClientProvider>
    );

    const { result, unmount } = renderHook(() => useGetCurrentGuild(guildId), {
      wrapper,
    });

    try {
      expect(result.current?.id).toEqual(guildId);
    } finally {
      unmount();
    }
  });

  it('returns undefined if it cannot find a guild with the given id', async () => {
    const guildId = mockGuild.id;
    const guild: Guild = {
      id: '12345676890345345',
      name: 'Guild Name',
      default_channel_id: '149587609049385',
      ownerId: '123941059157915',
    };

    const wrapper: React.FC = ({ children }) => (
      <QueryClientProvider client={createTestQueryClientWithData(gKey, [guild])}>{children}</QueryClientProvider>
    );

    const { result, unmount } = renderHook(() => useGetCurrentGuild(guildId), {
      wrapper,
    });

    try {
      expect(result.current).toBeUndefined();
    } finally {
      unmount();
    }
  });

  it('returns undefined if there is not initial data', async () => {
    const guildId = mockGuild.id;

    const wrapper: React.FC = ({ children }) => (
      <QueryClientProvider client={createTestQueryClientWithData(gKey, [])}>{children}</QueryClientProvider>
    );

    const { result, unmount } = renderHook(() => useGetCurrentGuild(guildId), {
      wrapper,
    });

    try {
      expect(result.current).toBeUndefined();
    } finally {
      unmount();
    }
  });
});
