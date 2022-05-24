import React, { useState } from 'react';
import { Avatar, Box, Flex, Icon, Text, useDisclosure } from '@chakra-ui/react';
import { Item, Menu, theme, useContextMenu } from 'react-contexify';
import { useParams } from 'react-router-dom';
import { MdEdit } from 'react-icons/md';
import { FaEllipsisH, FaRegTrashAlt } from 'react-icons/fa';
import { FiLink } from 'react-icons/fi';
import { MessageContent } from './MessageContent';
import { userStore } from '../../../lib/stores/userStore';
import { getShortenedTime, getTime } from '../../../lib/utils/dateUtils';
import { DeleteMessageModal } from '../../modals/DeleteMessageModal';
import { EditMessageModal } from '../../modals/EditMessageModal';
import { useGetCurrentGuild } from '../../../lib/utils/hooks/useGetCurrentGuild';
import { MemberContextMenu } from '../../menus/MemberContextMenu';
import { UserPopover } from '../../sections/UserPopover';
import { RouterProps } from '../../../lib/models/routerProps';
import { Message as MessageResponse } from '../../../lib/models/message';
import '../css/ContextMenu.css';

interface MessageProps {
  message: MessageResponse;
  isCompact?: boolean;
}

export const Message: React.FC<MessageProps> = ({ message, isCompact = false }) => {
  const [showSettings, setShowSettings] = useState(false);
  const current = userStore((state) => state.current);
  const isAuthor = current?.id === message.user.id;
  const { guildId } = useParams<keyof RouterProps>() as RouterProps;
  const guild = useGetCurrentGuild(guildId);
  const isOwner = guild !== undefined && guild.ownerId === current?.id;
  const showMenu = isAuthor || isOwner || message.attachment?.url;

  const { isOpen: isDeleteOpen, onOpen: onDeleteOpen, onClose: onDeleteClose } = useDisclosure();
  const { isOpen: isEditOpen, onOpen: onEditOpen, onClose: onEditClose } = useDisclosure();
  const id = `${message.user.id}-${Math.random().toString(36).substr(2, 5)}`;

  const { show } = useContextMenu({
    id: message.id,
  });

  const { show: profileShow } = useContextMenu({ id });

  const openInNewTab = (url: string): void => {
    const newWindow = window.open(url, '_blank', 'noopener,noreferrer');
    if (newWindow) newWindow.opener = null;
  };

  return (
    <>
      <Flex
        alignItems="center"
        mr="1"
        mt={isCompact ? '0' : '3'}
        _hover={{ bg: 'brandGray.hover' }}
        justify="space-between"
        onMouseLeave={() => setShowSettings(false)}
        onMouseEnter={() => setShowSettings(true)}
      >
        <Flex w="full">
          {isCompact ? (
            <>
              <Box ml="3" minW="44px" textAlign="center">
                <Text fontSize="10px" color="brandGray.accent" mt="1" hidden={!showSettings}>
                  {getShortenedTime(message.createdAt)}
                </Text>
              </Box>

              <Box ml="3" w="full" onContextMenu={show}>
                <MessageContent message={message} />
              </Box>
              {showSettings && showMenu ? (
                <Box onClick={show} mr="2" _hover={{ cursor: 'pointer' }} h="5px">
                  <FaEllipsisH />
                </Box>
              ) : (
                <Box mr="6" />
              )}
            </>
          ) : (
            <>
              <UserPopover member={message.user}>
                <Avatar
                  h="40px"
                  w="40px"
                  ml="4"
                  mt="1"
                  src={message.user.image}
                  _hover={{
                    cursor: 'pointer',
                  }}
                  onContextMenu={(e) => {
                    if (!isAuthor) profileShow(e);
                  }}
                />
              </UserPopover>
              <Box ml="3" w="full" onContextMenu={show}>
                <Flex alignItems="center" justify="space-between">
                  <Flex alignItems="center">
                    <Text color={message.user.color ?? undefined}>
                      {message.user.nickname ?? message.user.username}
                    </Text>
                    <Text fontSize="12px" color="brandGray.accent" ml="2">
                      {getTime(message.createdAt)}
                    </Text>
                  </Flex>
                  {showSettings && showMenu && (
                    <Box onClick={show} mr="2" _hover={{ cursor: 'pointer' }}>
                      <FaEllipsisH />
                    </Box>
                  )}
                </Flex>
                <MessageContent message={message} />
              </Box>
            </>
          )}
        </Flex>
      </Flex>
      {showMenu && (
        <>
          <Menu id={message.id} theme={theme.dark}>
            {message.attachment?.filetype ? (
              <Item
                className="menu-item"
                onClick={() => {
                  if (message.attachment?.url) openInNewTab(message.attachment.url);
                }}
              >
                <Flex align="center" justify="space-between" w="full">
                  <Text>Open Link</Text>
                  <Icon as={FiLink} />
                </Flex>
              </Item>
            ) : (
              isAuthor && (
                <Item className="menu-item" onClick={onEditOpen}>
                  <Flex align="center" justify="space-between" w="full">
                    <Text>Edit Message</Text>
                    <Icon as={MdEdit} />
                  </Flex>
                </Item>
              )
            )}
            {(isAuthor || isOwner) && (
              <Item onClick={onDeleteOpen} className="delete-item">
                <Flex align="center" justify="space-between" w="full">
                  <Text>Delete Message</Text>
                  <Icon as={FaRegTrashAlt} />
                </Flex>
              </Item>
            )}
          </Menu>
          {isDeleteOpen && <DeleteMessageModal message={message} isOpen={isDeleteOpen} onClose={onDeleteClose} />}
          {isEditOpen && <EditMessageModal message={message} isOpen={isEditOpen} onClose={onEditClose} />}
        </>
      )}
      {!isAuthor && <MemberContextMenu member={message.user} isOwner={isOwner} id={id} />}
    </>
  );
};
