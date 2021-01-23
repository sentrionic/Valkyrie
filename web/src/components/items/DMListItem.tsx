import { Avatar, AvatarBadge, Flex, ListItem, Text } from '@chakra-ui/react';
import React from 'react';

export const DMListItem: React.FC = () => {
  return (
    <ListItem
      p='2'
      mx='2'
      _hover={{ bg: '#36393f', borderRadius: '5px', cursor: 'pointer' }}
    >
      <Flex align='center'>
        <Avatar size='sm'>
          <AvatarBadge boxSize='1.25em' bg='green.500' />
        </Avatar>
        <Text ml='2'>sentrionic</Text>
      </Flex>
    </ListItem>
  );
};
