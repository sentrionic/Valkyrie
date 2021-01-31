import { Flex, ListItem, Text } from '@chakra-ui/react';
import { FaHashtag } from 'react-icons/fa';
import React from 'react';
import { Channel } from '../../lib/api/models';
import { Link, useLocation } from 'react-router-dom';

interface ChannelListItemProps {
  channel: Channel,
  guildId: string,
}

export const ChannelListItem: React.FC<ChannelListItemProps> = ({ channel, guildId }) => {

  const currentPath = `/channels/${guildId}/${channel.id}`;
  const location = useLocation();
  const isActive = location.pathname === currentPath;

  return (
    <Link to={currentPath}>
      <ListItem
        p='5px'
        m='0 10px'
        color={isActive ? '#fff' : 'brandGray.accent'}
        _hover={{ bg: '#36393f', borderRadius: '5px', cursor: 'pointer', color: '#fff' }}
        bg={isActive ? '#393c43' : undefined}
        mb='2px'
      >
        <Flex align='center'>
          <FaHashtag />
          <Text ml='2'>{channel.name}</Text>
        </Flex>
      </ListItem>
    </Link>
  );
};
