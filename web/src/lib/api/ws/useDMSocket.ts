import { useEffect } from 'react';
import { getSocket } from '../getSocket';
import { DMChannel } from '../models';
import { useQueryClient } from 'react-query';
import { userStore } from '../../stores/userStore';
import { dmKey } from '../../utils/querykeys';

export function useDMSocket() {

  const current = userStore(state => state.current);
  const cache = useQueryClient();

  useEffect(() => {
    const socket = getSocket();
    socket.emit('joinUser', current?.id);

    socket.on('push_to_top', (dmId: string) => {
      cache.setQueryData<DMChannel[]>(dmKey, (data) => {
        const index = data!.findIndex(d => d.id === dmId);
        if (index === 0 || index === -1) return [...data!];
        const dm = data![index];
        return [dm!, ...data!.filter(d => d.id !== dmId)];
      });
    });

    return () => {
      socket.emit('leaveRoom', current?.id);
      socket.disconnect();
    };
  }, [current, cache]);
}
