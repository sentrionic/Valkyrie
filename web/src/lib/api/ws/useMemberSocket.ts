import { useEffect } from 'react';
import { getSocket } from '../getSocket';
import { Member } from '../models';
import { useQueryClient } from 'react-query';

export function useMemberSocket(guildId: string, key: string) {
  const cache = useQueryClient();

  useEffect((): any => {
    const socket = getSocket();
    socket.emit('joinGuild', guildId);
    socket.on('add_member', (newMember: Member) => {
      cache.setQueryData<Member[]>(key, (data) => {
        return [...data!, newMember].sort((a, b) => {
          if (a.nickname && b.nickname) {
            return a.nickname.localeCompare(b.nickname);
          } else if (a.nickname && !b.nickname) {
            return a.nickname.localeCompare(b.username);
          } else if (!a.nickname && b.nickname) {
            return a.username.localeCompare(b.nickname);
          } else {
            return a.username.localeCompare(b.username);
          }
        });
      });
    });

    socket.on('remove_member', (memberId: string) => {
      cache.setQueryData<Member[]>(key, (data) => {
        return [...data!.filter((m) => m.id !== memberId)];
      });
    });

    socket.on('toggle_online', (memberId: string) => {
      cache.setQueryData<Member[]>(key, (data) => {
        const index = data!.findIndex((m) => m.id === memberId);
        if (index !== -1) data![index].isOnline = true;
        return data!;
      });
    });

    socket.on('toggle_offline', (memberId: string) => {
      cache.setQueryData<Member[]>(key, (data) => {
        const index = data!.findIndex((m) => m.id === memberId);
        if (index !== -1) data![index].isOnline = false;
        return data!;
      });
    });

    return () => {
      socket.emit('leaveRoom', guildId);
      socket.disconnect();
    };
  }, [key, cache, guildId]);
}
