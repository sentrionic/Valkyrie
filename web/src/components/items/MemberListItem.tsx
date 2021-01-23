import React from 'react';
import { Avatar, AvatarBadge, Flex, ListItem, Text } from '@chakra-ui/react';

export const MemberListItem: React.FC = () => {
  return (
    <ListItem
      p="5px"
      m="0 10px"
      _hover={{ bg: "#36393f", borderRadius: "5px", cursor: "pointer" }}
    >
      <Flex align="center">
        <Avatar size="sm">
          <AvatarBadge boxSize="1.25em" bg="green.500" />
        </Avatar>
        <Text ml="2">Username</Text>
      </Flex>
    </ListItem>
  );
}
