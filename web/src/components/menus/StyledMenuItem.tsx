import React from "react";
import { Flex, Icon, MenuItem, Text } from "@chakra-ui/react";
import { IconType } from "react-icons";

interface StyledMenuItemProps {
  label: string;
  icon: IconType;
  handleClick: () => void;
}

export const StyledMenuItem: React.FC<StyledMenuItemProps> = ({
  label,
  icon,
  handleClick,
}) => {
  return (
    <MenuItem
      _hover={{ bg: "#7289da", borderRadius: "2px" }}
      onClick={handleClick}
    >
      <Flex align="center" justify="space-between" w="full">
        <Text>{label}</Text>
        <Icon as={icon} />
      </Flex>
    </MenuItem>
  );
};

export const StyledRedMenuItem: React.FC<StyledMenuItemProps> = ({
  label,
  icon,
  handleClick,
}) => {
  return (
    <MenuItem
      _hover={{ bg: "#f04747", color: "#fff", borderRadius: "2px" }}
      onClick={handleClick}
    >
      <Flex align="center" justify="space-between" w="full">
        <Text>{label}</Text>
        <Icon as={icon} />
      </Flex>
    </MenuItem>
  );
};
