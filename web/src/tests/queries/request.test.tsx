import { renderHook } from '@testing-library/react-hooks';
import { useQuery } from 'react-query';
import { rest } from 'msw';
import { rKey } from '../../lib/utils/querykeys';
import { createQueryClientWrapper } from '../testUtils';
import { server } from '../../setupTests';
import { getPendingRequests } from '../../lib/api/handler/account';
import { mockRequest } from '../fixture/requestFixtures';

describe('useQuery - getPendingRequests', () => {
  it("successfully fetches the user's pending requests", async () => {
    const { result, waitForNextUpdate } = renderHook(
      () =>
        useQuery(rKey, async () => {
          const { data } = await getPendingRequests();
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

    const request = result.current.data?.[0];
    expect(request).toEqual(mockRequest);
    expect(request?.id).toBeDefined();
    expect(request?.username).toBeDefined();
    expect(request?.image).toBeDefined();
    expect(request?.type).toBeDefined();
  });

  it('returns an error when the server returns status 500', async () => {
    server.use(rest.get('*', (req, res, ctx) => res(ctx.status(500))));

    const { result, waitFor } = renderHook(
      () =>
        useQuery(rKey, async () => {
          const { data } = await getPendingRequests();
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
