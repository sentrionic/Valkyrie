import React from 'react';
import { Avatar, Flex, IconButton, Text, Tooltip, useClipboard } from '@chakra-ui/react';
import { RiSettings5Fill } from 'react-icons/ri';
import { Link } from 'react-router-dom';
import { userStore } from '../../lib/stores/userStore';

export const AccountBar: React.FC = () => {

  const user = userStore(state => state.current);
  const { hasCopied, onCopy } = useClipboard(user?.id || "");

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
      <Tooltip
        hasArrow
        label={hasCopied ? 'Copied!' : "Click to copy ID"}
        placement={"top"}
        bg={hasCopied ? '#43b581' : '#18191c'}
        color={"white"}
        closeOnClick={false}
      >
        <Flex align="center" w={"full"} mr={2} _hover={{ cursor: 'pointer' }} onClick={onCopy}>
          <Avatar size="sm" src={user?.image} />
          <Text ml="2">{user?.username}</Text>
        </Flex>
      </Tooltip>
      <Link to={'/account'}>
        <Tooltip
          hasArrow
          label={'User Settings'}
          placement={"top"}
          bg={'#18191c'}
          color={"white"}
        >
        <IconButton
          icon={<RiSettings5Fill />}
          aria-label="settings"
          size="sm"
          fontSize="20px"
          variant="ghost"
        />
        </Tooltip>
      </Link>
    </Flex>
  );
};
