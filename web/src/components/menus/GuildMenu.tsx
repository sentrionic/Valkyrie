import React from 'react';
import {
  Flex,
  GridItem,
  Heading,
  Icon,
  Menu,
  MenuButton,
} from '@chakra-ui/react';
import { FiChevronDown, FiX } from 'react-icons/fi';
import { FaUserPlus } from 'react-icons/fa';
import { MdAddCircle } from 'react-icons/md';
import { HiLogout } from 'react-icons/hi';
import { StyledMenuList } from './StyledMenuList';
import { StyledMenuItem, StyledRedMenuItem } from './StyledMenuItem';
import { useHistory, useParams } from 'react-router-dom';
import { useQuery } from 'react-query';
import { Guild } from '../../lib/api/models';
import { leaveGuild } from '../../lib/api/handler/guilds';

interface GuildMenuProps {
  channelOpen: () => void;
  inviteOpen: () => void;
}

interface RouterProps {
  guildId: string;
}

export const GuildMenu: React.FC<GuildMenuProps> = ({
  channelOpen,
  inviteOpen,
}) => {

  const { guildId } = useParams<RouterProps>();
  const { data } = useQuery<Guild[]>('guilds');
  const guild = data?.find(g => g.id === guildId);
  const history = useHistory();

  const handleLeave = async () => {
    const { data } = await leaveGuild(guildId);
    if (data) {
      history.replace('/channels/me');
    }
  }

  return (
    <GridItem
      gridColumn={2}
      gridRow={'1'}
      bg="brandGray.light"
      padding="10px"
      zIndex="2"
      boxShadow="md"
    >
      <Menu placement="bottom-end">
        {({ isOpen }) => (
          <>
            <Flex justify="space-between" align="center">
              <Heading fontSize="20px">{guild?.name}</Heading>
              <MenuButton>
                <Icon as={!isOpen ? FiChevronDown : FiX} />
              </MenuButton>
            </Flex>
            <StyledMenuList>
              <StyledMenuItem
                label={'Create Channel'}
                icon={MdAddCircle}
                handleClick={channelOpen}
              />
              <StyledMenuItem
                label={'Invite People'}
                icon={FaUserPlus}
                handleClick={inviteOpen}
              />
              <StyledRedMenuItem
                label={'Leave Server'}
                icon={HiLogout}
                handleClick={handleLeave}
              />
            </StyledMenuList>
          </>
        )}
      </Menu>
    </GridItem>
  );
};
