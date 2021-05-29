import React from 'react';
import { Box, Divider, Flex, GridItem, UnorderedList, useDisclosure } from '@chakra-ui/react';
import { useQuery } from 'react-query';
import { AddGuildModal } from '../../modals/AddGuildModal';
import { GuildListItem } from '../../items/GuildListItem';
import { AddGuildIcon } from '../../sections/AddGuildIcon';
import { HomeIcon } from '../../sections/HomeIcon';
import { getUserGuilds } from '../../../lib/api/handler/guilds';
import { gKey, nKey } from '../../../lib/utils/querykeys';
import { guildScrollbarCss } from './css/GuildScrollerCSS';
import { useGuildSocket } from '../../../lib/api/ws/useGuildSocket';
import { DMNotification } from '../../../lib/api/models';
import { NotificationListItem } from '../../items/NotificationListItem';

export const GuildList: React.FC = () => {
  const { isOpen, onOpen, onClose } = useDisclosure();

  const { data } = useQuery(
    gKey,
    () => {
      return getUserGuilds().then((response) => response.data);
    },
    {
      cacheTime: Infinity,
    }
  );

  const { data: dmData } = useQuery<DMNotification[]>(nKey, () => [], {
    cacheTime: Infinity,
  });

  useGuildSocket();

  return (
    <GridItem
      gridColumn={1}
      gridRow={'1 / 4'}
      bg="brandGray.darker"
      overflowY="auto"
      css={guildScrollbarCss}
      zIndex={2}
    >
      <HomeIcon />
      <UnorderedList listStyleType="none" ml="0">
        {dmData?.map((dm) => (
          <NotificationListItem notification={dm} key={dm.id} />
        ))}
      </UnorderedList>
      <Flex direction="column" my="2" align="center">
        <Divider w="40px" />
      </Flex>
      <UnorderedList listStyleType="none" ml="0">
        {data?.map((g) => (
          <GuildListItem guild={g} key={g.id} />
        ))}
      </UnorderedList>
      <AddGuildIcon onOpen={onOpen} />
      {isOpen && <AddGuildModal isOpen={isOpen} onClose={onClose} />}
      <Box h="20px" />
    </GridItem>
  );
};
