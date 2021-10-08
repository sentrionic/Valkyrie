import { useQuery } from 'react-query';
import { dmKey } from '../querykeys';
import { DMChannel } from '../../models/dm';

export function useGetCurrentDM(channelId: string): DMChannel | undefined {
  const { data } = useQuery<DMChannel[]>(dmKey);
  return data?.find((c) => c.id === channelId);
}
