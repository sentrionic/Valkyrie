import {
  Avatar,
  AvatarBadge,
  Flex,
  IconButton,
  ListItem,
  Text, useDisclosure
} from '@chakra-ui/react';
import React from 'react';
import { FaEllipsisV } from 'react-icons/fa';
import { DMChannel, Member } from '../../lib/api/models';
import { getOrCreateDirectMessage } from '../../lib/api/handler/dm';
import { useHistory } from 'react-router-dom';
import { RemoveFriendModal } from '../modals/RemoveFriendModal';
import { useQueryClient } from 'react-query';
import { dmKey } from '../../lib/utils/querykeys';

interface FriendsListItemProp {
  friend: Member
}

export const FriendsListItem: React.FC<FriendsListItemProp> = ({ friend }) => {

  const history = useHistory();
  const { isOpen, onOpen, onClose } = useDisclosure();
  const cache = useQueryClient();

  const getDMChannel = async () => {
    const { data } = await getOrCreateDirectMessage(friend.id);
    if (data) {
      cache.setQueryData<DMChannel[]>(dmKey, (d) => {
        const index = d!.findIndex(d => d.id);
        if (index === -1) return [data, ...d!];
        return d!;
      });
      history.push(`/channels/me/${data.id}`);
    }
  }

  return (
    <ListItem
      p="3"
      mx="3"
      _hover={{
        bg: 'brandGray.dark',
        borderRadius: '5px',
      }}
    >
      <Flex align="center" justify="space-between">
        <Flex align="center" w={"full"}
          onClick={getDMChannel}
          _hover={{ cursor: 'pointer' }}
        >
          <Avatar size="sm" src={friend.image}>
            <AvatarBadge boxSize="1.25em" bg={ friend.isOnline ? 'green.500' : 'gray.500'} />
          </Avatar>
          <Text ml="2">{friend.username}</Text>
        </Flex>
        <IconButton
          icon={<FaEllipsisV />}
          borderRadius="50%"
          aria-label="remove friend"
          onClick={(e) => {
            e.preventDefault();
            onOpen();
          }}
        />
      </Flex>
      {isOpen &&
        <RemoveFriendModal id={friend.id} isOpen onClose={onClose} />
      }
    </ListItem>
  );
};
