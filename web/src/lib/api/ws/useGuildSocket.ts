import { useEffect } from 'react';
import { useHistory, useLocation } from 'react-router-dom';
import { useQueryClient } from 'react-query';
import { getSocket } from '../getSocket';
import { Guild } from '../models';
import { userStore } from '../../stores/userStore';
import { gKey } from '../../utils/querykeys';

export function useGuildSocket() {

  const history = useHistory();
  const cache = useQueryClient();
  const current = userStore(state => state.current);
  const location = useLocation();

  useEffect((): any => {
    const socket = getSocket();
    socket.emit('joinUser', current?.id);

    socket.on('edit_guild', (editedGuild: Guild) => {
      cache.setQueryData<Guild[]>(gKey, (d) => {
        const index = d!.findIndex(c => c.id === editedGuild.id);
        if (index !== -1) {
          d![index] = { ...d![index], name: editedGuild.name, icon: editedGuild.icon };
        }
        return d!;
      });
    });

    socket.on('delete_guild', (deleteId: string) => {
      cache.setQueryData<Guild[]>(gKey, (d) => {
        const isActive = location.pathname.includes(deleteId);
        if (isActive) {
          history.replace('/channels/me');
        }
        return d!.filter(g => g.id !== deleteId);
      });
    });

    socket.on('new_notification', (id: string) => {
      if (!location.pathname.includes(id)) {
        cache.setQueryData<Guild[]>(gKey, (d) => {
          const index = d!.findIndex(c => c.id === id);
          if (index !== -1) {
            d![index] = { ...d![index], hasNotification: true };
          }
          return d!;
        });
      }
    });

    return () => {
      socket.emit('leaveRoom', current?.id);
      socket.disconnect();
    };
  }, [current, cache, history, location]);
}
