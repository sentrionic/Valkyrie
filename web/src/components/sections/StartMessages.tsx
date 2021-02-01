import { Box, Flex, Heading, Text } from '@chakra-ui/react';
import React from 'react';
import { useParams } from 'react-router-dom';
import { useQuery } from 'react-query';
import { Channel } from '../../lib/api/models';
import { RouterProps } from '../../routes/Routes';

export const StartMessages: React.FC = () => {

  const { guildId, channelId } = useParams<RouterProps>();
  const { data } = useQuery<Channel[]>(`channels-${guildId}`);
  const channel = data?.find(c => c.id === channelId);

  return (
    <Flex
      alignItems='center'
      mb='2'
      justify='center'
    >
      <Box textAlign={'center'}>
        <Heading>Welcome to #{channel?.name}</Heading>
        <Text>This is the start of the #{channel?.name} channel</Text>
      </Box>
    </Flex>
  );
};
