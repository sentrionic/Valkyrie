import React from 'react';
import { Flex, Text, useDisclosure } from '@chakra-ui/react';
import { Item, Menu, theme } from 'react-contexify';
import { useHistory } from 'react-router-dom';
import { getOrCreateDirectMessage } from '../../lib/api/handler/dm';
import { sendFriendRequest } from '../../lib/api/handler/account';
import { RemoveFriendModal } from '../modals/RemoveFriendModal';

interface MemberContextMenuProps {
  id: string;
  isFriend: boolean;
}

export const MemberContextMenu: React.FC<MemberContextMenuProps> = ({ id, isFriend }) => {

  const history = useHistory();
  const { isOpen, onOpen, onClose } = useDisclosure();

  const getOrCreateDM = async () => {
    const { data } = await getOrCreateDirectMessage(id);
    if (data) {
      history.push(`/channels/me/${data.id}`);
    }
  };

  const handleFriendClick = async () => {
    if (!isFriend) {
      await sendFriendRequest(id);
    } else {
      onOpen();
    }
  }

  return (
    <>
      <Menu id={id} theme={theme.dark}>
        <Item onClick={() => getOrCreateDM()} className={'menu-item'}>
          <Flex align='center' justify='space-between' w='full'>
            <Text>Message</Text>
          </Flex>
        </Item>
        <Item onClick={handleFriendClick} className={'menu-item'}>
          <Flex align='center' justify='space-between' w='full'>
            <Text>{isFriend ? 'Remove' : 'Add'} Friend</Text>
          </Flex>
        </Item>
      </Menu>
      {isOpen &&
        <RemoveFriendModal id={id} isOpen onClose={onClose} />
      }
    </>
  );
}
