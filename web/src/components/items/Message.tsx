import { Avatar, Box, Flex, Menu, MenuButton, Text } from "@chakra-ui/react";
import React, { useState } from "react";
import { FaEllipsisV, FaRegTrashAlt } from "react-icons/fa";
import { MdEdit } from "react-icons/md";
import { StyledMenuItem, StyledRedMenuItem } from "../menus/StyledMenuItem";
import { StyledMenuList } from "../menus/StyledMenuList";

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
          <StyledMenuList>
            <StyledMenuItem
              label={"Edit Message"}
              icon={MdEdit}
              handleClick={() => console.log("Edit")}
            />
            <StyledRedMenuItem
              label={"Delete Message"}
              icon={FaRegTrashAlt}
              handleClick={() => console.log("Delete")}
            />
          </StyledMenuList>
        </>
      )}
    </Menu>
  );
};
