import React from 'react';
import { GridItem } from '@chakra-ui/react';
import { FriendsListHeader } from './FriendsListHeader';
import { FriendsList } from './FriendsList';
import { PendingList } from './PendingList';
import { scrollbarCss } from '../../../../lib/utils/theme';
import { friendStore } from '../../../../lib/stores/friendStore';

export const FriendsDashboard: React.FC = () => {
  const isPending = friendStore(state => state.isPending);

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
        {isPending ? <PendingList /> : <FriendsList />}
      </GridItem>
    </>
  );
}
