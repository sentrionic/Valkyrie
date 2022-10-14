import { useEffect } from 'react';
import { InfiniteData, useQueryClient } from '@tanstack/react-query';
import { getSocket } from '../getSocket';
import { userStore } from '../../stores/userStore';
import { channelStore } from '../../stores/channelStore';
import { Message } from '../../models/message';
import { msgKey } from '../../utils/querykeys';

type WSMessage =
  | { action: 'new_message' | 'edit_message'; data: Message }
  | { action: 'addToTyping' | 'removeFromTyping' | 'delete_message'; data: string };

export function useMessageSocket(channelId: string): void {
  const current = userStore((state) => state.current);
  const store = channelStore();
  const cache = useQueryClient();

  useEffect((): any => {
    store.reset();
    const socket = getSocket();

    socket.send(
      JSON.stringify({
        action: 'joinChannel',
        room: channelId,
      })
    );

    socket.addEventListener('message', (event) => {
      const response: WSMessage = JSON.parse(event.data);
      switch (response.action) {
        case 'new_message': {
          cache.setQueryData<InfiniteData<Message[]>>([msgKey, channelId], (d) => {
            if (!d) return { pages: [], pageParams: [] };
            return {
              pages: d.pages.map((messages, i) => (i === 0 ? [response.data, ...messages] : messages)),
              pageParams: [...d.pageParams],
            };
          });
          break;
        }

        case 'edit_message': {
          const editMessage = response.data;
          cache.setQueryData<InfiniteData<Message[]>>([msgKey, channelId], (d) => {
            if (!d) return { pages: [], pageParams: [] };
            return {
              pages: d.pages.map((messages) =>
                messages.map((m) =>
                  m.id === editMessage.id ? { ...m, text: editMessage.text, updatedAt: editMessage.updatedAt } : m
                )
              ),
              pageParams: [...d.pageParams],
            };
          });
          break;
        }

        case 'delete_message': {
          const messageId = response.data;
          cache.setQueryData<InfiniteData<Message[]>>([msgKey, channelId], (d) => {
            if (!d) return { pages: [], pageParams: [] };
            return {
              pages: d.pages.map((messages) => messages.filter((m) => m.id !== messageId)),
              pageParams: [...d.pageParams],
            };
          });
          break;
        }

        case 'addToTyping': {
          const username = response.data;
          if (username !== current?.username) store.addTyping(username);
          break;
        }

        case 'removeFromTyping': {
          const username = response.data;
          if (username !== current?.username) store.removeTyping(username);
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
          room: channelId,
        })
      );

      socket.close();
    };
    // eslint-disable-next-line
  }, [channelId, cache, current]);
}
