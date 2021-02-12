import { GridItem, UnorderedList } from '@chakra-ui/react';
import React from 'react';
import { FriendsListHeader } from './FriendsListHeader';
import { scrollbarCss } from '../../../lib/utils/theme';
import { useQuery } from 'react-query';
import { fKey } from '../../../lib/utils/querykeys';
import { getFriends } from '../../../lib/api/handler/account';
import { FriendsListItem } from '../../items/FriendsListItem';

export const FriendList: React.FC = () => {

  const { data } = useQuery(fKey, () =>
    getFriends().then(response => response.data),
    {
      cacheTime: Infinity
    }
  );

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
          {data?.map((f) =>
            <FriendsListItem key={f.id} friend={f} />
          )}
        </UnorderedList>
      </GridItem>
    </>
  );
};
