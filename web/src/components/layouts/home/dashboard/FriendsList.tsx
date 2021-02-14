import React from 'react';
import { useQuery } from 'react-query';
import { UnorderedList } from '@chakra-ui/react';
import { fKey } from '../../../../lib/utils/querykeys';
import { getFriends } from '../../../../lib/api/handler/account';
import { OnlineLabel } from '../../../sections/OnlineLabel';
import { FriendsListItem } from '../../../items/FriendsListItem';
import { useFriendSocket } from '../../../../lib/api/ws/useFriendSocket';

export const FriendsList: React.FC = () => {
  const { data } = useQuery(fKey, () =>
      getFriends().then(response => response.data),
    {
      cacheTime: Infinity
    }
  );

  useFriendSocket();

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
