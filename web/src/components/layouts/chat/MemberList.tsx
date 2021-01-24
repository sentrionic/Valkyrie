import React from 'react';
import { GridItem, UnorderedList, Text } from '@chakra-ui/react';
import { MemberListItem } from '../../items/MemberListItem';

export const MemberList: React.FC = () => {
  return (
    <GridItem
      gridColumn={4}
      gridRow={'1 / 4'}
      bg="#2f3136"
      overflowY="hidden"
      _hover={{ overflowY: 'auto' }}
      css={{
        '&::-webkit-scrollbar': {
          width: '4px',
        },
        '&::-webkit-scrollbar-track': {
          width: '4px',
        },
        '&::-webkit-scrollbar-thumb': {
          background: '#202225',
          borderRadius: '18px',
        },
      }}
    >
      <UnorderedList listStyleType="none" ml="0">
        <Text fontSize="14" p="5px" m="5px 10px">
          Online
        </Text>
        <MemberListItem />
        {[...Array(15)].map((x, i) => (
          <MemberListItem key={`${i}`} />
        ))}
      </UnorderedList>
    </GridItem>
  );
};
