import { renderHook } from '@testing-library/react-hooks';
import { useQuery } from '@tanstack/react-query';
import { rest } from 'msw';
import { mKey } from '../../lib/utils/querykeys';
import { createQueryClientWrapper } from '../testUtils';
import { server } from '../../setupTests';
import { getGuildMembers } from '../../lib/api/handler/guilds';
import { mockGuild } from '../fixture/guildFixtures';
import { mockMember } from '../fixture/memberFixtures';

describe('useQuery - getGuildMembers', () => {
  it("successfully fetches the guild's member list", async () => {
    const guildId = mockGuild.id;

    const { result, waitForNextUpdate } = renderHook(
      () =>
        useQuery([mKey, guildId], async () => {
          const { data } = await getGuildMembers(guildId);
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

    const member = result.current.data?.[0];
    expect(member).toEqual(mockMember);
    expect(member?.id).toBeDefined();
    expect(member?.username).toBeDefined();
    expect(member?.image).toBeDefined();
    expect(member?.isOnline).toBeDefined();
    expect(member?.isFriend).toBeDefined();
  });

  it('returns an error when the server returns status 500', async () => {
    const guildId = mockGuild.id;
    server.use(rest.get('*', (req, res, ctx) => res(ctx.status(500))));

    const { result, waitFor } = renderHook(
      () =>
        useQuery([mKey, guildId], async () => {
          const { data } = await getGuildMembers(guildId);
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
