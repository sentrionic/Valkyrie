import { Flex, ListItem, Text } from '@chakra-ui/react';
import { FaHashtag } from 'react-icons/fa';
import React from 'react';

export const ChannelListItem: React.FC = () => {
  return (
    <ListItem
      p='5px'
      m='0 10px'
      _hover={{ bg: '#36393f', borderRadius: '5px', cursor: 'pointer' }}
    >
      <Flex align='center'>
        <FaHashtag />
        <Text ml='2'>general</Text>
      </Flex>
    </ListItem>
  );
};
