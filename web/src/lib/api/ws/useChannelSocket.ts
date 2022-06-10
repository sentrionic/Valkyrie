import { useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useQueryClient } from 'react-query';
import { getSocket } from '../getSocket';
import { useGetCurrentGuild } from '../../utils/hooks/useGetCurrentGuild';
import { userStore } from '../../stores/userStore';
import { Channel } from '../../models/channel';
import { VCMember, VoiceResponse } from '../../models/voice';
import { vcKey } from '../../utils/querykeys';

type WSMessage =
  | { action: 'delete_channel' | 'new_notification'; data: string }
  | { action: 'add_channel' | 'add_private_channel' | 'edit_channel'; data: Channel }
  | { action: 'joinVoice' | 'leaveVoice'; data: VoiceResponse };

export function useChannelSocket(guildId: string, key: string): void {
  const location = useLocation();
  const navigate = useNavigate();
  const cache = useQueryClient();
  const guild = useGetCurrentGuild(guildId);
  const current = userStore((state) => state.current);

  useEffect((): any => {
    const socket = getSocket();

    socket.send(
      JSON.stringify({
        action: 'joinGuild',
        room: guildId,
      })
    );
    socket.send(
      JSON.stringify({
        action: 'joinUser',
        room: current?.id,
      })
    );

    const disconnect = (): void => {
      socket.send(
        JSON.stringify({
          action: 'leaveGuild',
          room: guildId,
        })
      );
      socket.send(
        JSON.stringify({
          action: 'leaveRoom',
          room: current?.id,
        })
      );
      socket.close();
    };

    socket.addEventListener('message', (event) => {
      const response: WSMessage = JSON.parse(event.data);
      switch (response.action) {
        case 'add_channel': {
          cache.setQueryData<Channel[]>(key, (data) => [...data!, response.data]);
          break;
        }

        case 'add_private_channel': {
          cache.setQueryData<Channel[]>(key, (data) => [...data!, response.data]);
          break;
        }

        case 'edit_channel': {
          const editedChannel = response.data;
          cache.setQueryData<Channel[]>(key, (d) => {
            const data = d ?? [];
            const index = data.findIndex((c) => c.id === editedChannel.id);
            if (index !== -1) {
              data[index] = editedChannel;
            } else if (editedChannel.isPublic) {
              data.push(editedChannel);
            }
            return data;
          });
          break;
        }

        case 'delete_channel': {
          const deleteId = response.data;
          cache.setQueryData<Channel[]>(key, (d) => {
            const currentPath = `/channels/${guildId}/${deleteId}`;
            if (location.pathname === currentPath && guild) {
              if (deleteId === guild.default_channel_id) {
                navigate('/channels/me', { replace: true });
              } else {
                navigate(`/channels/${guild.id}/${guild.default_channel_id}`, { replace: true });
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
              const data = d ?? [];
              const index = data.findIndex((c) => c.id === id);
              if (index !== -1) {
                data[index] = {
                  ...d![index],
                  hasNotification: true,
                };
              }
              return data;
            });
          }
          break;
        }

        case 'joinVoice': {
          const { data } = response;
          // Remove the current user from the list
          cache.setQueryData<VCMember[]>(vcKey(guildId), (_) => data.clients.filter((e) => e.id !== current?.id));
          break;
        }

        case 'leaveVoice': {
          const { data } = response;
          // Remove the current user from the list
          cache.setQueryData<VCMember[]>(vcKey(guildId), (_) => data.clients.filter((e) => e.id !== current?.id));
          break;
        }

        default:
          break;
      }
    });

    window.addEventListener('beforeunload', disconnect);

    return () => disconnect();
  }, [guildId, key, cache, navigate, location, guild, current]);
}
