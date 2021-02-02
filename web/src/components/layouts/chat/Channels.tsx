import React, { useEffect } from 'react';
import { Box, GridItem, UnorderedList, useDisclosure } from '@chakra-ui/react';
import { useParams } from 'react-router-dom';
import { useQuery, useQueryClient } from 'react-query';
import { AccountBar } from '../AccountBar';
import { CreateChannelModal } from '../../modals/CreateChannelModal';
import { GuildMenu } from '../../menus/GuildMenu';
import { InviteModal } from '../../modals/InviteModal';
import { ChannelListItem } from '../../items/ChannelListItem';
import { getChannels } from '../../../lib/api/handler/guilds';
import { Channel } from '../../../lib/api/models';
import { getSocket } from '../../../lib/api/getSocket';
import { RouterProps } from '../../../routes/Routes';

export const Channels: React.FC = () => {
  const {
    isOpen: inviteIsOpen,
    onOpen: inviteOpen,
    onClose: inviteClose
  } = useDisclosure();
  const {
    isOpen: channelIsOpen,
    onOpen: channelOpen,
    onClose: channelClose
  } = useDisclosure();

  const { guildId } = useParams<RouterProps>();
  const key = `channels-${guildId}`;

  const cache = useQueryClient();

  const { data } = useQuery(key, () =>
    getChannels(guildId).then(response => response.data)
  );

  useEffect((): any => {
    const socket = getSocket();
    socket.emit('joinGuild', guildId);

    socket.on('add_channel', (newChannel: Channel) => {
      cache.setQueryData<Channel[]>(key, (_) => {
        return [...data!, newChannel];
      });
    });

    return () => {
      socket.emit('leaveRoom', guildId);
      socket.disconnect();
    };
  }, [guildId, data, key, cache]);

  return (
    <>
      <GuildMenu channelOpen={channelOpen} inviteOpen={inviteOpen} />
      <GridItem
        gridColumn={2}
        gridRow={'2/4'}
        bg='brandGray.dark'
        overflowY='hidden'
        _hover={{ overflowY: 'auto' }}
        css={{
          '&::-webkit-scrollbar': {
            width: '4px'
          },
          '&::-webkit-scrollbar-track': {
            width: '4px'
          },
          '&::-webkit-scrollbar-thumb': {
            background: '#202225',
            borderRadius: '18px'
          }
        }}
      >
        <InviteModal isOpen={inviteIsOpen} onClose={inviteClose} />
        <CreateChannelModal guildId={guildId} onClose={channelClose} isOpen={channelIsOpen} />
        <UnorderedList listStyleType='none' ml='0' mt='4'>
          {data?.map(c => (
            <ChannelListItem channel={c} guildId={guildId} key={`${c.id}`} />
          ))}
          <Box h='16' />
        </UnorderedList>
        <AccountBar />
      </GridItem>
    </>
  );
};
