import React, { useState } from 'react';
import { Divider, Flex, Text, useDisclosure } from '@chakra-ui/react';
import { Item, Menu, theme } from 'react-contexify';
import { useHistory } from 'react-router-dom';
import { getOrCreateDirectMessage } from '../../lib/api/handler/dm';
import { sendFriendRequest } from '../../lib/api/handler/account';
import { RemoveFriendModal } from '../modals/RemoveFriendModal';
import { ModActionModal } from '../modals/ModActionModal';
import { Member } from '../../lib/models/member';

interface MemberContextMenuProps {
  member: Member;
  isOwner: boolean;
  id: string;
}

export const MemberContextMenu: React.FC<MemberContextMenuProps> = ({ member, isOwner, id }) => {
  const history = useHistory();
  const { isOpen, onOpen, onClose } = useDisclosure();
  const { isOpen: modIsOpen, onOpen: modOnOpen, onClose: modOnClose } = useDisclosure();
  const [isBan, setIsBan] = useState(false);

  const getOrCreateDM = async (): Promise<void> => {
    try {
      const { data } = await getOrCreateDirectMessage(member.id);
      if (data) {
        history.push(`/channels/me/${data.id}`);
      }
    } catch (err) {}
  };

  const handleFriendClick = async (): Promise<void> => {
    if (!member.isFriend) {
      try {
        await sendFriendRequest(member.id);
      } catch (err) {}
    } else {
      onOpen();
    }
  };

  return (
    <>
      <Menu id={id} theme={theme.dark}>
        <Item onClick={() => getOrCreateDM()} className="menu-item">
          <Flex align="center" justify="space-between" w="full">
            <Text>Message</Text>
          </Flex>
        </Item>
        <Item onClick={handleFriendClick} className="menu-item">
          <Flex align="center" justify="space-between" w="full">
            <Text>{member.isFriend ? 'Remove' : 'Add'} Friend</Text>
          </Flex>
        </Item>
        {isOwner && (
          <>
            <Flex align="center" justify="center" w="full">
              <Divider my="1" w="90%" />
            </Flex>
            <Item
              onClick={() => {
                setIsBan(false);
                modOnOpen();
              }}
              className="delete-item"
            >
              <Flex align="center" justify="space-between" w="full">
                <Text>Kick {member.username}</Text>
              </Flex>
            </Item>
            <Item
              onClick={() => {
                setIsBan(true);
                modOnOpen();
              }}
              className="delete-item"
            >
              <Flex align="center" justify="space-between" w="full">
                <Text>Ban {member.username}</Text>
              </Flex>
            </Item>
          </>
        )}
      </Menu>
      {isOpen && <RemoveFriendModal member={member} isOpen onClose={onClose} />}
      {modIsOpen && <ModActionModal member={member} isOpen={modIsOpen} isBan={isBan} onClose={modOnClose} />}
    </>
  );
};
