import React from 'react';
import { Avatar, AvatarBadge, Flex, ListItem, Text } from '@chakra-ui/react';
import { Member } from '../../lib/api/models';

interface MemberListItemProps {
  member: Member;
}

export const MemberListItem: React.FC<MemberListItemProps> = ({ member }) => {

  return (
    <>
      <ListItem
        p='2'
        mx='10px'
        color={'brandGray.accent'}
        _hover={{ bg: 'brandGray.light', borderRadius: '5px', cursor: 'pointer', color: '#fff' }}
      >
        <Flex align='center'>
          <Avatar size='sm' src={member.image}>
            <AvatarBadge boxSize='1.25em' bg={member.isOnline ? 'green.500' : 'gray.500'} />
          </Avatar>
          <Text ml='2'>{member.username}</Text>
        </Flex>
      </ListItem>
    </>
  );
};
