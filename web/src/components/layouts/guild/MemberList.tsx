import {
  Avatar,
  AvatarBadge,
  Flex,
  GridItem,
  ListItem,
  UnorderedList,
  Text,
} from "@chakra-ui/react";
import React from "react";

export const MemberList: React.FC = () => {
  return (
    <GridItem gridColumn={4} gridRow={"1 / 4"} bg="#2f3136">
      <UnorderedList listStyleType="none" ml="0">
        <Text fontSize="14" p="5px" m="5px 10px">
          Online
        </Text>
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
      </UnorderedList>
    </GridItem>
  );
};
