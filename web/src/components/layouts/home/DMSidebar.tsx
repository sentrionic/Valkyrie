import { GridItem, Text, UnorderedList } from '@chakra-ui/react';
import React from 'react';
import { FriendsListButton } from '../../sections/FriendsListButton';
import { DMListItem } from '../../items/DMListItem';
import { AccountBar } from '../AccountBar';

export const DMSidebar: React.FC = () => {
  return (
    <GridItem gridColumn={'2'} gridRow={'1 / 4'} bg='brandGray.dark'>
      <FriendsListButton />
      <Text
        ml='4'
        textTransform='uppercase'
        fontSize='12px'
        fontWeight='semibold'
        color='brandGray.accent'
      >
        DIRECT MESSAGES
      </Text>
      <UnorderedList listStyleType='none' ml='0' mt='4'>
        <DMListItem />
        <DMListItem />
      </UnorderedList>
      <AccountBar />
    </GridItem>
  );
};
