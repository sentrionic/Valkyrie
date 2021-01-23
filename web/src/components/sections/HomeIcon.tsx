import React from 'react';
import { Avatar, Divider, Flex } from '@chakra-ui/react';
import { Link } from 'react-router-dom';

export const HomeIcon: React.FC = () => {
  return (
    <Flex direction='column' my='2' align='center'>
      <Link to='/channels/me'>
        <Avatar
          src={`${process.env.PUBLIC_URL}/icon.png`}
          size='md'
          _hover={{ cursor: 'pointer' }}
        />
      </Link>
      <Divider mt='2' w='40px' />
    </Flex>
  );
};
