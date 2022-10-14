import React from 'react';
import { GridItem, UnorderedList } from '@chakra-ui/react';
import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { MemberListItem } from '../../items/MemberListItem';
import { getGuildMembers } from '../../../lib/api/handler/guilds';
import { mKey } from '../../../lib/utils/querykeys';
import { memberScrollbarCss } from './css/MemberScrollerCSS';
import { useMemberSocket } from '../../../lib/api/ws/useMemberSocket';
import { OnlineLabel } from '../../sections/OnlineLabel';
import { RouterProps } from '../../../lib/models/routerProps';
import { Member } from '../../../lib/models/member';

export const MemberList: React.FC = () => {
  const { guildId } = useParams<keyof RouterProps>() as RouterProps;
  const key = [mKey, guildId];

  const { data } = useQuery(key, () => getGuildMembers(guildId).then((response) => response.data));

  const online: Member[] = [];
  const offline: Member[] = [];

  if (data) {
    data.forEach((m) => {
      if (m.isOnline) {
        online.push(m);
      } else {
        offline.push(m);
      }
    });
  }

  useMemberSocket(guildId);

  return (
    <GridItem
      gridColumn={4}
      gridRow="1 / 4"
      bg="memberList"
      overflowY="hidden"
      _hover={{ overflowY: 'auto' }}
      css={memberScrollbarCss}
    >
      <UnorderedList listStyleType="none" ml="0" id="member-list">
        <OnlineLabel label={`online—${online.length}`} />
        {online.map((m) => (
          <MemberListItem key={`${m.id}`} member={m} />
        ))}
        <OnlineLabel label={`offline—${offline.length}`} />
        {offline.map((m) => (
          <MemberListItem key={`${m.id}`} member={m} />
        ))}
      </UnorderedList>
    </GridItem>
  );
};
