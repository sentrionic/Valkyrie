import React from 'react';
import { Item, Menu, theme } from 'react-contexify';
import { Flex, Text } from '@chakra-ui/react';
import { useHistory } from 'react-router-dom';
import { getOrCreateDirectMessage } from '../../lib/api/handler/dm';
import '../items/css/ContextMenu.css';

interface MemberContextMenuProps {
  id: string;
  isFriend: boolean;
}

export const MemberContextMenu: React.FC<MemberContextMenuProps> = ({ id, isFriend }) => {

  const history = useHistory();

  const getOrCreateDM = async () => {
    const { data } = await getOrCreateDirectMessage(id);
    if (data) {
      history.push(`/channels/me/${data.id}`);
    }
  };

  return (
    <Menu id={id} theme={theme.dark}>
      <Item onClick={() => getOrCreateDM()} className={'menu-item'}>
        <Flex align='center' justify='space-between' w='full'>
          <Text>Message</Text>
        </Flex>
      </Item>
      <Item onClick={() => console.log('Friend')} className={'menu-item'}>
        <Flex align='center' justify='space-between' w='full'>
          <Text>{isFriend ? 'Remove' : 'Add'} Friend</Text>
        </Flex>
      </Item>
    </Menu>
  );
}