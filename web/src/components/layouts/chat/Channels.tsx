import React from 'react';
import { Box, GridItem, UnorderedList, useDisclosure } from '@chakra-ui/react';
import { CreateChannelModal } from '../../modals/CreateChannelModal';
import { GuildMenu } from '../../menus/GuildMenu';
import { InviteModal } from '../../modals/InviteModal';
import { ChannelListItem } from '../../items/ChannelListItem';
import { AccountBar } from '../AccountBar';
import { useQuery } from 'react-query';
import { getChannels } from '../../../lib/api/handler/guilds';
import { useParams } from 'react-router-dom';

interface RouterProps {
  guildId: string;
}

export const Channels: React.FC = () => {
  const {
    isOpen: inviteIsOpen,
    onOpen: inviteOpen,
    onClose: inviteClose,
  } = useDisclosure();
  const {
    isOpen: channelIsOpen,
    onOpen: channelOpen,
    onClose: channelClose,
  } = useDisclosure();

  const { guildId } = useParams<RouterProps>();

  const { data } = useQuery(`channels-${guildId}`, () =>
    getChannels(guildId).then(response => response.data),
    {
      refetchOnWindowFocus: false,
    },
  );

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
            width: '4px',
          },
          '&::-webkit-scrollbar-track': {
            width: '4px',
          },
          '&::-webkit-scrollbar-thumb': {
            background: '#202225',
            borderRadius: '18px',
          },
        }}
      >
        <InviteModal isOpen={inviteIsOpen} onClose={inviteClose} />
        <CreateChannelModal onClose={channelClose} isOpen={channelIsOpen} />
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
