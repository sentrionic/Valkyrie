import { useEffect } from 'react';
import { useQueryClient } from 'react-query';
import { getSocket } from '../getSocket';
import { RequestResponse } from '../models';
import { userStore } from '../../stores/userStore';
import { rKey } from '../../utils/querykeys';
import { homeStore } from '../../stores/homeStore';

export function useRequestSocket() {
  const current = userStore((state) => state.current);
  const setRequests = homeStore((state) => state.setRequests);
  const cache = useQueryClient();

  useEffect((): any => {
    const socket = getSocket();
    socket.emit('joinUser', current?.id);
    socket.on('add_request', (newRequest: RequestResponse) => {
      cache.setQueryData<RequestResponse[]>(rKey, (data) => {
        return [...data!, newRequest].sort((a, b) => a.username.localeCompare(b.username));
      });
    });

    return () => {
      socket.emit('leaveRoom', current?.id);
      socket.disconnect();
    };
  }, [cache, current, setRequests]);
}
