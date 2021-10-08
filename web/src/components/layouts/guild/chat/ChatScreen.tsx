import React, { useState } from 'react';
import { Box, Flex, Spinner } from '@chakra-ui/react';
import InfiniteScroll from 'react-infinite-scroll-component';
import { useInfiniteQuery } from 'react-query';
import { useParams } from 'react-router-dom';
import { Message } from '../../../items/message/Message';
import { StartMessages } from '../../../sections/StartMessages';
import { getMessages } from '../../../../lib/api/handler/messages';
import { checkNewDay, getTimeDifference } from '../../../../lib/utils/dateUtils';
import { guildScrollbarCss } from '../css/GuildScrollerCSS';
import { useMessageSocket } from '../../../../lib/api/ws/useMessageSocket';
import { DateDivider } from '../../../sections/DateDivider';
import { ChatGrid } from './ChatGrid';
import { RouterProps } from '../../../../lib/models/routerProps';
import { Message as MessageResponse } from '../../../../lib/models/message';

export const ChatScreen: React.FC = () => {
  const { channelId } = useParams<RouterProps>();
  const [hasMore, setHasMore] = useState(true);
  const qKey = `messages-${channelId}`;

  const { data, isLoading, fetchNextPage } = useInfiniteQuery<MessageResponse[]>(
    qKey,
    async ({ pageParam = null }) => {
      const { data: messageData } = await getMessages(channelId, pageParam);
      if (messageData.length !== 35) setHasMore(false);
      return messageData;
    },
    {
      staleTime: 0,
      cacheTime: 0,
      getNextPageParam: (lastPage) => (hasMore && lastPage.length ? lastPage[lastPage.length - 1].createdAt : ''),
    }
  );

  useMessageSocket(channelId, qKey);

  if (isLoading) {
    return (
      <ChatGrid>
        <Flex align="center" justify="center" h="full">
          <Spinner size="xl" thickness="4px" />
        </Flex>
      </ChatGrid>
    );
  }

  const checkIfWithinTime = (message1: MessageResponse, message2: MessageResponse): boolean => {
    if (message1.user.id !== message2.user.id) return false;
    if (message1.createdAt === message2.createdAt) return false;
    return getTimeDifference(message1.createdAt, message2.createdAt) <= 5;
  };

  const messages = data ? data!.pages.map((p) => p.map((mr) => mr)).flat() : [];

  return (
    <ChatGrid>
      <Box h="10px" mt={4} />
      <Box
        as={InfiniteScroll as any}
        css={guildScrollbarCss}
        dataLength={messages.length}
        next={() => fetchNextPage()}
        style={{
          display: 'flex',
          flexDirection: 'column-reverse',
        }}
        inverse
        hasMore={hasMore}
        loader={
          messages.length > 0 && (
            <Flex align="center" justify="center" h="50px">
              <Spinner />
            </Flex>
          )
        }
        scrollableTarget="chatGrid"
      >
        {messages.map((m, i) => (
          <React.Fragment key={m.id}>
            <Message message={m} isCompact={checkIfWithinTime(m, messages[Math.min(i + 1, messages.length - 1)])} />
            {checkNewDay(m.createdAt, messages[Math.min(i + 1, messages.length - 1)].createdAt) && (
              <DateDivider date={m.createdAt} />
            )}
          </React.Fragment>
        ))}
      </Box>
      {!hasMore && <StartMessages />}
    </ChatGrid>
  );
};
