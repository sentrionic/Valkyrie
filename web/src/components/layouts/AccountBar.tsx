import React from 'react';
import { Avatar, Flex, IconButton, Text } from '@chakra-ui/react';
import { RiSettings5Fill } from 'react-icons/ri';
import { Link } from 'react-router-dom';

export const AccountBar = () => {
  return (
    <Flex
      p='10px'
      pos='absolute'
      bottom='0'
      w='240px'
      bg='#292b2f'
      align='center'
      justify='space-between'
    >
      <Flex align='center'>
        <Avatar size='sm' />
        <Text ml='2'>Username</Text>
      </Flex>
      <Link to={'/account'}>
        <IconButton
          icon={<RiSettings5Fill />}
          aria-label='settings'
          size='sm'
          fontSize='20px'
          variant='ghost'
        />
      </Link>
    </Flex>
  );
};
