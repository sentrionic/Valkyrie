import { Avatar, Box, Divider, Flex, Heading, Text } from '@chakra-ui/react';
import React from 'react';
import { useParams } from 'react-router-dom';
import { RouterProps } from '../../routes/Routes';
import { useGetCurrentChannel } from '../../lib/utils/hooks/useGetCurrentChannel';
import { cKey } from '../../lib/utils/querykeys';
import { useGetCurrentDM } from '../../lib/utils/hooks/useGetCurrentDM';

export const StartMessages: React.FC = () => {
  const { guildId } = useParams<RouterProps>();
  return (guildId === undefined) ? <DMStartMessages /> : <ChannelStartMessages />;
};

const ChannelStartMessages: React.FC = () => {

  const { guildId, channelId } = useParams<RouterProps>();
  const channel = useGetCurrentChannel(channelId, cKey(guildId));

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

const DMStartMessages: React.FC = () => {

  const { channelId } = useParams<RouterProps>();
  const channel = useGetCurrentDM(channelId);

  return (
    <Box m='4'>
      <Avatar h={'80px'} w={'80px'} src={channel?.user.image} />
      <Heading mt={2}>{channel?.user.username}</Heading>
      <Text textColor={'brandGray.accent'}>
        This is the beginning of your direct message history with
        @{channel?.user.username}
      </Text>
      <Divider mt={2} />
    </Box>
  );
};

