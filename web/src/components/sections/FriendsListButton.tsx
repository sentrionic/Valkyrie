import React from 'react';
import { Flex, Icon, Text } from '@chakra-ui/react';
import { FiUsers } from 'react-icons/fi';
import { Link } from 'react-router-dom';

export const FriendsListButton: React.FC = () => {
  return (
    <Link to={'/channels/me'}>
      <Flex
        m='2'
        p='3'
        align='center'
        _hover={{ cursor: 'pointer', bg: 'brandGray.light' }}
      >
        <Icon as={FiUsers} fontSize='20px' />
        <Text fontSize='16px' ml='4' fontWeight='semibold'>
          Friends
        </Text>
      </Flex>
    </Link>
  );
};
