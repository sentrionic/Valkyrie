import React from 'react';
import { Flex, Icon, Text } from '@chakra-ui/react';
import { FiUsers } from 'react-icons/fi';

export const FriendsListButton: React.FC = () => {
  return(
    <Flex
      m="2"
      p="3"
      align="center"
      _hover={{ cursor: "pointer", bg: "brandGray.light" }}
    >
      <Icon as={FiUsers} fontSize="20px" />
      <Text fontSize="16px" ml="4" fontWeight="semibold">
        Friends
      </Text>
    </Flex>
  );
}
