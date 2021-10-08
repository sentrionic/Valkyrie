import { useEffect } from 'react';
import { useQueryClient } from 'react-query';
import { getSocket } from '../getSocket';
import { userStore } from '../../stores/userStore';
import { rKey } from '../../utils/querykeys';
import { homeStore } from '../../stores/homeStore';
import { FriendRequest } from '../../models/friend';

type WSMessage = { action: 'add_request'; data: FriendRequest };

export function useRequestSocket(): void {
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

    socket.addEventListener('message', (event) => {
      const response: WSMessage = JSON.parse(event.data);
      switch (response.action) {
        case 'add_request': {
          cache.setQueryData<FriendRequest[]>(rKey, (data) =>
            [...data!, response.data].sort((a, b) => a.username.localeCompare(b.username))
          );
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
