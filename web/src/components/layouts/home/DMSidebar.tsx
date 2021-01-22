import {
  GridItem,
  Flex,
  Heading,
  Icon,
  UnorderedList,
  ListItem,
  Avatar,
  IconButton,
  Text,
  AvatarBadge,
} from "@chakra-ui/react";
import React from "react";
import { FaAt } from "react-icons/fa";
import { FiUsers } from "react-icons/fi";
import { RiSettings5Fill } from "react-icons/ri";

export const DMSidebar: React.FC = () => {
  return (
    <GridItem gridColumn={2} gridRow={"1 / 4"} bg="brandGray.dark">
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
      <Text
        ml="4"
        textTransform="uppercase"
        fontSize="12px"
        fontWeight="semibold"
        color="brandGray.accent"
      >
        DIRECT MESSAGES
      </Text>
      <UnorderedList listStyleType="none" ml="0" mt="4">
        <ListItem
          p="5px"
          m="0 10px"
          _hover={{ bg: "#36393f", borderRadius: "5px", cursor: "pointer" }}
        >
          <Flex align="center">
            <Avatar size="sm">
              <AvatarBadge boxSize="1.25em" bg="green.500" />
            </Avatar>
            <Text ml="2">sentrionic</Text>
          </Flex>
        </ListItem>
      </UnorderedList>
      <Flex
        p="10px"
        pos="absolute"
        bottom="0"
        w="240px"
        bg="#292b2f"
        align="center"
        justify="space-between"
      >
        <Flex align="center">
          <Avatar size="sm" />
          <Text ml="2">Username</Text>
        </Flex>
        <IconButton
          icon={<RiSettings5Fill />}
          aria-label="settings"
          size="sm"
          fontSize="20px"
          variant="ghost"
        />
      </Flex>
    </GridItem>
  );
};
