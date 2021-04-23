import React, { useEffect } from 'react';
import { useQuery } from 'react-query';
import { Flex, UnorderedList, Text } from '@chakra-ui/react';
import { rKey } from '../../../../lib/utils/querykeys';
import { getPendingRequests } from '../../../../lib/api/handler/account';
import { OnlineLabel } from '../../../sections/OnlineLabel';
import { RequestListItem } from '../../../items/RequestListItem';
import { homeStore } from '../../../../lib/stores/homeStore';
import { useRequestSocket } from '../../../../lib/api/ws/useRequestSocket';

export const PendingList: React.FC = () => {
  const { data } = useQuery(rKey, () =>
      getPendingRequests().then(response => response.data),
    {
      staleTime: 0
    }
  );

  useRequestSocket();

  const reset = homeStore(state => state.resetRequest);

  useEffect(() => {
    reset();
  });

  if (data?.length === 0) {
    return (
      <Flex justify={'center'} align={"center"} w={'full'}>
        <Text textColor={"brandGray.accent"}>
          There are no pending friend requests
        </Text>
      </Flex>
    );
  }

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
