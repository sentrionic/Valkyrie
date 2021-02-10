import React from 'react';
import { Avatar, Flex, IconButton, Text } from '@chakra-ui/react';
import { RiSettings5Fill } from 'react-icons/ri';
import { Link } from 'react-router-dom';
import { userStore } from '../../lib/stores/userStore';

export const AccountBar: React.FC = () => {

  const user = userStore(state => state.current);

  return (
    <Flex
      p="10px"
      pos="absolute"
      bottom="0"
      w="240px"
      bg="#292b2f"
      align="center"
      justify="space-between"
    >
      <Flex align="center">
        <Avatar size="sm" src={user?.image} />
        <Text ml="2">{user?.username}</Text>
      </Flex>
      <Link to={'/account'}>
        <IconButton
          icon={<RiSettings5Fill />}
          aria-label="settings"
          size="sm"
          fontSize="20px"
          variant="ghost"
        />
      </Link>
    </Flex>
  );
};
