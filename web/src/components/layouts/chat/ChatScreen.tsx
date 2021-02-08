import React, { useEffect, useState } from 'react';
import { Box, Divider, Flex, GridItem, Spinner, Text } from '@chakra-ui/react';
import InfiniteScroll from 'react-infinite-scroll-component';
import { InfiniteData, useInfiniteQuery, useQueryClient } from 'react-query';
import { Message } from '../../items/message/Message';
import { StartMessages } from '../../sections/StartMessages';
import { scrollbarCss } from '../../../lib/utils/theme';
import { useParams } from 'react-router-dom';
import { getMessages } from '../../../lib/api/handler/messages';
import { Message as MessageResponse } from '../../../lib/api/models';
import { getSocket } from '../../../lib/api/getSocket';
import { RouterProps } from '../../../routes/Routes';
import { channelStore } from '../../../lib/stores/channelStore';
import { userStore } from '../../../lib/stores/userStore';
import { checkNewDay, formatDivider, getTimeDifference } from '../../../lib/utils/dateUtils';

export const ChatScreen: React.FC = () => {

  const { channelId } = useParams<RouterProps>();
  const [hasMore, setHasMore] = useState(true);
  const qKey = `messages-${channelId}`;
  const cache = useQueryClient();
  const current = userStore(state => state.current);
  const store = channelStore();

  const { data, isLoading, fetchNextPage } = useInfiniteQuery<MessageResponse[]>(qKey, async ({ pageParam = null }) => {
    const { data } = await getMessages(channelId, pageParam);
    if (data.length !== 35) setHasMore(false);
    return data;
  }, {
    staleTime: 0,
    cacheTime: 0,
    getNextPageParam: lastPage => hasMore && lastPage.length ? lastPage[lastPage.length - 1].createdAt : ''
  });

  useEffect((): any => {

    store.reset();
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

    socket.on('addToTyping', (username: string) => {
      if (username !== current?.username)
        store.addTyping(username);
    });

    socket.on('removeFromTyping', (username: string) => {
      if (username !== current?.username)
        store.removeTyping(username);
    });

    return () => {
      socket.emit('leaveRoom', channelId);
      socket.disconnect();
    };
    // eslint-disable-next-line
  }, [channelId, data, cache, qKey, current]);

  if (isLoading) {
    return (
      <ChatGrid>
        <Flex align={'center'} justify={'center'} h={'full'} />
      </ChatGrid>
    );
  }

  const checkIfWithinTime = (message1: MessageResponse, message2: MessageResponse): boolean => {
    if (message1.user.id !== message2.user.id) return false;
    if (message1.createdAt === message2.createdAt) return false;
    return getTimeDifference(message1.createdAt, message2.createdAt) <= 5;
  };

  const messages = data ? data!.pages.map(p => p.map(p => p)).flat() : [];

  return (
    <ChatGrid>
      <Box h={'10px'} mt={4} />
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
        {messages.map((m, i) =>
          <React.Fragment key={i}>
            <Message
              key={m.id}
              message={m}
              isCompact={
                checkIfWithinTime(
                  m,
                  messages[Math.min(i + 1, messages.length - 1)]
                )}
            />
            {checkNewDay(m.createdAt, messages[Math.min(i + 1, messages.length - 1)].createdAt) &&
            <Flex textAlign='center' align='center' mt={'2'} mx={'4'} key={m.createdAt}>
              <Divider />
              <Text
                w={['75%', '75%', '75%', '40%', '25%']}
                fontSize={'12px'}
                color={'brandGray.accent'}
              >
                {formatDivider(m.createdAt)}
              </Text>
              <Divider />
            </Flex>
            }
          </React.Fragment>
        )
        }
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

