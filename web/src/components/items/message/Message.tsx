import React, { useState } from 'react';
import { Avatar, Box, Flex, Icon, Text, useDisclosure } from '@chakra-ui/react';
import { Item, Menu, theme, useContextMenu } from 'react-contexify';
import { useHistory } from 'react-router-dom';
import { MdEdit } from 'react-icons/md';
import { FaEllipsisH, FaRegTrashAlt } from 'react-icons/fa';
import { FiLink } from 'react-icons/fi';
import { MessageContent } from './MessageContent';
import { Message as MessageResponse } from '../../../lib/api/models';
import { userStore } from '../../../lib/stores/userStore';
import { getShortenedTime, getTime } from '../../../lib/utils/dateUtils';
import { DeleteMessageModal } from '../../modals/DeleteMessageModal';
import { EditMessageModal } from '../../modals/EditMessageModal';
import { getOrCreateDirectMessage } from '../../../lib/api/handler/dm';
import '../css/ContextMenu.css';

interface MessageProps {
  message: MessageResponse;
  isCompact?: boolean;
}

export const Message: React.FC<MessageProps> = ({ message, isCompact = false }) => {

  const [showSettings, setShowSettings] = useState(false);
  const current = userStore(state => state.current);
  const isAuthor = current?.id === message.user.id;
  const history = useHistory();

  const { isOpen: isDeleteOpen, onOpen: onDeleteOpen, onClose: onDeleteClose } = useDisclosure();
  const { isOpen: isEditOpen, onOpen: onEditOpen, onClose: onEditClose } = useDisclosure();

  const { show } = useContextMenu({
    id: message.id
  });

  const { show: profileShow } = useContextMenu({
    id: message.user.id
  });

  const openInNewTab = (url: string) => {
    const newWindow = window.open(url, '_blank', 'noopener,noreferrer');
    if (newWindow) newWindow.opener = null;
  };

  const getOrCreateDM = async () => {
    const { data } = await getOrCreateDirectMessage(message.user.id);
    if (data) {
      history.push(`/channels/me/${data.id}`)
    }
  }

  return (
    <>
      <Flex
        alignItems='center'
        mr='1'
        mt={isCompact ? '0' : '3'}
        _hover={{ bg: '#32353b' }}
        justify='space-between'
        onContextMenu={(e) => {
          if (isAuthor) show(e);
        }}
        onMouseLeave={() => setShowSettings(false)}
        onMouseEnter={() => setShowSettings(true)}
      >
        <Flex w={'full'}>
          {isCompact ?
            <>
              <Box ml={'3'} minW={'44px'} textAlign={'center'}>
                <Text fontSize={'10px'} color='brandGray.accent' mt={'1'} hidden={!showSettings}>
                  {getShortenedTime(message.createdAt)}
                </Text>
              </Box>

              <Box ml='3' w={'full'}>
                <MessageContent message={message} />
              </Box>
              {(isAuthor && (showSettings)) ?
                <Box onClick={show} mr='2' _hover={{ cursor: 'pointer' }} h={'5px'}>
                  <FaEllipsisH />
                </Box>
                :
                <Box mr={'6'} />
              }
            </>
            :
            <>
              <Avatar
                h='40px' w='40px' ml='4' mt={'1'}
                src={message.user.image}
                onContextMenu={(e) => {
                  if (!isAuthor) profileShow(e);
                }} />
              <Box ml='3' w={'full'}>
                <Flex alignItems='center' justify={'space-between'}>
                  <Flex alignItems={'center'}>
                    <Text>{message.user.username}</Text>
                    <Text fontSize='12px' color='brandGray.accent' ml='2'>
                      {getTime(message.createdAt)}
                    </Text>
                  </Flex>
                  {(isAuthor && (showSettings)) && (
                    <Box onClick={show} mr='2' _hover={{ cursor: 'pointer' }}>
                      <FaEllipsisH />
                    </Box>
                  )}
                </Flex>
                <MessageContent message={message} />
              </Box>
            </>
          }
        </Flex>
      </Flex>
      {isAuthor &&
      <>
        <Menu id={message.id} theme={theme.dark}>
          {message.filetype ?
            <Item className={'menu-item'} onClick={() => {
              if (message.url) openInNewTab(message.url);
            }}>
              <Flex align='center' justify='space-between' w='full'>
                <Text>Open Link</Text>
                <Icon as={FiLink} />
              </Flex>
            </Item> :
            <Item className={'menu-item'} onClick={onEditOpen}>
              <Flex align='center' justify='space-between' w='full'>
                <Text>Edit Message</Text>
                <Icon as={MdEdit} />
              </Flex>
            </Item>
          }
          <Item onClick={onDeleteOpen} className={'delete-item'}>
            <Flex align='center' justify='space-between' w='full'>
              <Text>Delete Message</Text>
              <Icon as={FaRegTrashAlt} />
            </Flex>
          </Item>
        </Menu>
        {isDeleteOpen && <DeleteMessageModal message={message} isOpen={isDeleteOpen} onClose={onDeleteClose} />}
        {isEditOpen && <EditMessageModal message={message} isOpen={isEditOpen} onClose={onEditClose} />}
      </>
      }
      {!isAuthor &&
        <Menu id={message.user.id} theme={theme.dark}>
          <Item onClick={() => getOrCreateDM()} className={'menu-item'}>
            <Flex align='center' justify='space-between' w='full'>
              <Text>Message</Text>
            </Flex>
          </Item>
        </Menu>
      }
    </>
  );
};
