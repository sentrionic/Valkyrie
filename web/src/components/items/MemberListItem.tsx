import React from 'react';
import { Avatar, AvatarBadge, Flex, ListItem, Text } from '@chakra-ui/react';
import { useContextMenu } from 'react-contexify';
import { Member } from '../../lib/api/models';
import { useGetCurrentGuild } from '../../lib/utils/hooks/useGetCurrentGuild';
import { userStore } from '../../lib/stores/userStore';
import { useParams } from 'react-router-dom';
import { RouterProps } from '../../routes/Routes';
import { MemberContextMenu } from '../menus/MemberContextMenu';

interface MemberListItemProps {
  member: Member;
}

export const MemberListItem: React.FC<MemberListItemProps> = ({ member }) => {

  const current = userStore(state => state.current);
  const { guildId } = useParams<RouterProps>();
  const guild = useGetCurrentGuild(guildId);
  const isOwner = guild !== undefined && guild.ownerId === current?.id;

  const { show } = useContextMenu({
    id: member.id
  });

  return (
    <>
      <ListItem
        p='2'
        mx='10px'
        color={'brandGray.accent'}
        _hover={{ bg: 'brandGray.light', borderRadius: '5px', cursor: 'pointer', color: '#fff' }}
        onContextMenu={show}
      >
        <Flex align='center'>
          <Avatar size='sm' src={member.image}>
            <AvatarBadge boxSize='1.25em' bg={member.isOnline ? 'green.500' : 'gray.500'} />
          </Avatar>
          <Text ml='2'>{member.username}</Text>
        </Flex>
      </ListItem>
      {member.id !== current?.id &&
        <MemberContextMenu member={member} isOwner={isOwner} id={member.id} />
      }
    </>
  );
}
