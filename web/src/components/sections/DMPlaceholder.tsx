import React from 'react';
import { Box, Flex, Stack } from '@chakra-ui/react';

export const DMPlaceholder: React.FC = () =>
  <Flex align={"center"} m='3'>
    <Box w={'32px'} h={'32px'} borderRadius={'50%'} bg={'#36393f'} />
    <Box ml={2} height="20px" w={'144px'} bg={'#36393f'} borderRadius={'10px'} />
  </Flex>
