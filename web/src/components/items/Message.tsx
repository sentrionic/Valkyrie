import { Avatar, Box, Flex, Menu, MenuButton, Text } from '@chakra-ui/react';
import React, { useState } from 'react';
import { FaEllipsisV, FaRegTrashAlt } from 'react-icons/fa';
import { MdEdit } from 'react-icons/md';
import { StyledMenuItem, StyledRedMenuItem } from '../menus/StyledMenuItem';
import { StyledMenuList } from '../menus/StyledMenuList';
import { Message as MessageResponse } from '../../lib/api/models';
import { userStore } from '../../lib/stores/userStore';
import { getTime } from '../../lib/utils/dateUtils';

interface MessageProps {
  message: MessageResponse;
}

export const Message: React.FC<MessageProps> = ({ message }) => {

  const current = userStore(state => state.current);
  const isAuthor = current?.id === message.user.id;
  const [showSettings, setShowSettings] = useState(false);

  return (
    // <Menu>
    //   {({ isOpen }) => (
    //     <>
          <Flex
            alignItems='center'
            my='2'
            mr='1'
            _hover={{ bg: 'brandGray.dark' }}
            justify='space-between'
            onMouseLeave={() => setShowSettings(false)}
            onMouseEnter={() => setShowSettings(true)}
          >
            <Flex alignItems='center'>
              <Avatar h='40px' w='40px' ml='4' src={message.user.image} />
              <Box ml='3'>
                <Flex alignItems='center'>
                  <Text>{message.user.username}</Text>
                  <Text fontSize='12px' color='brandGray.accent' ml='3'>
                    {getTime(message.createdAt)}
                  </Text>
                </Flex>
                <Text>{message.text}</Text>
              </Box>
            </Flex>
            {(isAuthor && (showSettings)) && (
              <Box  mr='2' _hover={{ cursor: "pointer" }}>
                <FaEllipsisV />
              </Box>
            )}
          </Flex>
    //       <StyledMenuList>
    //         <StyledMenuItem
    //           label={'Edit Message'}
    //           icon={MdEdit}
    //           handleClick={() => console.log('Edit')}
    //         />
    //         <StyledRedMenuItem
    //           label={'Delete Message'}
    //           icon={FaRegTrashAlt}
    //           handleClick={() => console.log('Delete')}
    //         />
    //       </StyledMenuList>
    //     </>
    //   )}
    // </Menu>
  );
};
