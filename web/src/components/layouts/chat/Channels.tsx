import React, { useEffect } from 'react';
import { Box, GridItem, UnorderedList, useDisclosure } from '@chakra-ui/react';
import { useLocation, useParams, useHistory } from 'react-router-dom';
import { useQuery, useQueryClient } from 'react-query';
import { AccountBar } from '../AccountBar';
import { CreateChannelModal } from '../../modals/CreateChannelModal';
import { GuildMenu } from '../../menus/GuildMenu';
import { InviteModal } from '../../modals/InviteModal';
import { ChannelListItem } from '../../items/ChannelListItem';
import { getChannels } from '../../../lib/api/handler/guilds';
import { Channel, Guild } from '../../../lib/api/models';
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
  const location = useLocation();
  const history = useHistory();

  const cache = useQueryClient();

  const { data } = useQuery(key, () =>
    getChannels(guildId).then(response => response.data)
  );

  const { data: guildData } = useQuery<Guild[]>('guilds');
  const guild = guildData?.find(g => g.id === guildId);

  useEffect((): any => {
    const socket = getSocket();
    socket.emit('joinGuild', guildId);

    socket.on('add_channel', (newChannel: Channel) => {
      cache.setQueryData<Channel[]>(key, (_) => {
        return [...data!, newChannel];
      });
    });

    socket.on('edit_channel', (editedChannel: Channel) => {
      cache.setQueryData<Channel[]>(key, (d) => {
        const index = d!.findIndex(c => c.id === editedChannel.id);
        if (index !== -1) {
          d![index] = editedChannel;
        } else if (editedChannel.isPublic) {
          d!.push(editedChannel);
        }
        return d!;
      });
    });

    socket.on('delete_channel', (deleteId: string) => {
      cache.setQueryData<Channel[]>(key, (d) => {
        const currentPath = `/channels/${guildId}/${deleteId}`;
        if (location.pathname === currentPath && guild) {
          if (deleteId === guild.default_channel_id) {
            history.replace('/channels/me')
          } else {
            history.replace(`${guild.default_channel_id}`);
          }
        }
        return d!.filter(c => c.id !== deleteId);
      });
    });

    return () => {
      socket.emit('leaveRoom', guildId);
      socket.disconnect();
    };
  }, [guildId, data, key, cache, history, location, guild]);

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
