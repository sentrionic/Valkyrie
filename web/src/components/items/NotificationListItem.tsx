import React, { useEffect, useState } from 'react';
import { Flex, Avatar } from '@chakra-ui/react';
import { Link, useLocation } from 'react-router-dom';
import { DMChannel, DMNotification } from '../../lib/api/models';
import { StyledTooltip } from '../sections/StyledTooltip';
import { ActiveGuildPill, HoverGuildPill, NotificationIndicator } from '../common/GuildPills';
import { useQueryClient } from 'react-query';
import { NotificationIcon } from '../common/NotificationIcon';
import { dmKey, nKey } from '../../lib/utils/querykeys';

interface NotificationListItemProps {
  notification: DMNotification;
}

export const NotificationListItem: React.FC<NotificationListItemProps> = ({ notification }) => {
  const location = useLocation();
  const isActive = location.pathname.includes(notification.id);
  const [isHover, setHover] = useState(false);
  const cache = useQueryClient();

  useEffect(() => {
    if (isActive) {
      cache.setQueryData<DMNotification[]>(nKey, (d) => {
        return d!.filter((c) => c.id !== notification.id);
      });
    }
  });

  const handleClick = () => {
    if (window.location.pathname.includes('/channels/me')) {
      const newChannel: DMChannel = {
        id: notification.id,
        user: notification.user,
      };

      cache.setQueryData<DMChannel[]>(dmKey, (d) => {
        const index = d!.findIndex((d) => d.id === notification.id);
        if (index === -1) return [newChannel, ...d!];
        return d!;
      });
    }
  };

  return (
    <Flex mb={'2'} justify={'center'} position={'relative'}>
      {isActive && <ActiveGuildPill />}
      {isHover && <HoverGuildPill />}
      <NotificationIndicator />
      <StyledTooltip label={notification.user.username} position={'right'}>
        <Link to={`/channels/me/${notification.id}`}>
          <Avatar
            src={notification.user.image}
            borderRadius={isActive || isHover ? '35%' : '50%'}
            name={notification.user.username}
            color={'#fff'}
            bg={'brandGray.light'}
            onMouseLeave={() => setHover(false)}
            onMouseEnter={() => setHover(true)}
            onClick={() => handleClick()}
          >
            <NotificationIcon count={notification.count} />
          </Avatar>
        </Link>
      </StyledTooltip>
    </Flex>
  );
};
