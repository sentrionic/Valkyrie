import { Avatar, AvatarBadge, Flex, IconButton, ListItem, Text } from '@chakra-ui/react';
import React from 'react';
import { FaEllipsisV } from 'react-icons/fa';

export const FriendsListItem: React.FC = () => {
  return (
    <ListItem
      p='3'
      mx='3'
      _hover={{
        bg: 'brandGray.dark',
        borderRadius: '5px',
        cursor: 'pointer',
      }}
    >
      <Flex align='center' justify='space-between'>
        <Flex align='center'>
          <Avatar size='sm'>
            <AvatarBadge boxSize='1.25em' bg='green.500' />
          </Avatar>
          <Text ml='2'>Username</Text>
        </Flex>
        <IconButton
          icon={<FaEllipsisV />}
          borderRadius='50%'
          aria-label='remove friend'
        />
      </Flex>
    </ListItem>
  );
};
