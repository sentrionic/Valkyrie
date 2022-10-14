import { renderHook } from '@testing-library/react-hooks';
import { QueryClientProvider, useQuery } from '@tanstack/react-query';
import { rest } from 'msw';
import * as React from 'react';
import { createQueryClientWrapper, createTestQueryClientWithData, IQueryWrapperProps } from '../testUtils';
import { fKey } from '../../lib/utils/querykeys';
import { server } from '../../setupTests';
import { getFriends } from '../../lib/api/handler/account';
import { mockFriend, mockFriendList } from '../fixture/friendFixture';
import { useGetFriend } from '../../lib/utils/hooks/useGetFriend';

describe('useQuery - getFriends', () => {
  it("successfully fetches the user's friend list", async () => {
    const { result, waitForNextUpdate } = renderHook(
      () =>
        useQuery([fKey], async () => {
          const { data } = await getFriends();
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

    const friend = result.current.data?.[0];
    expect(friend).toEqual(mockFriend);
    expect(friend?.id).toBeDefined();
    expect(friend?.username).toBeDefined();
    expect(friend?.image).toBeDefined();
    expect(friend?.isOnline).toBeDefined();
  });

  it('returns an error when the server returns status 500', async () => {
    server.use(rest.get('*', (req, res, ctx) => res(ctx.status(500))));

    const { result, waitFor } = renderHook(
      () =>
        useQuery([fKey], async () => {
          const { data } = await getFriends();
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

describe('useGetFriend', () => {
  it('successfully fetches the friend for the given ID', async () => {
    const friendId = mockFriend.id;

    const wrapper: React.FC<IQueryWrapperProps> = ({ children }) => (
      <QueryClientProvider client={createTestQueryClientWithData([fKey], mockFriendList)}>
        {children}
      </QueryClientProvider>
    );

    const { result, unmount } = renderHook(() => useGetFriend(friendId), {
      wrapper,
    });

    try {
      expect(result.current?.id).toEqual(friendId);
      expect(result.current).toEqual(mockFriend);
    } finally {
      unmount();
    }
  });

  it('returns undefined if it cannot find a friend with the given id', async () => {
    const friendId = mockFriend.id;
    const friend = {
      id: '12345676890345345',
      username: 'Test User',
      image: 'https://gravatar.com/avatar/c160f8cc69a4f0bf2b0362752353d060?d=identicon',
      isOnline: false,
    };

    const wrapper: React.FC<IQueryWrapperProps> = ({ children }) => (
      <QueryClientProvider client={createTestQueryClientWithData([fKey], [friend])}>{children}</QueryClientProvider>
    );

    const { result, unmount } = renderHook(() => useGetFriend(friendId), {
      wrapper,
    });

    try {
      expect(result.current).toBeUndefined();
    } finally {
      unmount();
    }
  });

  it('returns undefined if there is not initial data', async () => {
    const friendId = mockFriend.id;

    const wrapper: React.FC<IQueryWrapperProps> = ({ children }) => (
      <QueryClientProvider client={createTestQueryClientWithData([fKey], [])}>{children}</QueryClientProvider>
    );

    const { result, unmount } = renderHook(() => useGetFriend(friendId), {
      wrapper,
    });

    try {
      expect(result.current).toBeUndefined();
    } finally {
      unmount();
    }
  });
});
