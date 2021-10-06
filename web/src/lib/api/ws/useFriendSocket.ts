import { useEffect } from 'react';
import { useQueryClient } from 'react-query';
import { getSocket } from '../getSocket';
import { Member } from '../models';
import { userStore } from '../../stores/userStore';
import { fKey } from '../../utils/querykeys';
import { homeStore } from '../../stores/homeStore';

type WSMessage =
  | { action: 'toggle_online' | 'toggle_offline' | 'remove_friend'; data: string }
  | { action: 'requestCount'; data: number }
  | { action: 'add_friend'; data: Member };

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
          cache.setQueryData<Member[]>(fKey, (d) => {
            const data = d ?? [];
            const index = data.findIndex((m) => m.id === response.data);
            if (index !== -1) data[index].isOnline = true;
            return data;
          });
          break;
        }

        case 'toggle_offline': {
          cache.setQueryData<Member[]>(fKey, (d) => {
            const data = d ?? [];
            const index = data.findIndex((m) => m.id === response.data);
            if (index !== -1) data[index].isOnline = false;
            return data;
          });
          break;
        }

        case 'requestCount': {
          setRequests(response.data);
          break;
        }

        case 'add_friend': {
          cache.setQueryData<Member[]>(fKey, (data) =>
            [...data!, response.data].sort((a, b) => a.username.localeCompare(b.username))
          );
          break;
        }

        case 'remove_friend': {
          cache.setQueryData<Member[]>(fKey, (data) => [...data!.filter((m) => m.id !== response.data)]);
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
