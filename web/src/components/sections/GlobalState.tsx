import React, { useEffect } from 'react';
import { useQueryClient } from 'react-query';
import { userStore } from '../../lib/stores/userStore';
import { getSocket } from '../../lib/api/getSocket';
import { homeStore } from '../../lib/stores/homeStore';
import { DMChannel, DMNotification } from '../../lib/api/models';
import { nKey } from '../../lib/utils/querykeys';

type WSMessage = { action: 'new_dm_notification'; data: DMChannel } | { action: 'send_request' };

export const GlobalState: React.FC = ({ children }) => {
  const current = userStore((state) => state.current);
  const inc = homeStore((state) => state.increment);
  const cache = useQueryClient();

  // eslint-disable-next-line consistent-return
  useEffect(() => {
    if (current) {
      const disconnect = (): void => {
        socket.send(JSON.stringify({ action: 'toggleOffline' }));
        socket.close();
      };

      const socket = getSocket();
      socket.send(JSON.stringify({ action: 'toggleOnline' }));
      socket.send(
        JSON.stringify({
          action: 'joinUser',
          room: current?.id,
        })
      );

      socket.addEventListener('message', (event) => {
        const response: WSMessage = JSON.parse(event.data);
        switch (response.action) {
          case 'new_dm_notification': {
            const channel = response.data;
            if (channel.user.id !== current.id) {
              cache.setQueryData<DMNotification[]>(nKey, (data) => {
                const index = data?.findIndex((c) => c.id === channel.id);
                if (index !== -1 && index !== undefined) {
                  return [
                    {
                      ...channel,
                      count: data![index].count + 1,
                    },
                    ...data!.filter((c) => c.id !== channel.id),
                  ];
                }
                return [
                  {
                    ...channel,
                    count: 1,
                  },
                  ...(data || []),
                ];
              });
            }
            break;
          }

          case 'send_request': {
            if (!window.location.pathname.includes('/channels/me')) {
              inc();
            }
            break;
          }

          default:
            break;
        }
      });

      window.addEventListener('beforeunload', disconnect);

      return () => disconnect();
    }
  }, [current, inc, cache]);

  return <>{children}</>;
};
