import React from 'react';
import { Avatar, AvatarBadge, Flex, ListItem, Text } from '@chakra-ui/react';
import { Member } from '../../lib/api/models';
import { userStore } from '../../lib/stores/userStore';
import { useContextMenu } from 'react-contexify';
import { MemberContextMenu } from '../menus/MemberContextMenu';

interface MemberListItemProps {
  member: Member;
}

export const MemberListItem: React.FC<MemberListItemProps> = ({ member }) => {

  const current = userStore(state => state.current);
  const self = current?.id !== member.id;

  const { show } = useContextMenu({
    id: member.id,
  });

  return (
    <>
      <ListItem
        p='5px'
        m='0 10px'
        color={'brandGray.accent'}
        onContextMenu={(e) => {
          if (!self) show(e)
        }}
        _hover={{ bg: 'brandGray.light', borderRadius: '5px', cursor: 'pointer', color: '#fff' }}
      >
        <Flex align='center'>
          <Avatar size='sm' src={member.image}>
            <AvatarBadge boxSize='1.25em' bg={member.isOnline ? 'green.500' : 'gray.500'} />
          </Avatar>
          <Text ml='2'>{member.username}</Text>
        </Flex>
      </ListItem>
      {!self &&
        <MemberContextMenu id={member.id} isFriend={member.isFriend} />
      }
    </>
  );
};
