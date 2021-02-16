import React, { useState } from 'react';
import { Flex, Avatar } from '@chakra-ui/react';
import { Link, useLocation } from 'react-router-dom';
import { Guild } from '../../lib/api/models';
import { StyledTooltip } from '../sections/StyledTooltip';
import { ActiveGuildPill, HoverGuildPill } from '../common/GuildPills';

interface GuildListItemProps {
  guild: Guild
}

export const GuildListItem: React.FC<GuildListItemProps> = ({ guild }) => {

  const location = useLocation();
  const isActive = location.pathname.includes(guild.id);
  const [isHover, setHover] = useState(false);

  return (
    <Flex mb={'2'} justify={'center'}>
      { isActive && <ActiveGuildPill />}
      { isHover && <HoverGuildPill /> }
      <StyledTooltip label={guild.name} position={'right'}>
        <Link to={`/channels/${guild.id}/${guild.default_channel_id}`}>
          {guild.icon ?
            <Avatar
              src={guild.icon}
              borderRadius={isActive ? '35%' : '50%'}
              _hover={{ borderRadius: "35%"}}
              name={guild.name}
              color={'#fff'}
              bg={'brandGray.light'}
              onMouseLeave={() => setHover(false)}
              onMouseEnter={() => setHover(true)}
            />
            :
            <Flex
              justify={'center'}
              align={'center'}
              bg={isActive ? 'highlight.standard' : 'brandGray.light'}
              borderRadius={isActive ? '35%' : '50%'}
              h={'48px'}
              w={'48px'}
              color={isActive ? 'white' : undefined}
              fontSize='20px'
              _hover={{
                borderRadius: '35%',
                bg: 'highlight.standard',
                color: 'white'
              }}
              onMouseLeave={() => setHover(false)}
              onMouseEnter={() => setHover(true)}
            >
              {guild.name[0]}
            </Flex>
          }
        </Link>
      </StyledTooltip>
    </Flex>
  );
}
