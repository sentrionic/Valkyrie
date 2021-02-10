import React from 'react';
import { Text } from '@chakra-ui/react';

interface LabelProps {
  label: string;
}

export const OnlineLabel: React.FC<LabelProps> = ({ label }) => {
  return (
    <Text
      fontSize='12px'
      color={'brandGray.accent'}
      textTransform={'uppercase'}
      fontWeight={'semibold'}
      mx={'3'}
      mt={'4'}
      mb={'1'}
    >
      {label}
    </Text>
  );
};
