import React from 'react';
import { Avatar, AvatarBadge, Flex, ListItem, Text } from '@chakra-ui/react';
import { Member } from '../../lib/api/models';

interface MemberListItemProps {
  member: Member;
}

export const MemberListItem: React.FC<MemberListItemProps> = ({ member }) => {
  return (
    <ListItem
      p="5px"
      m="0 10px"
      _hover={{ bg: "#36393f", borderRadius: "5px", cursor: "pointer" }}
    >
      <Flex align="center">
        <Avatar size="sm" src={member.image}>
          <AvatarBadge boxSize="1.25em" bg="green.500" />
        </Avatar>
        <Text ml="2">{member.username}</Text>
      </Flex>
    </ListItem>
  );
}
