import React, { useEffect, useState } from 'react';
import { Flex, Avatar } from '@chakra-ui/react';
import { Link, useLocation } from 'react-router-dom';
import { Guild } from '../../lib/api/models';
import { StyledTooltip } from '../sections/StyledTooltip';
import { ActiveGuildPill, HoverGuildPill, NotificationIndicator } from '../common/GuildPills';
import { useQueryClient } from 'react-query';
import { gKey } from '../../lib/utils/querykeys';

interface GuildListItemProps {
  guild: Guild;
}

export const GuildListItem: React.FC<GuildListItemProps> = ({ guild }) => {
  const location = useLocation();
  const isActive = location.pathname.includes(guild.id);
  const [isHover, setHover] = useState(false);
  const cache = useQueryClient();

  useEffect(() => {
    if (guild.hasNotification && isActive) {
      cache.setQueryData<Guild[]>(gKey, (d) => {
        const index = d!.findIndex((c) => c.id === guild.id);
        if (index !== -1) {
          d![index] = { ...d![index], hasNotification: false };
        }
        return d!;
      });
    }
  });

  return (
    <Flex mb={'2'} justify={'center'} position={'relative'}>
      {isActive && <ActiveGuildPill />}
      {isHover && <HoverGuildPill />}
      {guild.hasNotification && <NotificationIndicator />}
      <StyledTooltip label={guild.name} position={'right'}>
        <Link to={`/channels/${guild.id}/${guild.default_channel_id}`}>
          {guild.icon ? (
            <Avatar
              src={guild.icon}
              borderRadius={isActive || isHover ? '35%' : '50%'}
              name={guild.name}
              color={'#fff'}
              bg={'brandGray.light'}
              onMouseLeave={() => setHover(false)}
              onMouseEnter={() => setHover(true)}
            />
          ) : (
            <Flex
              justify={'center'}
              align={'center'}
              bg={isActive ? 'highlight.standard' : 'brandGray.light'}
              borderRadius={isActive ? '35%' : '50%'}
              h={'48px'}
              w={'48px'}
              color={isActive ? 'white' : undefined}
              fontSize="20px"
              _hover={{
                borderRadius: '35%',
                bg: 'highlight.standard',
                color: 'white',
              }}
              onMouseLeave={() => setHover(false)}
              onMouseEnter={() => setHover(true)}
            >
              {guild.name[0]}
            </Flex>
          )}
        </Link>
      </StyledTooltip>
    </Flex>
  );
};
