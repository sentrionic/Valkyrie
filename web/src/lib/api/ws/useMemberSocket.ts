import { useEffect } from 'react';
import { getSocket } from '../getSocket';
import { Member } from '../models';
import { useQueryClient } from 'react-query';

type WSMessage =
  | { action: 'remove_member' | 'toggle_online' | 'toggle_offline'; data: string }
  | { action: 'add_member'; data: Member };

export function useMemberSocket(guildId: string, key: string) {
  const cache = useQueryClient();
  useEffect((): any => {
    const socket = getSocket();

    socket.send(JSON.stringify({ action: 'joinGuild', room: guildId }));

    socket.addEventListener('message', (event) => {
      const response: WSMessage = JSON.parse(event.data);
      switch (response.action) {
        case 'add_member': {
          cache.setQueryData<Member[]>(key, (data) => {
            return [...data!, response.data].sort((a, b) => {
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
          break;
        }

        case 'remove_member': {
          cache.setQueryData<Member[]>(key, (data) => {
            return [...data!.filter((m) => m.id !== response.data)];
          });
          break;
        }

        case 'toggle_online': {
          const memberId = response.data;
          cache.setQueryData<Member[]>(key, (data) => {
            const index = data!.findIndex((m) => m.id === memberId);
            if (index !== -1) data![index].isOnline = true;
            return data!;
          });
          break;
        }

        case 'toggle_offline': {
          const memberId = response.data;
          cache.setQueryData<Member[]>(key, (data) => {
            const index = data!.findIndex((m) => m.id === memberId);
            if (index !== -1) data![index].isOnline = false;
            return data!;
          });
          break;
        }
      }
    });

    return () => {
      socket.send(JSON.stringify({ action: 'leaveRoom', room: guildId }));
      socket.close();
    };
  }, [key, cache, guildId]);
}
