import { renderHook } from '@testing-library/react-hooks';
import { useQuery } from 'react-query';
import { rest } from 'msw';
import { aKey } from '../../lib/utils/querykeys';
import { createQueryClientWrapper } from '../testUtils';
import { server } from '../../setupTests';
import { getAccount } from '../../lib/api/handler/account';
import { mockAccount } from '../fixture/accountFixture';

describe('useQuery - getAccount', () => {
  it("successfully fetches the user's info", async () => {
    const { result, waitForNextUpdate } = renderHook(
      () =>
        useQuery(aKey, async () => {
          const { data } = await getAccount();
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

    const account = result.current.data;
    expect(account).toEqual(mockAccount);
    expect(account?.id).toBeDefined();
    expect(account?.username).toBeDefined();
    expect(account?.email).toBeDefined();
    expect(account?.image).toBeDefined();
  });

  it('returns an error when the server returns status 500', async () => {
    server.use(rest.get('*', (req, res, ctx) => res(ctx.status(500))));

    const { result, waitFor } = renderHook(
      () =>
        useQuery(aKey, async () => {
          const { data } = await getAccount();
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
