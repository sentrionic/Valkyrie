import React, { useEffect, useState } from 'react';
import { GridItem, Flex, Box, Spinner } from '@chakra-ui/react';
import { Message } from '../../items/Message';
import { StartMessages } from '../../sections/StartMessages';
import { scrollbarCss } from '../../../lib/utils/theme';
import { useParams } from 'react-router-dom';
import { getMessages } from '../../../lib/api/handler/messages';
import InfiniteScroll from 'react-infinite-scroll-component';
import { Message as MessageResponse } from '../../../lib/api/models';
import { InfiniteData, useInfiniteQuery, useQueryClient } from 'react-query';
import { getSocket } from '../../../lib/api/getSocket';
import { RouterProps } from '../../../routes/Routes';

export const ChatScreen: React.FC = () => {

  const { channelId } = useParams<RouterProps>();
  const [hasMore, setHasMore] = useState(true);
  const qKey = `messages-${channelId}`;
  const cache = useQueryClient();

  const { data, isLoading, fetchNextPage } = useInfiniteQuery<MessageResponse[]>(qKey, async ({ pageParam = null }) => {
    const { data } = await getMessages(channelId, pageParam);
    if (data.length !== 35) setHasMore(false);
    return data;
  }, {
    staleTime: 0,
    cacheTime: 0,
    getNextPageParam: lastPage => hasMore && lastPage.length ? lastPage[lastPage.length - 1].createdAt : '',
    refetchOnWindowFocus: false,
    refetchIntervalInBackground: false
  });

  useEffect((): any => {

    const socket = getSocket();
    socket.emit('joinChannel', channelId);

    socket.on('new_message', (newMessage: MessageResponse) => {
      cache.setQueryData<InfiniteData<MessageResponse[]>>(qKey, (d) => {
        d!.pages[0].unshift(newMessage);
        return d!;
      });
    });

    socket.on('edit_message', (editMessage: MessageResponse) => {
      cache.setQueryData<InfiniteData<MessageResponse[]>>(qKey, (d) => {
        let index = -1;
        let editId = -1;
        d!.pages.forEach((p, i) => {
          editId = p.findIndex(m => m.id === editMessage.id);
          if (editId !== -1) index = i;
        });

        if (index !== -1 && editId !== -1) {
          d!.pages[index][editId] = editMessage;
        }
        return d!;
      });
    });

    socket.on('delete_message', (toBeRemoved: MessageResponse) => {
      cache.setQueryData<InfiniteData<MessageResponse[]>>(qKey, (d) => {
        let index = -1;
        d!.pages.forEach((p, i) => {
          if (p.findIndex(m => m.id === toBeRemoved.id) !== -1) index = i;
        });
        if (index !== -1) d!.pages[index] = d!.pages[index].filter(m => m.id !== toBeRemoved.id);
        return d!;
      });
    });

    return () => {
      socket.emit('leaveRoom', channelId);
      socket.disconnect();
    };
  }, [channelId, data, cache, qKey]);

  if (isLoading) {
    return (
      <ChatGrid>
        <Flex align={'center'} justify={'center'} h={'full'} />
      </ChatGrid>
    );
  }

  const messages = data ? data!.pages.map(p => p.map(p => p)).flat() : [];

  return (
    <ChatGrid>
      <Box h={'10px'} />
      <Box
        as={InfiniteScroll}
        css={{
          '&::-webkit-scrollbar': {
            width: '0'
          }
        }}
        dataLength={messages.length}
        next={() => fetchNextPage()}
        style={{ display: 'flex', flexDirection: 'column-reverse' }}
        inverse={true}
        hasMore={hasMore}
        loader={
          messages.length > 0 &&
          <Flex align={'center'} justify={'center'} h={'50px'}>
            <Spinner />
          </Flex>
        }
        scrollableTarget='chatGrid'
      >
        {messages.map(m => <Message key={m.id} message={m} />)}
      </Box>
      {!hasMore && <StartMessages />}
    </ChatGrid>
  );
};

const ChatGrid: React.FC = ({ children }) =>
  <GridItem
    id={'chatGrid'}
    gridColumn={3}
    gridRow={'2'}
    bg='brandGray.light'
    mr='5px'
    display='flex'
    flexDirection='column-reverse'
    overflowY='auto'
    css={scrollbarCss}
  >
    {children}
  </GridItem>;

