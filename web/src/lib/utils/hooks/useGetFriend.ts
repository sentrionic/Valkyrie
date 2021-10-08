import { useQuery } from 'react-query';
import { fKey } from '../querykeys';
import { Friend } from '../../models/friend';

export function useGetFriend(id: string): Friend | undefined {
  const { data } = useQuery<Friend[]>(fKey);
  return data?.find((f) => f.id === id);
}
