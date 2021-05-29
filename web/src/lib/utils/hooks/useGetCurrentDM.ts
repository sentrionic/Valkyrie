import { useQuery } from 'react-query';
import { DMChannel } from '../../api/models';
import { dmKey } from '../querykeys';

export function useGetCurrentDM(channelId: string): DMChannel | undefined {
  const { data } = useQuery<DMChannel[]>(dmKey);
  return data?.find((c) => c.id === channelId);
}
