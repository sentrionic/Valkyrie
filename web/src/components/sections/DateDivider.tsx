import { Divider, Flex, Text } from '@chakra-ui/react';
import React from 'react';
import { formatDivider } from '../../lib/utils/dateUtils';

interface DateDividerProps {
  date: string;
}

export const DateDivider: React.FC<DateDividerProps> = ({ date }) => (
  <Flex textAlign="center" align="center" mt="2" mx="4" key={date}>
    <Divider />
    <Text w={['75%', '75%', '75%', '40%', '25%']} fontSize="12px" color="brandGray.accent">
      {formatDivider(date)}
    </Text>
    <Divider />
  </Flex>
);
