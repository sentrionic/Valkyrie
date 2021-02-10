import React from 'react';
import { Flex, GridItem, Heading, Icon, Menu, MenuButton } from '@chakra-ui/react';
import { FiChevronDown, FiX } from 'react-icons/fi';
import { FaUserPlus } from 'react-icons/fa';
import { MdAddCircle } from 'react-icons/md';
import { HiLogout } from 'react-icons/hi';
import { StyledMenuList } from './StyledMenuList';
import { StyledMenuItem, StyledRedMenuItem } from './StyledMenuItem';
import { useHistory, useParams } from 'react-router-dom';
import { leaveGuild } from '../../lib/api/handler/guilds';
import { RouterProps } from '../../routes/Routes';
import { userStore } from '../../lib/stores/userStore';
import { useGetCurrentGuild } from '../../lib/utils/hooks/useGetCurrentGuild';

interface GuildMenuProps {
  channelOpen: () => void;
  inviteOpen: () => void;
}

export const GuildMenu: React.FC<GuildMenuProps> = ({ channelOpen, inviteOpen }) => {

  const { guildId } = useParams<RouterProps>();
  const guild = useGetCurrentGuild(guildId);
  const history = useHistory();

  const user = userStore(state => state.current);
  const isOwner = guild?.ownerId === user?.id;

  const handleLeave = async () => {
    const { data } = await leaveGuild(guildId);
    if (data) {
      history.replace('/channels/me');
    }
  };

  return (
    <GridItem
      gridColumn={2}
      gridRow={'1'}
      bg='brandGray.light'
      padding='10px'
      zIndex='2'
      boxShadow='md'
    >
      <Menu placement='bottom-end' isLazy>
        {({ isOpen }) => (
          <>
            <Flex justify='space-between' align='center'>
              <Heading fontSize='20px'>{guild?.name}</Heading>
              <MenuButton>
                <Icon as={!isOpen ? FiChevronDown : FiX} />
              </MenuButton>
            </Flex>
            <StyledMenuList>
              {isOwner &&
                <StyledMenuItem
                  label={'Create Channel'}
                  icon={MdAddCircle}
                  handleClick={channelOpen}
                />
              }
              <StyledMenuItem
                label={'Invite People'}
                icon={FaUserPlus}
                handleClick={inviteOpen}
              />
              {!isOwner &&
                <StyledRedMenuItem
                  label={'Leave Server'}
                  icon={HiLogout}
                  handleClick={handleLeave}
                />
              }
            </StyledMenuList>
          </>
        )}
      </Menu>
    </GridItem>
  );
};
