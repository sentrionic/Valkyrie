import { useEffect } from 'react';
import { useQueryClient } from 'react-query';
import { getSocket } from '../getSocket';
import { RequestResponse } from '../models';
import { userStore } from '../../stores/userStore';
import { rKey } from '../../utils/querykeys';
import { homeStore } from '../../stores/homeStore';

type WSMessage = { action: 'add_request'; data: RequestResponse };

export function useRequestSocket() {
  const current = userStore((state) => state.current);
  const setRequests = homeStore((state) => state.setRequests);
  const cache = useQueryClient();

  useEffect((): any => {
    const socket = getSocket();

    socket.send(JSON.stringify({ action: 'joinUser', room: current?.id }));

    socket.addEventListener('message', (event) => {
      const response: WSMessage = JSON.parse(event.data);
      switch (response.action) {
        case 'add_request': {
          cache.setQueryData<RequestResponse[]>(rKey, (data) => {
            return [...data!, response.data].sort((a, b) => a.username.localeCompare(b.username));
          });
          break;
        }
      }
    });

    return () => {
      socket.send(JSON.stringify({ action: 'leaveRoom', room: current?.id }));
      socket.close();
    };
  }, [cache, current, setRequests]);
}
