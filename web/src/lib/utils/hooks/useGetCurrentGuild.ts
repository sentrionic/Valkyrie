import { useQuery } from 'react-query';
import { Guild } from '../../api/models';
import { gKey } from '../querykeys';

export function useGetCurrentGuild(guildId: string): Guild | undefined {
  const { data: guildData } = useQuery<Guild[]>(gKey);
  return guildData?.find(g => g.id === guildId);
}

