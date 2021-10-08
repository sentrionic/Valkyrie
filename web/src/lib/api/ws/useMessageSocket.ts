import { useEffect } from 'react';
import { InfiniteData, useQueryClient } from 'react-query';
import { getSocket } from '../getSocket';
import { userStore } from '../../stores/userStore';
import { channelStore } from '../../stores/channelStore';
import { Message } from '../../models/message';

type WSMessage =
  | { action: 'new_message' | 'edit_message'; data: Message }
  | { action: 'addToTyping' | 'removeFromTyping' | 'delete_message'; data: string };

export function useMessageSocket(channelId: string, key: string): void {
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
          cache.setQueryData<InfiniteData<Message[]>>(key, (d) => {
            d!.pages[0].unshift(response.data);
            return d!;
          });
          break;
        }

        case 'edit_message': {
          const editMessage = response.data;
          cache.setQueryData<InfiniteData<Message[]>>(key, (d) => {
            let index = -1;
            let editId = -1;
            const data = d!;
            data.pages.forEach((p, i) => {
              editId = p.findIndex((m) => m.id === editMessage.id);
              if (editId !== -1) index = i;
            });
            if (index !== -1 && editId !== -1) {
              data.pages[index][editId].text = editMessage.text;
              data.pages[index][editId].updatedAt = editMessage.updatedAt;
            }
            return data;
          });
          break;
        }

        case 'delete_message': {
          const messageId = response.data;
          cache.setQueryData<InfiniteData<Message[]>>(key, (d) => {
            let index = -1;
            const data = d!;
            data.pages.forEach((p, i) => {
              if (p.findIndex((m) => m.id === messageId) !== -1) index = i;
            });
            if (index !== -1) data.pages[index] = data.pages[index].filter((m) => m.id !== messageId);
            return data;
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
  }, [channelId, cache, key, current]);
}
