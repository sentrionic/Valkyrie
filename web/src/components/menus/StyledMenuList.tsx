import { MenuList } from '@chakra-ui/react';
import React from 'react';

interface IProps {
  children: React.ReactNode;
}

export const StyledMenuList: React.FC<IProps> = ({ children }) => (
  <MenuList bg="brandGray.darkest" px="2">
    {children}
  </MenuList>
);
