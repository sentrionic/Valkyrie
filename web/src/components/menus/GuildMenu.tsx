import React from 'react';
import { Flex, GridItem, Heading, Icon, Menu, MenuButton, MenuDivider, useDisclosure } from '@chakra-ui/react';
import { FiChevronDown, FiX } from 'react-icons/fi';
import { FaUserPlus, FaUserEdit } from 'react-icons/fa';
import { MdAddCircle } from 'react-icons/md';
import { HiLogout } from 'react-icons/hi';
import { RiSettings5Fill } from 'react-icons/ri';
import { StyledMenuList } from './StyledMenuList';
import { StyledMenuItem, StyledRedMenuItem } from './StyledMenuItem';
import { useHistory, useParams } from 'react-router-dom';
import { leaveGuild } from '../../lib/api/handler/guilds';
import { RouterProps } from '../../routes/Routes';
import { userStore } from '../../lib/stores/userStore';
import { useGetCurrentGuild } from '../../lib/utils/hooks/useGetCurrentGuild';
import { GuildSettingsModal } from '../modals/GuildSettingsModal';
import { EditMemberModal } from "../modals/EditMemberModal";
import { Guild, Member } from '../../lib/api/models';
import { gKey, mKey } from '../../lib/utils/querykeys';
import { useQueryClient } from 'react-query';

interface GuildMenuProps {
  channelOpen: () => void;
  inviteOpen: () => void;
}

export const GuildMenu: React.FC<GuildMenuProps> = ({ channelOpen, inviteOpen }) => {

  const { guildId } = useParams<RouterProps>();
  const guild = useGetCurrentGuild(guildId);
  const history = useHistory();
  const cache = useQueryClient();

  const user = userStore(state => state.current);
  const isOwner = guild?.ownerId === user?.id;

  const { isOpen, onOpen, onClose } = useDisclosure();
  const { isOpen: memberOpen, onOpen: memberOnOpen, onClose: memberOnClose } = useDisclosure();

  const handleLeave = async () => {
    const { data } = await leaveGuild(guildId);
    if (data) {
      cache.setQueryData<Guild[]>(gKey, (d) => {
        return d!.filter(g => g.id !== guild?.id);
      });
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
              <StyledMenuItem
                label={'Invite People'}
                icon={FaUserPlus}
                handleClick={inviteOpen}
              />
              {isOwner &&
                <StyledMenuItem
                  label={'Server Settings'}
                  icon={RiSettings5Fill}
                  handleClick={onOpen}
                />
              }
              {isOwner &&
                <StyledMenuItem
                  label={'Create Channel'}
                  icon={MdAddCircle}
                  handleClick={channelOpen}
                />
              }
              <MenuDivider />
              <StyledMenuItem
                label={'Change Appearance'}
                icon={FaUserEdit}
                handleClick={memberOnOpen}
              />
              {!isOwner &&
                <>
                  <MenuDivider />
                  <StyledRedMenuItem
                    label={'Leave Server'}
                    icon={HiLogout}
                    handleClick={handleLeave}
                  />
                </>
              }
            </StyledMenuList>
          </>
        )}
      </Menu>
      {isOpen &&
        <GuildSettingsModal guildId={guildId} isOpen={isOpen} onClose={onClose} />
      }
      {memberOpen &&
        <EditMemberModal guildId={guildId} isOpen={memberOpen} onClose={memberOnClose} />
      }
    </GridItem>
  );
}
