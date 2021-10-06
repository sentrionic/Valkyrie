import { useEffect } from 'react';
import { useQueryClient } from 'react-query';
import { getSocket } from '../getSocket';
import { DMChannel } from '../models';
import { userStore } from '../../stores/userStore';
import { dmKey } from '../../utils/querykeys';

type WSMessage = { action: 'push_to_top'; data: string };

export function useDMSocket(): void {
  const current = userStore((state) => state.current);
  const cache = useQueryClient();

  useEffect(() => {
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
        case 'push_to_top': {
          const dmId = response.data;
          cache.setQueryData<DMChannel[]>(dmKey, (d) => {
            const data = d ?? [];
            const index = data.findIndex((dm) => dm.id === dmId);
            if (index === 0 || index === -1) return [...data];
            const dm = data[index];
            return [dm, ...data.filter((dc) => dc.id !== dmId)];
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
          room: current?.id,
        })
      );
      socket.close();
    };
  }, [current, cache]);
}
