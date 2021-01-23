import { Flex, Heading, Icon, Menu, MenuButton, MenuItem, MenuList } from '@chakra-ui/react';
import { FiChevronDown, FiX } from 'react-icons/fi';
import React from 'react';

interface GuildMenuProps {
  channelOpen: () => void;
  inviteOpen: () => void;
}

export const GuildMenu: React.FC<GuildMenuProps> = ({ channelOpen, inviteOpen }) => {
  return (
    <Menu placement='bottom-end'>
      {({ isOpen }) => (
        <>
          <Flex
            justify='space-between'
            align='center'
            boxShadow='md'
            p='10px'
          >
            <Heading fontSize='20px'>Harmony</Heading>
            <MenuButton>
              <Icon as={!isOpen ? FiChevronDown : FiX} />
            </MenuButton>
          </Flex>
          <MenuList bg='#18191c'>
            <MenuItem onClick={channelOpen}>Create Channel</MenuItem>
            <MenuItem onClick={inviteOpen}>Invite People</MenuItem>
            <MenuItem>Leave Server</MenuItem>
          </MenuList>
        </>
      )}
    </Menu>
  );
};
