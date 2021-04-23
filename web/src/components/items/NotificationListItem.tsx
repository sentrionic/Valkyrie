import React, { useEffect, useState } from 'react';
import { Flex, Avatar } from '@chakra-ui/react';
import { Link, useLocation } from 'react-router-dom';
import { DMNotification } from '../../lib/api/models';
import { StyledTooltip } from '../sections/StyledTooltip';
import { ActiveGuildPill, HoverGuildPill, NotificationIndicator } from '../common/GuildPills';
import { useQueryClient } from 'react-query';
import { NotificationIcon } from '../common/NotificationIcon';
import { nKey } from '../../lib/utils/querykeys';

interface NotificationListItemProps {
  notification: DMNotification
}

export const NotificationListItem: React.FC<NotificationListItemProps> = ({ notification }) => {

  const location = useLocation();
  const isActive = location.pathname.includes(notification.id);
  const [isHover, setHover] = useState(false);
  const cache = useQueryClient();

  useEffect(() => {
    if (isActive) {
      cache.setQueryData<DMNotification[]>(nKey, (d) => {
        return d!.filter(c => c.id !== notification.id);
      });
    }
  });

  return (
    <Flex mb={'2'} justify={'center'} position={'relative'}>
      { isActive && <ActiveGuildPill />}
      { isHover && <HoverGuildPill /> }
      <NotificationIndicator />
      <StyledTooltip label={notification.user.username} position={'right'}>
        <Link to={`/channels/me/${notification.id}`}>
          <Avatar
            src={notification.user.image}
            borderRadius={(isActive || isHover) ? '35%' : '50%'}
            name={notification.user.username}
            color={'#fff'}
            bg={'brandGray.light'}
            onMouseLeave={() => setHover(false)}
            onMouseEnter={() => setHover(true)}
          >
            <NotificationIcon count={notification.count} />
          </Avatar>
        </Link>
      </StyledTooltip>
    </Flex>
  );
}
