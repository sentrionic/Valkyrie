import React from 'react';
import { GridItem, UnorderedList, Text } from '@chakra-ui/react';
import { useQuery } from 'react-query';
import { FriendsListHeader } from './FriendsListHeader';
import { scrollbarCss } from '../../../lib/utils/theme';
import { fKey, rKey } from '../../../lib/utils/querykeys';
import { getFriends, getPendingRequests } from '../../../lib/api/handler/account';
import { FriendsListItem } from '../../items/FriendsListItem';
import { OnlineLabel } from '../../sections/OnlineLabel';
import { friendStore } from '../../../lib/stores/friendStore';
import { RequestListItem } from '../../items/RequestListItem';

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

const FriendsList: React.FC = () => {
  const { data } = useQuery(fKey, () =>
      getFriends().then(response => response.data),
    {
      cacheTime: Infinity
    }
  );

  return (
    <>
      <UnorderedList listStyleType='none' ml='0' w='full' mt='2'>
        <OnlineLabel label={`friends — ${data?.length || 0}`} />
        {data?.map((f) =>
          <FriendsListItem key={f.id} friend={f} />
        )}
      </UnorderedList>
    </>
  );
}

const PendingList: React.FC = () => {
  const { data } = useQuery(rKey, () =>
      getPendingRequests().then(response => response.data),
    {
      staleTime: 0
    }
  );

  return (
    <>
      <UnorderedList listStyleType='none' ml='0' w='full' mt='2'>
        <OnlineLabel label={`Pending — ${data?.length || 0}`} />
        {data?.map((r) =>
          <RequestListItem request={r} key={r.id} />
        )}
      </UnorderedList>
    </>
  );
}
