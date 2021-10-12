import React, { useState } from 'react';
import { Avatar, AvatarBadge, Flex, Icon, ListItem, Text } from '@chakra-ui/react';
import { Link, useHistory, useLocation } from 'react-router-dom';
import { IoMdClose } from 'react-icons/io';
import { useQueryClient } from 'react-query';
import { closeDirectMessage } from '../../lib/api/handler/dm';
import { dmKey } from '../../lib/utils/querykeys';
import { DMChannel } from '../../lib/models/dm';

interface DMListItemProps {
  dm: DMChannel;
}

export const DMListItem: React.FC<DMListItemProps> = ({ dm }) => {
  const currentPath = `/channels/me/${dm.id}`;
  const location = useLocation();
  const isActive = location.pathname === currentPath;
  const [showCloseButton, setShowButton] = useState(false);
  const history = useHistory();
  const cache = useQueryClient();

  const closeDM = async (): Promise<void> => {
    try {
      await closeDirectMessage(dm.id);
      cache.setQueryData<DMChannel[]>(dmKey, (d) => d!.filter((c) => c.id !== dm.id));
      if (isActive) {
        history.replace('/channels/me');
      }
    } catch (err) {}
  };

  return (
    <Link to={`/channels/me/${dm.id}`}>
      <ListItem
        p="2"
        mx="2"
        color={isActive ? '#fff' : 'brandGray.accent'}
        _hover={{
          bg: 'brandGray.light',
          borderRadius: '5px',
          cursor: 'pointer',
          color: '#fff',
        }}
        bg={isActive ? 'brandGray.active' : undefined}
        onMouseLeave={() => setShowButton(false)}
        onMouseEnter={() => setShowButton(true)}
      >
        <Flex align="center" justify="space-between">
          <Flex align="center">
            <Avatar size="sm" src={dm.user.image}>
              <AvatarBadge boxSize="1.25em" bg={dm.user.isOnline ? 'green.500' : 'gray.500'} />
            </Avatar>
            <Text ml="2">{dm.user.username}</Text>
          </Flex>
          {showCloseButton && (
            <Icon
              aria-label="close dm"
              as={IoMdClose}
              onClick={async (e) => {
                e.preventDefault();
                await closeDM();
              }}
            />
          )}
        </Flex>
      </ListItem>
    </Link>
  );
};
