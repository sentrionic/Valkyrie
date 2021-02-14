import { useEffect } from 'react';
import { useQueryClient } from 'react-query';
import { getSocket } from '../getSocket';
import { Member } from '../models';
import { userStore } from '../../stores/userStore';
import { fKey } from '../../utils/querykeys';

export function useFriendSocket() {

  const current = userStore(state => state.current);
  const cache = useQueryClient();

  useEffect((): any => {
    const socket = getSocket();
    socket.emit('joinUser', current?.id);
    socket.on('add_friend', (newFriend: Member) => {
      cache.setQueryData<Member[]>(fKey, (data) => {
        return [...data!, newFriend].sort((a, b) => a.username.localeCompare(b.username));
      });
    });

    socket.on('remove_friend', (memberId: string) => {
      cache.setQueryData<Member[]>(fKey, (data) => {
        return [...data!.filter(m => m.id !== memberId)];
      });
    });

    socket.on('toggle_online', (memberId: string) => {
      cache.setQueryData<Member[]>(fKey, (data) => {
        const index = data!.findIndex(m => m.id === memberId);
        data![index].isOnline = true;
        return data!;
      });
    });

    socket.on('toggle_offline', (memberId: string) => {
      cache.setQueryData<Member[]>(fKey, (data) => {
        const index = data!.findIndex(m => m.id === memberId);
        data![index].isOnline = false;
        return data!;
      });
    });

    return () => {
      socket.emit('leaveRoom', current?.id);
      socket.disconnect();
    };
  }, [cache, current]);
}
