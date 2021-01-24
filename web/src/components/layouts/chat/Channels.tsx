import React from 'react';
import { GridItem, UnorderedList, useDisclosure } from '@chakra-ui/react';
import { CreateChannelModal } from '../../modals/CreateChannelModal';
import { GuildMenu } from '../../menus/GuildMenu';
import { InviteModal } from '../../modals/InviteModal';
import { ChannelListItem } from '../../items/ChannelListItem';
import { AccountBar } from '../AccountBar';

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

  return (
    <GridItem
      gridColumn={2}
      gridRow={'1 / 4'}
      bg="brandGray.dark"
      overflowY="hidden"
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
      <GuildMenu channelOpen={channelOpen} inviteOpen={inviteOpen} />
      <InviteModal isOpen={inviteIsOpen} onClose={inviteClose} />
      <CreateChannelModal onClose={channelClose} isOpen={channelIsOpen} />
      <UnorderedList listStyleType="none" ml="0" mt="4">
        {[...Array(15)].map((x, i) => (
          <ChannelListItem key={`${i}`} />
        ))}
      </UnorderedList>
      <AccountBar />
    </GridItem>
  );
};
