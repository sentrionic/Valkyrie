import React, { useEffect, useState } from 'react';
import { GridItem, Flex, Box, Spinner } from '@chakra-ui/react';
import { Message } from '../../items/Message';
import { StartMessages } from '../../sections/StartMessages';
import { scrollbarCss } from '../../../lib/utils/theme';
import { useParams } from 'react-router-dom';
import { getMessages } from '../../../lib/api/handler/messages';
import InfiniteScroll from 'react-infinite-scroll-component';
import { Message as MessageResponse } from '../../../lib/api/models';
import socketIOClient from 'socket.io-client';

interface RouterProps {
  channelId: string;
}

export const ChatScreen: React.FC = () => {

  const { channelId } = useParams<RouterProps>();
  const [isLoading, setIsLoading] = useState(true);
  const [data, setData] = useState<MessageResponse[]>([]);
  const [hasMore, setHasMore] = useState(true);

  useEffect(() => {
    setData([]);
    setIsLoading(true);
    const fetchData = async () => {
      const { data: messages } = await getMessages(channelId);
      setData(messages);
      setIsLoading(false);
    };
    fetchData();
  }, [channelId]);

  useEffect((): any => {

    if (data.length === 0) setHasMore(false);

    const socket = socketIOClient(process.env.REACT_APP_API_WS!);
    socket.emit('joinChannel', channelId);

    socket.on('new_message', (newMessage: MessageResponse) => {
      setData([newMessage, ...data]);
    });

    socket.on('edit_message', (editMessage: MessageResponse) => {
      const index = data.findIndex(m => m.id === editMessage.id);
      const messages = [...data];
      messages[index] = editMessage;
      setData(messages);
    });

    socket.on('delete_message', (toBeRemoved: MessageResponse) => {
      setData(data.filter(m => m.id !== toBeRemoved.id));
    });

    return () => {
      socket.emit('leaveRoom', channelId);
      socket.disconnect();
    };
  }, [channelId, setData, data]);

  const fetchMore = async () => {
    const cursor = data[data.length - 1].createdAt;
    const { data: messages } = await getMessages(channelId, cursor);
    setData([...data, ...messages]);
    setHasMore(messages.length === 35);
  };

  if (isLoading) {
    return (
      <ChatGrid>
        <Flex align={'center'} justify={'center'} h={'full'} />
      </ChatGrid>
    );
  }

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
        dataLength={data.length}
        next={() => fetchMore()}
        style={{ display: 'flex', flexDirection: 'column-reverse' }}
        inverse={true}
        hasMore={hasMore}
        loader={
          data.length > 0 &&
          <Flex align={'center'} justify={'center'} h={'50px'}>
            <Spinner />
          </Flex>
        }
        scrollableTarget='chatGrid'
      >
        {data.map(m => <Message key={m.id} message={m} />)}
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

