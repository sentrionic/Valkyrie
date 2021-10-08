import { renderHook } from '@testing-library/react-hooks';
import { useInfiniteQuery } from 'react-query';
import { rest } from 'msw';
import { createQueryClientWrapper } from '../testUtils';
import { server } from '../../setupTests';
import { getMessages } from '../../lib/api/handler/messages';
import { mockChannel } from '../fixture/channelFixtures';
import { mockMessageList } from '../fixture/messageFixtures';
import { Message } from '../../lib/models/message';

describe('useQuery - getMessages', () => {
  it("successfully fetches the channel's message list", async () => {
    const channelId = mockChannel.id;
    const qKey = `messages-${channelId}`;

    const { result, waitForNextUpdate } = renderHook(
      () =>
        useInfiniteQuery<Message[]>(
          qKey,
          async ({ pageParam = null }) => {
            const { data: messageData } = await getMessages(channelId, pageParam);
            return messageData;
          },
          {
            staleTime: 0,
            cacheTime: 0,
            getNextPageParam: (lastPage) => (lastPage.length ? lastPage[lastPage.length - 1].createdAt : ''),
          }
        ),
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

    expect(result.current.data).toEqual({ pageParams: [undefined], pages: [mockMessageList] });
  });

  it('returns an error when the server returns status 500', async () => {
    const channelId = mockChannel.id;
    const qKey = `messages-${channelId}`;
    server.use(rest.get('*', (req, res, ctx) => res(ctx.status(500))));

    const { result, waitFor } = renderHook(
      () =>
        useInfiniteQuery<Message[]>(
          qKey,
          async ({ pageParam = null }) => {
            const { data: messageData } = await getMessages(channelId, pageParam);
            return messageData;
          },
          {
            staleTime: 0,
            cacheTime: 0,
            getNextPageParam: (lastPage) => (lastPage.length ? lastPage[lastPage.length - 1].createdAt : ''),
          }
        ),
      {
        wrapper: createQueryClientWrapper,
      }
    );

    await waitFor(() => result.current.isError);

    expect(result.current.error).toBeDefined();
    expect(result.current.data).toBeUndefined();
  });
});
