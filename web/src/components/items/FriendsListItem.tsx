import {
  Avatar,
  AvatarBadge,
  Flex,
  IconButton,
  ListItem,
  Text,
} from '@chakra-ui/react';
import React from 'react';
import { FaEllipsisV } from 'react-icons/fa';
import { Member } from '../../lib/api/models';
import { getOrCreateDirectMessage } from '../../lib/api/handler/dm';
import { useHistory } from 'react-router-dom';

interface FriendsListItemProp {
  friend: Member
}

export const FriendsListItem: React.FC<FriendsListItemProp> = ({ friend }) => {

  const history = useHistory();

  const getDMChannel = async () => {
    const { data } = await getOrCreateDirectMessage(friend.id);
    if (data) {
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
        cursor: 'pointer',
      }}
      onClick={getDMChannel}
    >
      <Flex align="center" justify="space-between">
        <Flex align="center">
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
          }}
        />
      </Flex>
    </ListItem>
  );
};
