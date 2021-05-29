import React, { useEffect } from 'react';
import { userStore } from '../../lib/stores/userStore';
import { getSocket } from '../../lib/api/getSocket';
import { homeStore } from '../../lib/stores/homeStore';
import { DMChannel, DMNotification } from '../../lib/api/models';
import { useQueryClient } from 'react-query';
import { nKey } from '../../lib/utils/querykeys';

export const GlobalState: React.FC = ({ children }) => {
  const current = userStore((state) => state.current);
  const inc = homeStore((state) => state.increment);
  const cache = useQueryClient();

  useEffect(() => {
    if (current) {
      const disconnect = () => {
        socket.emit('toggleOffline');
        socket.disconnect();
      };

      const socket = getSocket();
      socket.emit('toggleOnline');
      socket.emit('joinUser', current?.id);

      socket.on('new_dm_notification', (channel: DMChannel) => {
        if (channel.user.id !== current.id) {
          cache.setQueryData<DMNotification[]>(nKey, (data) => {
            const index = data?.findIndex((c) => c.id === channel.id);
            if (index !== -1 && index !== undefined) {
              return [{ ...channel, count: data![index].count + 1 }, ...data!.filter((c) => c.id !== channel.id)];
            } else {
              return [{ ...channel, count: 1 }, ...(data || [])];
            }
          });
        }
      });

      socket.on('send_request', () => {
        if (!window.location.pathname.includes('/channels/me')) {
          inc();
        }
      });

      window.addEventListener('beforeunload', disconnect);

      return () => disconnect();
    }
  }, [current, inc, cache]);

  return <>{children}</>;
};
