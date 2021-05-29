import { useEffect } from 'react';
import { getSocket } from '../getSocket';
import { Message as MessageResponse } from '../models';
import { InfiniteData, useQueryClient } from 'react-query';
import { userStore } from '../../stores/userStore';
import { channelStore } from '../../stores/channelStore';

export function useMessageSocket(channelId: string, key: string) {
  const current = userStore((state) => state.current);
  const store = channelStore();
  const cache = useQueryClient();

  useEffect((): any => {
    store.reset();
    const socket = getSocket();
    socket.emit('joinChannel', channelId);

    socket.on('new_message', (newMessage: MessageResponse) => {
      cache.setQueryData<InfiniteData<MessageResponse[]>>(key, (d) => {
        d!.pages[0].unshift(newMessage);
        return d!;
      });
    });

    socket.on('edit_message', (editMessage: MessageResponse) => {
      cache.setQueryData<InfiniteData<MessageResponse[]>>(key, (d) => {
        let index = -1;
        let editId = -1;
        d!.pages.forEach((p, i) => {
          editId = p.findIndex((m) => m.id === editMessage.id);
          if (editId !== -1) index = i;
        });

        if (index !== -1 && editId !== -1) {
          d!.pages[index][editId].text = editMessage.text;
          d!.pages[index][editId].updatedAt = editMessage.updatedAt;
        }
        return d!;
      });
    });

    socket.on('delete_message', (toBeRemoved: MessageResponse) => {
      cache.setQueryData<InfiniteData<MessageResponse[]>>(key, (d) => {
        let index = -1;
        d!.pages.forEach((p, i) => {
          if (p.findIndex((m) => m.id === toBeRemoved.id) !== -1) index = i;
        });
        if (index !== -1) d!.pages[index] = d!.pages[index].filter((m) => m.id !== toBeRemoved.id);
        return d!;
      });
    });

    socket.on('addToTyping', (username: string) => {
      if (username !== current?.username) store.addTyping(username);
    });

    socket.on('removeFromTyping', (username: string) => {
      if (username !== current?.username) store.removeTyping(username);
    });

    return () => {
      socket.emit('leaveRoom', channelId);
      socket.disconnect();
    };
    // eslint-disable-next-line
  }, [channelId, cache, key, current]);
}
