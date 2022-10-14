import { useEffect } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { getSocket } from '../getSocket';
import { Member } from '../../models/member';
import { mKey } from '../../utils/querykeys';

type WSMessage =
  | { action: 'remove_member' | 'toggle_online' | 'toggle_offline'; data: string }
  | { action: 'add_member'; data: Member };

export function useMemberSocket(guildId: string): void {
  const cache = useQueryClient();

  useEffect((): any => {
    const socket = getSocket();
    const key = [mKey, guildId];

    socket.send(
      JSON.stringify({
        action: 'joinGuild',
        room: guildId,
      })
    );

    socket.addEventListener('message', (event) => {
      const response: WSMessage = JSON.parse(event.data);
      switch (response.action) {
        case 'add_member': {
          cache.setQueryData<Member[]>(key, (data) =>
            // Add member and sort array by nickname, then username
            [...(data ?? []), response.data].sort((a, b) => {
              if (a.nickname && b.nickname) {
                return a.nickname.localeCompare(b.nickname);
              }
              if (a.nickname && !b.nickname) {
                return a.nickname.localeCompare(b.username);
              }
              if (!a.nickname && b.nickname) {
                return a.username.localeCompare(b.nickname);
              }
              return a.username.localeCompare(b.username);
            })
          );
          break;
        }

        case 'remove_member': {
          cache.setQueryData<Member[]>(key, (data) => [...(data?.filter((m) => m.id !== response.data) ?? [])]);
          break;
        }

        case 'toggle_online': {
          const memberId = response.data;
          cache.setQueryData<Member[]>(key, (d) => {
            if (!d) return [];
            return d.map((m) => (m.id === memberId ? { ...m, isOnline: true } : m));
          });
          break;
        }

        case 'toggle_offline': {
          const memberId = response.data;
          cache.setQueryData<Member[]>(key, (d) => {
            if (!d) return [];
            return d.map((m) => (m.id === memberId ? { ...m, isOnline: false } : m));
          });
          break;
        }

        default:
          break;
      }
    });

    return () => {
      socket.send(
        JSON.stringify({
          action: 'leaveRoom',
          room: guildId,
        })
      );
      socket.close();
    };
  }, [cache, guildId]);
}
