import {
  Avatar,
  AvatarBadge,
  Button,
  Flex,
  GridItem,
  Icon,
  IconButton,
  LightMode,
  ListItem,
  Text,
  UnorderedList,
} from "@chakra-ui/react";
import React from "react";
import { FaEllipsisV } from "react-icons/fa";
import { FiUsers } from "react-icons/fi";

export const FriendList: React.FC = () => {
  return (
    <>
      <GridItem
        gridColumn={3}
        gridRow={"1"}
        bg="brandGray.light"
        padding="10px"
        zIndex="2"
        boxShadow="md"
      >
        <Flex align="center" justify="space-between">
          <Flex align="center">
            <Icon as={FiUsers} fontSize="20px" />
            <Text ml="2" fontWeight="semibold">
              Friends
            </Text>
          </Flex>
          <LightMode>
            <Button size="sm" colorScheme="blue">
              Add Friend
            </Button>
          </LightMode>
        </Flex>
      </GridItem>
      <GridItem
        gridColumn={3}
        gridRow={"2"}
        bg="brandGray.light"
        mr="5px"
        display="flex"
        overflowY="auto"
        css={{
          "&::-webkit-scrollbar": {
            width: "8px",
          },
          "&::-webkit-scrollbar-track": {
            background: "#2f3136",
            width: "10px",
          },
          "&::-webkit-scrollbar-thumb": {
            background: "#202225",
            borderRadius: "18px",
          },
        }}
      >
        <UnorderedList listStyleType="none" ml="0" w="full" mt="2">
          <ListItem
            p="3"
            m="0 10px"
            _hover={{
              bg: "brandGray.dark",
              borderRadius: "5px",
              cursor: "pointer",
            }}
          >
            <Flex align="center" justify="space-between">
              <Flex align="center">
                <Avatar size="sm">
                  <AvatarBadge boxSize="1.25em" bg="green.500" />
                </Avatar>
                <Text ml="2">Username</Text>
              </Flex>
              <IconButton
                icon={<FaEllipsisV />}
                borderRadius="50%"
                aria-label="remove friend"
              />
            </Flex>
          </ListItem>
        </UnorderedList>
      </GridItem>
    </>
  );
};
