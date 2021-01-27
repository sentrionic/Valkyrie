import React from 'react';
import { GridItem, UnorderedList, Text } from '@chakra-ui/react';
import { MemberListItem } from '../../items/MemberListItem';
import { useParams } from 'react-router-dom';
import { useQuery } from 'react-query';
import { getGuildMembers } from '../../../lib/api/handler/guilds';

interface RouterProps {
  guildId: string;
}

export const MemberList: React.FC = () => {

  const { guildId } = useParams<RouterProps>();

  const { data } = useQuery(`members-${guildId}`, () =>
      getGuildMembers(guildId).then(response => response.data),
    {
      refetchOnWindowFocus: false,
    },
  );

  return (
    <GridItem
      gridColumn={4}
      gridRow={'1 / 4'}
      bg="#2f3136"
      overflowY="hidden"
      _hover={{ overflowY: 'auto' }}
      css={{
        '&::-webkit-scrollbar': {
          width: '4px',
        },
        '&::-webkit-scrollbar-track': {
          width: '4px',
        },
        '&::-webkit-scrollbar-thumb': {
          background: '#202225',
          borderRadius: '18px',
        },
      }}
    >
      <UnorderedList listStyleType="none" ml="0">
        <Text fontSize="14" p="5px" m="5px 10px">
          Online
        </Text>
        {data?.map(m => (
          <MemberListItem key={`${m.id}`} member={m} />
        ))}
      </UnorderedList>
    </GridItem>
  );
};
