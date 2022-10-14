import React from 'react';
import { Flex, GridItem, Heading, Icon, Menu, MenuButton, MenuDivider, useDisclosure } from '@chakra-ui/react';
import { FiChevronDown, FiX } from 'react-icons/fi';
import { FaUserEdit, FaUserPlus } from 'react-icons/fa';
import { MdAddCircle } from 'react-icons/md';
import { HiLogout } from 'react-icons/hi';
import { RiSettings5Fill } from 'react-icons/ri';
import { useNavigate, useParams } from 'react-router-dom';
import { useQueryClient } from '@tanstack/react-query';
import { StyledMenuList } from './StyledMenuList';
import { StyledMenuItem, StyledRedMenuItem } from './StyledMenuItem';
import { leaveGuild } from '../../lib/api/handler/guilds';
import { userStore } from '../../lib/stores/userStore';
import { useGetCurrentGuild } from '../../lib/utils/hooks/useGetCurrentGuild';
import { GuildSettingsModal } from '../modals/GuildSettingsModal';
import { EditMemberModal } from '../modals/EditMemberModal';
import { gKey } from '../../lib/utils/querykeys';
import { RouterProps } from '../../lib/models/routerProps';
import { Guild } from '../../lib/models/guild';

interface GuildMenuProps {
  channelOpen: () => void;
  inviteOpen: () => void;
}

export const GuildMenu: React.FC<GuildMenuProps> = ({ channelOpen, inviteOpen }) => {
  const { guildId } = useParams<keyof RouterProps>() as RouterProps;
  const guild = useGetCurrentGuild(guildId);
  const navigate = useNavigate();
  const cache = useQueryClient();

  const user = userStore((state) => state.current);
  const isOwner = guild?.ownerId === user?.id;

  const { isOpen, onOpen, onClose } = useDisclosure();
  const { isOpen: memberOpen, onOpen: memberOnOpen, onClose: memberOnClose } = useDisclosure();

  const handleLeave = async (): Promise<void> => {
    try {
      const { data } = await leaveGuild(guildId);
      if (data) {
        cache.setQueryData<Guild[]>([gKey], (d) => d?.filter((g) => g.id !== guild?.id) ?? []);
        navigate('/channels/me', { replace: true });
      }
    } catch (err) {}
  };

  return (
    <GridItem gridColumn={2} gridRow="1" bg="brandGray.light" padding="10px" zIndex="2" boxShadow="md">
      <Menu placement="bottom-end" isLazy>
        {({ isOpen: menuIsOpen }) => (
          <>
            <Flex justify="space-between" align="center">
              <Heading fontSize="20px" noOfLines={1}>
                {guild?.name}
              </Heading>
              <MenuButton>
                <Icon as={!menuIsOpen ? FiChevronDown : FiX} />
              </MenuButton>
            </Flex>
            <StyledMenuList>
              <StyledMenuItem label="Invite People" icon={FaUserPlus} handleClick={inviteOpen} />
              {isOwner && <StyledMenuItem label="Server Settings" icon={RiSettings5Fill} handleClick={onOpen} />}
              {isOwner && <StyledMenuItem label="Create Channel" icon={MdAddCircle} handleClick={channelOpen} />}
              <MenuDivider />
              <StyledMenuItem label="Change Appearance" icon={FaUserEdit} handleClick={memberOnOpen} />
              {!isOwner && (
                <>
                  <MenuDivider />
                  <StyledRedMenuItem label="Leave Server" icon={HiLogout} handleClick={handleLeave} />
                </>
              )}
            </StyledMenuList>
          </>
        )}
      </Menu>
      {isOpen && <GuildSettingsModal guildId={guildId} isOpen={isOpen} onClose={onClose} />}
      {memberOpen && <EditMemberModal guildId={guildId} isOpen={memberOpen} onClose={memberOnClose} />}
    </GridItem>
  );
};
