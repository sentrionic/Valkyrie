import React from 'react';
import { Flex } from '@chakra-ui/react';
import { NavBar } from '../sections/NavBar';

export const LandingLayout: React.FC = ({ children }) => {
  return (
    <Flex direction="column" align="center" maxW={{ xl: '1200px' }} m="0 auto">
      <NavBar />
      {children}
    </Flex>
  );
};
