import React from 'react';
import { GridItem, Box, Text, UnorderedList } from '@chakra-ui/react';
import { useQuery } from '@tanstack/react-query';
import { AccountBar } from '../AccountBar';
import { FriendsListButton } from '../../sections/FriendsListButton';
import { DMListItem } from '../../items/DMListItem';
import { getUserDMs } from '../../../lib/api/handler/dm';
import { dmKey } from '../../../lib/utils/querykeys';
import { dmScrollerCss } from './css/dmScrollerCSS';
import { useDMSocket } from '../../../lib/api/ws/useDMSocket';
import { DMPlaceholder } from '../../sections/DMPlaceholder';

export const DMSidebar: React.FC = () => {
  const { data } = useQuery([dmKey], () => getUserDMs().then((result) => result.data));

  useDMSocket();

  return (
    <GridItem
      gridColumn="2"
      gridRow="1 / 4"
      bg="brandGray.dark"
      overflowY="hidden"
      _hover={{ overflowY: 'auto' }}
      css={dmScrollerCss}
    >
      <FriendsListButton />
      <Text ml="4" textTransform="uppercase" fontSize="12px" fontWeight="semibold" color="brandGray.accent">
        DIRECT MESSAGES
      </Text>
      <UnorderedList listStyleType="none" ml="0" mt="4" id="dm-list">
        {data?.map((dm) => (
          <DMListItem dm={dm} key={dm.id} />
        ))}
        {data?.length === 0 && (
          <Box>
            <DMPlaceholder />
            <DMPlaceholder />
            <DMPlaceholder />
            <DMPlaceholder />
            <DMPlaceholder />
          </Box>
        )}
      </UnorderedList>
      <AccountBar />
    </GridItem>
  );
};
