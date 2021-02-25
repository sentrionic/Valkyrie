import { useEffect } from 'react';
import { getSocket } from '../getSocket';
import { Channel } from '../models';
import { useHistory, useLocation } from 'react-router-dom';
import { useQueryClient } from 'react-query';
import { useGetCurrentGuild } from '../../utils/hooks/useGetCurrentGuild';

export function useChannelSocket(guildId: string, key: string) {

  const location = useLocation();
  const history = useHistory();
  const cache = useQueryClient();
  const guild = useGetCurrentGuild(guildId);

  useEffect((): any => {
    const socket = getSocket();
    socket.emit('joinGuild', guildId);

    const disconnect = () => {
      socket.emit('leaveGuild', guildId);
      socket.disconnect();
    }

    socket.on('add_channel', (newChannel: Channel) => {
      cache.setQueryData<Channel[]>(key, (data) => {
        return [...data!, newChannel];
      });
    });

    socket.on('edit_channel', (editedChannel: Channel) => {
      cache.setQueryData<Channel[]>(key, (d) => {
        const index = d!.findIndex(c => c.id === editedChannel.id);
        if (index !== -1) {
          d![index] = editedChannel;
        } else if (editedChannel.isPublic) {
          d!.push(editedChannel);
        }
        return d!;
      });
    });

    socket.on('delete_channel', (deleteId: string) => {
      cache.setQueryData<Channel[]>(key, (d) => {
        const currentPath = `/channels/${guildId}/${deleteId}`;
        if (location.pathname === currentPath && guild) {
          if (deleteId === guild.default_channel_id) {
            history.replace('/channels/me')
          } else {
            history.replace(`${guild.default_channel_id}`);
          }
        }
        return d!.filter(c => c.id !== deleteId);
      });
    });

    socket.on('new_notification', (id: string) => {
      const currentPath = `/channels/${guildId}/${id}`;
      if (location.pathname !== currentPath) {
        cache.setQueryData<Channel[]>(key, (d) => {
          const index = d!.findIndex(c => c.id === id);
          if (index !== -1) {
            d![index] = { ...d![index], hasNotification: true };
          }
          return d!;
        });
      }
    });

    window.addEventListener('beforeunload', disconnect);

    return () => disconnect();
  }, [guildId, key, cache, history, location, guild]);
}
