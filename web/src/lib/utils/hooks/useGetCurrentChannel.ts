import { useQuery } from 'react-query';
import { Channel } from '../../models/channel';

export function useGetCurrentChannel(channelId: string, key: string): Channel | undefined {
  const { data } = useQuery<Channel[]>(key);
  return data?.find((c) => c.id === channelId);
}
