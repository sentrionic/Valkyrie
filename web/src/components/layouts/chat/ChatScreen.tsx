import React, { useEffect, useState } from 'react';
import { GridItem, Flex, Box, Spinner } from '@chakra-ui/react';
import { Message } from '../../items/Message';
import { StartMessages } from '../../sections/StartMessages';
import { scrollbarCss } from '../../../lib/utils/theme';
import { useParams } from 'react-router-dom';
import { getMessages } from '../../../lib/api/handler/messages';
import InfiniteScroll from 'react-infinite-scroll-component';
import { Message as MessageResponse } from '../../../lib/api/models';

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
      setHasMore(messages.length === 35);
    };
    fetchData();
  }, [channelId]);

  const fetchMore = async () => {
    const cursor = data[data.length - 1].createdAt;
    const { data: messages } = await getMessages(channelId, cursor);
    setData([...data, ...messages]);
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
      <Box as={InfiniteScroll}
           css={{
             '&::-webkit-scrollbar': {
               width: '0',
             },
           }}
           dataLength={data.length}
           next={() => fetchMore()}
           style={{ display: 'flex', flexDirection: 'column-reverse' }}
           inverse={true}
           hasMore={hasMore}
           loader={
             <Flex align={'center'} justify={'center'} h={'50px'}>
               <Spinner />
             </Flex>
           }
           scrollableTarget='chatGrid'
      >
        {data.map(m => <Message key={m.id} message={m} />)}
      </Box>
      <StartMessages />
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

