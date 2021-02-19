import React from 'react';
import { Flex, Text, useDisclosure } from '@chakra-ui/react';
import { Item, Menu, theme } from 'react-contexify';
import { useHistory } from 'react-router-dom';
import { getOrCreateDirectMessage } from '../../lib/api/handler/dm';
import { sendFriendRequest } from '../../lib/api/handler/account';
import { RemoveFriendModal } from '../modals/RemoveFriendModal';
import { Member } from '../../lib/api/models';

interface MemberContextMenuProps {
  member: Member;
  isOwner: boolean;
}

export const MemberContextMenu: React.FC<MemberContextMenuProps> = ({ member, isOwner }) => {

  const history = useHistory();
  const { isOpen, onOpen, onClose } = useDisclosure();
  const { isOpen: modIsOpen, onOpen: modOnOpen, onClose: modOnClose } = useDisclosure();

  const getOrCreateDM = async () => {
    const { data } = await getOrCreateDirectMessage(member.id);
    if (data) {
      history.push(`/channels/me/${data.id}`);
    }
  };

  const handleFriendClick = async () => {
    if (!member.isFriend) {
      await sendFriendRequest(member.id);
    } else {
      onOpen();
    }
  }

  return (
    <>
      <Menu id={member.id} theme={theme.dark}>
        <Item onClick={() => getOrCreateDM()} className={'menu-item'}>
          <Flex align='center' justify='space-between' w='full'>
            <Text>Message</Text>
          </Flex>
        </Item>
        <Item onClick={handleFriendClick} className={'menu-item'}>
          <Flex align='center' justify='space-between' w='full'>
            <Text>{member.isFriend ? 'Remove' : 'Add'} Friend</Text>
          </Flex>
        </Item>
        {isOwner &&
          <>
            <Item onClick={handleFriendClick} className={'delete-item'}>
              <Flex align='center' justify='space-between' w='full'>
                <Text>Kick {member.username}</Text>
              </Flex>
            </Item>
            <Item onClick={handleFriendClick} className={'delete-item'}>
              <Flex align='center' justify='space-between' w='full'>
                <Text>Ban {member.username}</Text>
              </Flex>
            </Item>
          </>
        }
      </Menu>
      {isOpen &&
        <RemoveFriendModal id={member.id} isOpen onClose={onClose} />
      }
    </>
  );
}
