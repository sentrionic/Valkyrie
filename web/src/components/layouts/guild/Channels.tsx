import React from 'react';
import { Box, GridItem, UnorderedList, useDisclosure } from '@chakra-ui/react';
import { useParams } from 'react-router-dom';
import { useQuery } from 'react-query';
import { AccountBar } from '../AccountBar';
import { CreateChannelModal } from '../../modals/CreateChannelModal';
import { GuildMenu } from '../../menus/GuildMenu';
import { InviteModal } from '../../modals/InviteModal';
import { ChannelListItem } from '../../items/ChannelListItem';
import { cKey } from '../../../lib/utils/querykeys';
import { channelScrollbarCss } from './css/ChannelScrollerCSS';
import { useChannelSocket } from '../../../lib/api/ws/useChannelSocket';
import { getChannels } from '../../../lib/api/handler/channel';
import { RouterProps } from '../../../lib/api/models';

export const Channels: React.FC = () => {
  const { isOpen: inviteIsOpen, onOpen: inviteOpen, onClose: inviteClose } = useDisclosure();
  const { isOpen: channelIsOpen, onOpen: channelOpen, onClose: channelClose } = useDisclosure();

  const { guildId } = useParams<RouterProps>();
  const key = cKey(guildId);

  const { data } = useQuery(key, () => getChannels(guildId).then((response) => response.data));

  useChannelSocket(guildId, key);

  return (
    <>
      <GuildMenu channelOpen={channelOpen} inviteOpen={inviteOpen} />
      <GridItem
        gridColumn={2}
        gridRow="2/4"
        bg="brandGray.dark"
        overflowY="hidden"
        _hover={{ overflowY: 'auto' }}
        css={channelScrollbarCss}
      >
        {inviteIsOpen && <InviteModal isOpen={inviteIsOpen} onClose={inviteClose} />}
        {channelIsOpen && <CreateChannelModal guildId={guildId} onClose={channelClose} isOpen={channelIsOpen} />}
        <UnorderedList listStyleType="none" ml="0" mt="4">
          {data?.map((c) => (
            <ChannelListItem channel={c} guildId={guildId} key={`${c.id}`} />
          ))}
          <Box h="16" />
        </UnorderedList>
        <AccountBar />
      </GridItem>
    </>
  );
};
