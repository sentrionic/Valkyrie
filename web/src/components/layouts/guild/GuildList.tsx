import React from 'react';
import { Box, GridItem, UnorderedList, useDisclosure } from '@chakra-ui/react';
import { useQuery } from 'react-query';
import { AddGuildModal } from '../../modals/AddGuildModal';
import { GuildListItem } from '../../items/GuildListItem';
import { AddGuildIcon } from '../../sections/AddGuildIcon';
import { HomeIcon } from '../../sections/HomeIcon';
import { getUserGuilds } from '../../../lib/api/handler/guilds';
import { gKey } from '../../../lib/utils/querykeys';
import { guildScrollbarCss } from './css/GuildScrollerCSS';

export const GuildList: React.FC = () => {
  const { isOpen, onOpen, onClose } = useDisclosure();

  const { data } = useQuery(gKey, () => {
      return getUserGuilds().then(response => response.data);
    },
    {
      cacheTime: Infinity
    }
  );

  return (
    <GridItem
      gridColumn={1}
      gridRow={'1 / 4'}
      bg='#202225'
      overflowY='auto'
      css={guildScrollbarCss}
    >
      <HomeIcon />
      <UnorderedList listStyleType='none' ml='0'>
        {data?.map(g => <GuildListItem guild={g} key={g.id} />)}
      </UnorderedList>
      <AddGuildIcon onOpen={onOpen} />
      {isOpen &&
        <AddGuildModal isOpen={isOpen} onClose={onClose} />
      }
      <Box h='20px' />
    </GridItem>
  );
};
