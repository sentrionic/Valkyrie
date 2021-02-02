import React, { useState } from 'react';
import { Message as MessageResponse } from '../../lib/api/models';
import { userStore } from '../../lib/stores/userStore';
import { Avatar, Box, Flex, Icon, Text, useDisclosure } from '@chakra-ui/react';
import { Item, Menu, theme, useContextMenu } from 'react-contexify';
import { getTime } from '../../lib/utils/dateUtils';
import { MdEdit } from 'react-icons/md';
import { FaEllipsisV, FaRegTrashAlt } from 'react-icons/fa';
import { DeleteMessageModal } from '../modals/DeleteMessageModal';
import { EditMessageModal } from '../modals/EditMessageModal';
import './css/ContextMenu.css';

interface MessageProps {
  message: MessageResponse;
}

export const Message: React.FC<MessageProps> = ({ message }) => {

  const [showSettings, setShowSettings] = useState(false);
  const current = userStore(state => state.current);
  const isAuthor = current?.id === message.user.id;

  const { isOpen: isDeleteOpen, onOpen: onDeleteOpen, onClose: onDeleteClose } = useDisclosure();
  const { isOpen: isEditOpen, onOpen: onEditOpen, onClose: onEditClose } = useDisclosure();

  const { show } = useContextMenu({
    id: message.id
  });

  return (
    <>
      <Flex
        alignItems='center'
        my='2'
        mr='1'
        _hover={{ bg: '#32353b' }}
        justify='space-between'
        onContextMenu={(e) => {
          if (isAuthor) show(e);
        }}
        onMouseLeave={() => setShowSettings(false)}
        onMouseEnter={() => setShowSettings(true)}
      >
        <Flex alignItems='center'>
          <Avatar h='40px' w='40px' ml='4' src={message.user.image} />
          <Box ml='3'>
            <Flex alignItems='center'>
              <Text>{message.user.username}</Text>
              <Text fontSize='12px' color='brandGray.accent' ml='2'>
                {getTime(message.createdAt)}
              </Text>
            </Flex>
            <Flex alignItems={'center'}>
              <Text>{message.text}</Text>
              {message.createdAt !== message.updatedAt &&
              <Text fontSize={'10px'} ml={'1'} color={'#72767d'}>(edited)</Text>}
            </Flex>
          </Box>
        </Flex>
        {(isAuthor && (showSettings)) && (
          <Box onClick={show} mr='2' _hover={{ cursor: 'pointer' }}>
            <FaEllipsisV />
          </Box>
        )}
      </Flex>
      {isAuthor &&
      <>
        <Menu id={message.id} theme={theme.dark}>
          <Item className={'menu-item'} onClick={onEditOpen}>
            <Flex align='center' justify='space-between' w='full'>
              <Text>Edit Message</Text>
              <Icon as={MdEdit} />
            </Flex>
          </Item>
          <Item onClick={onDeleteOpen} className={'delete-item'}>
            <Flex align='center' justify='space-between' w='full'>
              <Text>Delete Message</Text>
              <Icon as={FaRegTrashAlt} />
            </Flex>
          </Item>
        </Menu>
        <DeleteMessageModal message={message} isOpen={isDeleteOpen} onClose={onDeleteClose} />
        <EditMessageModal message={message} isOpen={isEditOpen} onClose={onEditClose} />
      </>
      }
    </>
  );
}
