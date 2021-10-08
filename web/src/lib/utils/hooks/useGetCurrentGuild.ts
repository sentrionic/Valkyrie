import { useQuery } from 'react-query';
import { gKey } from '../querykeys';
import { Guild } from '../../models/guild';

export function useGetCurrentGuild(guildId: string): Guild | undefined {
  const { data: guildData } = useQuery<Guild[]>(gKey);
  return guildData?.find((g) => g.id === guildId);
}
