import { useEffect } from 'react';
import { useHistory, useLocation } from 'react-router-dom';
import { useQueryClient } from 'react-query';
import { getSocket } from '../getSocket';
import { Guild } from '../models';
import { userStore } from '../../stores/userStore';
import { gKey } from '../../utils/querykeys';

type WSMessage =
  | { action: 'delete_guild' | 'remove_from_guild' | 'new_notification'; data: string }
  | { action: 'edit_guild'; data: Guild };

export function useGuildSocket() {
  const history = useHistory();
  const cache = useQueryClient();
  const current = userStore((state) => state.current);
  const location = useLocation();

  useEffect((): any => {
    const socket = getSocket();

    socket.send(JSON.stringify({ action: 'joinUser', room: current?.id }));

    socket.addEventListener('message', (event) => {
      const response: WSMessage = JSON.parse(event.data);
      switch (response.action) {
        case 'edit_guild': {
          const editedGuild = response.data;
          cache.setQueryData<Guild[]>(gKey, (d) => {
            const index = d!.findIndex((c) => c.id === editedGuild.id);
            if (index !== -1) {
              d![index] = { ...d![index], name: editedGuild.name, icon: editedGuild.icon };
            }
            return d!;
          });
          break;
        }

        case 'delete_guild': {
          const deleteId = response.data;
          cache.setQueryData<Guild[]>(gKey, (d) => {
            const isActive = location.pathname.includes(deleteId);
            if (isActive) {
              history.replace('/channels/me');
            }
            return d!.filter((g) => g.id !== deleteId);
          });
          break;
        }

        case 'new_notification': {
          const id = response.data;
          if (!location.pathname.includes(id)) {
            cache.setQueryData<Guild[]>(gKey, (d) => {
              const index = d!.findIndex((c) => c.id === id);
              if (index !== -1) {
                d![index] = { ...d![index], hasNotification: true };
              }
              return d!;
            });
          }
          break;
        }

        case 'remove_from_guild': {
          cache.setQueryData<Guild[]>(gKey, (d) => {
            const guildId = response.data;
            const isActive = location.pathname.includes(guildId);
            if (isActive) {
              history.replace('/channels/me');
            }
            return d!.filter((g) => g.id !== guildId);
          });
          break;
        }
      }
    });

    return () => {
      socket.send(JSON.stringify({ action: 'leaveRoom', room: current?.id }));
      socket.close();
    };
  }, [current, cache, history, location]);
}
