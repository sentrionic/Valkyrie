import { Avatar, AvatarBadge, Flex, ListItem, Text } from '@chakra-ui/react';
import React from 'react';
import { DMChannel } from '../../lib/api/models';
import { Link } from 'react-router-dom';

interface DMListItemProps {
  dm: DMChannel
}

export const DMListItem: React.FC<DMListItemProps> = ({ dm }) => {
  return (
    <Link to={`/channels/me/${dm.id}`}>
      <ListItem
        p='2'
        mx='2'
        _hover={{ bg: '#36393f', borderRadius: '5px', cursor: 'pointer' }}
      >
        <Flex align='center'>
          <Avatar size='sm' src={dm.user.image}>
            <AvatarBadge boxSize='1.25em' bg={ dm.user.isOnline ? 'green.500' : 'gray.500'}  />
          </Avatar>
          <Text ml='2'>{dm.user.username}</Text>
        </Flex>
      </ListItem>
    </Link>
  );
};
