import { MenuList } from "@chakra-ui/react";
import React from "react";

export const StyledMenuList: React.FC = ({ children }) => {
  return (
    <MenuList bg="brandGray.darkest" px="2">
      {children}
    </MenuList>
  );
};
