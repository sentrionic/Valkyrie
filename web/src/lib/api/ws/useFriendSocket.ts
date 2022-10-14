import { useEffect } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { getSocket } from '../getSocket';
import { userStore } from '../../stores/userStore';
import { fKey } from '../../utils/querykeys';
import { homeStore } from '../../stores/homeStore';
import { Friend } from '../../models/friend';

type WSMessage =
  | { action: 'toggle_online' | 'toggle_offline' | 'remove_friend'; data: string }
  | { action: 'requestCount'; data: number }
  | { action: 'add_friend'; data: Friend };

export function useFriendSocket(): void {
  const current = userStore((state) => state.current);
  const setRequests = homeStore((state) => state.setRequests);
  const cache = useQueryClient();

  useEffect((): any => {
    const socket = getSocket();
    socket.send(
      JSON.stringify({
        action: 'joinUser',
        room: current?.id,
      })
    );

    socket.send(JSON.stringify({ action: 'getRequestCount' }));

    socket.addEventListener('message', (event) => {
      const response: WSMessage = JSON.parse(event.data);
      switch (response.action) {
        case 'toggle_online': {
          cache.setQueryData<Friend[]>([fKey], (d) => {
            if (!d) return [];
            return d.map((f) => (f.id === response.data ? { ...f, isOnline: true } : f));
          });
          break;
        }

        case 'toggle_offline': {
          cache.setQueryData<Friend[]>([fKey], (d) => {
            if (!d) return [];
            return d.map((f) => (f.id === response.data ? { ...f, isOnline: false } : f));
          });
          break;
        }

        case 'requestCount': {
          setRequests(response.data);
          break;
        }

        case 'add_friend': {
          cache.setQueryData<Friend[]>([fKey], (data) =>
            [...(data ?? []), response.data].sort((a, b) => a.username.localeCompare(b.username))
          );
          break;
        }

        case 'remove_friend': {
          cache.setQueryData<Friend[]>([fKey], (data) => [...(data?.filter((m) => m.id !== response.data) ?? [])]);
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
          room: current?.id,
        })
      );

      socket.close();
    };
  }, [cache, current, setRequests]);
}
