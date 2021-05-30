import { useEffect } from 'react';
import { getSocket } from '../getSocket';
import { Channel } from '../models';
import { useHistory, useLocation } from 'react-router-dom';
import { useQueryClient } from 'react-query';
import { useGetCurrentGuild } from '../../utils/hooks/useGetCurrentGuild';

type WSMessage =
  | { action: 'delete_channel' | 'new_notification'; data: string }
  | { action: 'add_channel' | 'edit_channel'; data: Channel };

export function useChannelSocket(guildId: string, key: string) {
  const location = useLocation();
  const history = useHistory();
  const cache = useQueryClient();
  const guild = useGetCurrentGuild(guildId);
  useEffect((): any => {
    const socket = getSocket();

    socket.send(JSON.stringify({ action: 'joinGuild', room: guildId }));

    const disconnect = () => {
      socket.send(JSON.stringify({ action: 'leaveGuild', room: guildId }));
      socket.close();
    };

    socket.addEventListener('message', (event) => {
      const response: WSMessage = JSON.parse(event.data);
      switch (response.action) {
        case 'add_channel': {
          cache.setQueryData<Channel[]>(key, (data) => {
            return [...data!, response.data];
          });
          break;
        }

        case 'edit_channel': {
          const editedChannel = response.data;
          cache.setQueryData<Channel[]>(key, (d) => {
            const index = d!.findIndex((c) => c.id === editedChannel.id);
            if (index !== -1) {
              d![index] = editedChannel;
            } else if (editedChannel.isPublic) {
              d!.push(editedChannel);
            }
            return d!;
          });
          break;
        }

        case 'delete_channel': {
          const deleteId = response.data;
          cache.setQueryData<Channel[]>(key, (d) => {
            const currentPath = `/channels/${guildId}/${deleteId}`;
            if (location.pathname === currentPath && guild) {
              if (deleteId === guild.default_channel_id) {
                history.replace('/channels/me');
              } else {
                history.replace(`${guild.default_channel_id}`);
              }
            }
            return d!.filter((c) => c.id !== deleteId);
          });
          break;
        }

        case 'new_notification': {
          const id = response.data;
          const currentPath = `/channels/${guildId}/${id}`;
          if (location.pathname !== currentPath) {
            cache.setQueryData<Channel[]>(key, (d) => {
              const index = d!.findIndex((c) => c.id === id);
              if (index !== -1) {
                d![index] = { ...d![index], hasNotification: true };
              }
              return d!;
            });
          }
          break;
        }
      }
    });

    window.addEventListener('beforeunload', disconnect);

    return () => disconnect();
  }, [guildId, key, cache, history, location, guild]);
}
