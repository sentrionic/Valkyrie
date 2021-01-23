import { Box, Flex, Heading, Text } from '@chakra-ui/react';
import React from 'react';

export const StartMessages: React.FC = () => {
  return (
    <Flex
      alignItems='center'
      mb='2'
      justify='center'
    >
      <Box textAlign={"center"}>
        <Heading>Welcome to #general</Heading>
        <Text>This is the start of #general channel</Text>
      </Box>
    </Flex>
  );
};
