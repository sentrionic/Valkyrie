import React from 'react';
import { Flex, Icon, Text } from '@chakra-ui/react';
import { FiUsers } from 'react-icons/fi';
import { Link, useLocation } from 'react-router-dom';

export const FriendsListButton: React.FC = () => {

  const currentPath = `/channels/me`;
  const location = useLocation();
  const isActive = location.pathname === currentPath;

  return (
    <Link to={'/channels/me'}>
      <Flex
        m='2'
        p='3'
        align='center'
        color={isActive ? '#fff' : 'brandGray.accent'}
        _hover={{ bg: '#36393f', borderRadius: '5px', cursor: 'pointer', color: '#fff' }}
        bg={isActive ? '#393c43' : undefined}
      >
        <Icon as={FiUsers} fontSize='20px' />
        <Text fontSize='16px' ml='4' fontWeight='semibold'>
          Friends
        </Text>
      </Flex>
    </Link>
  );
};
