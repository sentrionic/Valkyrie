import React, { useEffect } from 'react';
import { GridItem, Text, UnorderedList } from '@chakra-ui/react';
import { MemberListItem } from '../../items/MemberListItem';
import { useParams } from 'react-router-dom';
import { useQuery, useQueryClient } from 'react-query';
import { getGuildMembers } from '../../../lib/api/handler/guilds';
import { Member } from '../../../lib/api/models';
import socketIOClient from 'socket.io-client';

interface RouterProps {
  guildId: string;
}

export const MemberList: React.FC = () => {

  const { guildId } = useParams<RouterProps>();
  const key = `members-${guildId}`;
  const cache = useQueryClient();

  const { data } = useQuery(key, () =>
      getGuildMembers(guildId).then(response => response.data),
    {
      refetchOnWindowFocus: false
    }
  );

  console.log(data);

  const online: Member[] = [];
  const offline: Member[] = [];

  if (data) data.forEach(m => {
    if (m.isOnline) online.push(m);
    else offline.push(m);
  });

  useEffect((): any => {
    const socket = socketIOClient(process.env.REACT_APP_API_WS!);
    socket.emit('joinGuild', guildId);
    socket.on('add_member', (newMember: Member) => {
      console.log(newMember);
      cache.setQueryData<Member[]>(key, (_) => {
        return [...data!, newMember].sort((a, b) => a.username.localeCompare(b.username));
      });
    });

    socket.on('remove_member', (memberId: string) => {
      cache.setQueryData<Member[]>(key, (_) => {
        return [...data!.filter(m => m.id !== memberId)];
      });
    });

    socket.on('toggle_online', (memberId: string) => {
      cache.setQueryData<Member[]>(key, (_) => {
        const index = data!.findIndex(m => m.id === memberId);
        data![index].isOnline = true;
        return data!;
      });
    });

    socket.on('toggle_offline', (memberId: string) => {
      cache.setQueryData<Member[]>(key, (_) => {
        const index = data!.findIndex(m => m.id === memberId);
        data![index].isOnline = false;
        return data!;
      });
    });

    return () => {
      socket.emit('leaveRoom', guildId);
      socket.disconnect();
    };
  }, [data, key, cache, guildId]);

  return (
    <GridItem
      gridColumn={4}
      gridRow={'1 / 4'}
      bg='#2f3136'
      overflowY='hidden'
      _hover={{ overflowY: 'auto' }}
      css={{
        '&::-webkit-scrollbar': {
          width: '4px'
        },
        '&::-webkit-scrollbar-track': {
          width: '4px'
        },
        '&::-webkit-scrollbar-thumb': {
          background: '#202225',
          borderRadius: '18px'
        }
      }}
    >
      <UnorderedList listStyleType='none' ml='0'>
        <OnlineStatusLabel label={`online-${online.length}`} />
        {online.map(m => <MemberListItem key={`${m.id}`} member={m} />)}
        <OnlineStatusLabel label={`offline-${offline.length}`} />
        {offline.map(m => <MemberListItem key={`${m.id}`} member={m} />)}
      </UnorderedList>
    </GridItem>
  );
};

interface LabelProps {
  label: string;
}

const OnlineStatusLabel: React.FC<LabelProps> = ({ label }) => {
  return (
    <Text
      fontSize='12px'
      color={'brandGray.accent'}
      textTransform={'uppercase'}
      fontWeight={'semibold'}
      mx={'3'}
      mt={'4'}
      mb={'1'}
    >
      {label}
    </Text>
  );
};
