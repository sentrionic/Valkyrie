import { Avatar, Box, Flex, IconButton, ListItem, Text } from '@chakra-ui/react';
import React from 'react';
import { BiCheck } from 'react-icons/bi';
import { AiOutlineClose } from 'react-icons/ai';
import { useQueryClient } from '@tanstack/react-query';
import { StyledTooltip } from '../sections/StyledTooltip';
import { acceptFriendRequest, declineFriendRequest } from '../../lib/api/handler/account';
import { fKey, rKey } from '../../lib/utils/querykeys';
import { FriendRequest, RequestType } from '../../lib/models/friend';

interface RequestListItemProps {
  request: FriendRequest;
}

export const RequestListItem: React.FC<RequestListItemProps> = ({ request }) => {
  const cache = useQueryClient();

  const acceptRequest = async (): Promise<void> => {
    try {
      const { data } = await acceptFriendRequest(request.id);
      if (data) {
        cache.setQueryData<FriendRequest[]>([rKey], (d) => d?.filter((r) => r.id !== request.id) ?? []);
        await cache.invalidateQueries([fKey]);
      }
    } catch (err) {}
  };

  const declineRequest = async (): Promise<void> => {
    try {
      const { data } = await declineFriendRequest(request.id);
      if (data) {
        cache.setQueryData<FriendRequest[]>([rKey], (d) => d?.filter((r) => r.id !== request.id) ?? []);
      }
    } catch (err) {}
  };

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
        <Flex align="center">
          <Avatar size="sm" src={request.image} />
          <Box ml="2">
            <Text>{request.username}</Text>
            <Text fontSize="12px">
              {request.type === RequestType.INCOMING ? 'Incoming Friend Request' : 'Outgoing Friend Request'}
            </Text>
          </Box>
        </Flex>
        <Flex align="center">
          {request.type === 1 && (
            <StyledTooltip label="Accept" position="top">
              <IconButton
                icon={<BiCheck />}
                borderRadius="50%"
                aria-label="accept request"
                fontSize="28px"
                onClick={acceptRequest}
                mr="2"
              />
            </StyledTooltip>
          )}
          <StyledTooltip label="Decline" position="top">
            <IconButton
              icon={<AiOutlineClose />}
              borderRadius="50%"
              aria-label="decline request"
              fontSize="20px"
              onClick={declineRequest}
            />
          </StyledTooltip>
        </Flex>
      </Flex>
    </ListItem>
  );
};
