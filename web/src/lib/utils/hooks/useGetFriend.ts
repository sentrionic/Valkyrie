import { useQuery } from 'react-query';
import { Member } from '../../api/models';
import { fKey } from '../querykeys';

export function useGetFriend(id: string): Member | undefined {
  const { data } = useQuery<Member[]>(fKey);
  return data?.find(f => f.id === id);
}
