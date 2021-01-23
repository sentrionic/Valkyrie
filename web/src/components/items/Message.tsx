import {
  Avatar,
  Box,
  Flex,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
  Text,
} from "@chakra-ui/react";
import React, { useState } from "react";
import { FaEllipsisV, FaRegTrashAlt } from "react-icons/fa";
import { MdEdit } from "react-icons/md";

export const Message: React.FC = () => {
  const [showSettings, setShowSettings] = useState(false);

  return (
    <Menu>
      {({ isOpen }) => (
        <>
          <Flex
            alignItems="center"
            my="2"
            mr="1"
            _hover={{ bg: "brandGray.dark" }}
            justify="space-between"
            onMouseLeave={() => setShowSettings(false)}
            onMouseEnter={() => setShowSettings(true)}
          >
            <Flex alignItems="center">
              <Avatar h="40px" w="40px" ml="4" />
              <Box ml="3">
                <Flex alignItems="center">
                  <Text>sentrionic</Text>
                  <Text fontSize="12px" color="brandGray.accent" ml="3">
                    Today at 7:40 PM
                  </Text>
                </Flex>
                <Text>Hello World</Text>
              </Box>
            </Flex>
            {(showSettings || isOpen) && (
              <MenuButton as={Box} mr="2">
                <FaEllipsisV />
              </MenuButton>
            )}
          </Flex>
          <MenuList bg="#18191c">
            <MenuItem _hover={{ bg: "#7289da" }}>
              <Flex align="center" justify="space-between" w="full">
                <Text>Edit Message</Text>
                <MdEdit />
              </Flex>
            </MenuItem>
            <MenuItem _hover={{ bg: "#f04747", color: "#fff" }}>
              <Flex align="center" justify="space-between" w="full">
                <Text>Delete Message</Text>
                <FaRegTrashAlt />
              </Flex>
            </MenuItem>
          </MenuList>
        </>
      )}
    </Menu>
  );
};
