import { GridItem, UnorderedList } from '@chakra-ui/react';
import React from 'react';
import { FriendsListHeader } from './FriendsListHeader';
import { scrollbarCss } from '../../../lib/utils/theme';
import { FriendsListItem } from '../../items/FriendsListItem';

export const FriendList: React.FC = () => {
  return (
    <>
      <FriendsListHeader />
      <GridItem
        gridColumn={3}
        gridRow={'2'}
        bg='brandGray.light'
        mr='5px'
        display='flex'
        overflowY='auto'
        css={scrollbarCss}
      >
        <UnorderedList listStyleType='none' ml='0' w='full' mt='2'>
          {[...Array(15)].map((x, i) =>
            <FriendsListItem key={`${i}`} />
          )}
        </UnorderedList>
      </GridItem>
    </>
  );
};
