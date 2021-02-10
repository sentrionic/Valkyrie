import React, { useState } from 'react';
import { Flex, Icon, ListItem, Text, useDisclosure } from '@chakra-ui/react';
import { FaHashtag, FaUserLock } from 'react-icons/fa';
import { MdSettings } from 'react-icons/md';
import { Channel } from '../../lib/api/models';
import { Link, useLocation } from 'react-router-dom';
import { userStore } from '../../lib/stores/userStore';
import { ChannelSettingsModal } from '../modals/ChannelSettingsModal';
import { useGetCurrentGuild } from '../../lib/utils/hooks/useGetCurrentGuild';

interface ChannelListItemProps {
  channel: Channel,
  guildId: string,
}

export const ChannelListItem: React.FC<ChannelListItemProps> = ({ channel, guildId }) => {

  const currentPath = `/channels/${guildId}/${channel.id}`;
  const location = useLocation();
  const isActive = location.pathname === currentPath;
  const [showSettings, setShowSettings] = useState(false);

  const current = userStore(state => state.current);
  const guild = useGetCurrentGuild(guildId);

  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <Link to={currentPath}>
      <ListItem
        p='5px'
        m='0 10px'
        color={isActive ? '#fff' : 'brandGray.accent'}
        _hover={{ bg: '#36393f', borderRadius: '5px', cursor: 'pointer', color: '#fff' }}
        bg={isActive ? '#393c43' : undefined}
        mb='2px'
        onMouseLeave={() => setShowSettings(false)}
        onMouseEnter={() => setShowSettings(true)}
      >
        <Flex align='center' justify={'space-between'}>
          <Flex align='center'>
            <Icon as={channel.isPublic ? FaHashtag : FaUserLock} />
            <Text ml='2'>{channel.name}</Text>
          </Flex>
          {(current?.id === guild?.ownerId && (showSettings || isOpen)) &&
          <>
            <Icon
              as={MdSettings}
              color={'brandGray.accent'}
              fontSize={'12px'}
              _hover={{ color: '#fff' }}
              onClick={(e) => {
                e.preventDefault();
                onOpen();
              }}
            />
            {isOpen &&
              <ChannelSettingsModal
                guildId={guildId}
                channelId={channel.id}
                isOpen={isOpen}
                onClose={onClose}
              />
            }
          </>
          }
        </Flex>
      </ListItem>
    </Link>
  );
};
