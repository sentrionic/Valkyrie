import { useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useQueryClient } from '@tanstack/react-query';
import { getSocket } from '../getSocket';
import { userStore } from '../../stores/userStore';
import { gKey } from '../../utils/querykeys';
import { Guild } from '../../models/guild';

type WSMessage =
  | { action: 'delete_guild' | 'remove_from_guild' | 'new_notification'; data: string }
  | { action: 'edit_guild'; data: Guild };

export function useGuildSocket(): void {
  const navigate = useNavigate();
  const cache = useQueryClient();
  const current = userStore((state) => state.current);
  const location = useLocation();

  useEffect((): any => {
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
        case 'edit_guild': {
          const editedGuild = response.data;
          cache.setQueryData<Guild[]>([gKey], (d) => {
            if (!d) return [];
            return d.map((g) =>
              g.id === editedGuild.id
                ? {
                    ...g,
                    name: editedGuild.name,
                    icon: editedGuild.icon,
                  }
                : g
            );
          });
          break;
        }

        case 'delete_guild': {
          const deleteId = response.data;
          cache.setQueryData<Guild[]>([gKey], (d) => {
            const isActive = location.pathname.includes(deleteId);
            if (isActive) {
              navigate('/channels/me', { replace: true });
            }
            return d?.filter((g) => g.id !== deleteId) ?? [];
          });
          break;
        }

        case 'new_notification': {
          const id = response.data;
          if (!location.pathname.includes(id)) {
            cache.setQueryData<Guild[]>([gKey], (d) => {
              if (!d) return [];
              return d.map((g) =>
                g.id === id
                  ? {
                      ...g,
                      hasNotification: true,
                    }
                  : g
              );
            });
          }
          break;
        }

        case 'remove_from_guild': {
          cache.setQueryData<Guild[]>([gKey], (d) => {
            const guildId = response.data;
            const isActive = location.pathname.includes(guildId);
            if (isActive) {
              navigate('/channels/me', { replace: true });
            }
            return d?.filter((g) => g.id !== guildId) ?? [];
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
  }, [current, cache, navigate, location]);
}
