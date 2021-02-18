import React from 'react';
import { Button, Flex, GridItem, Icon, LightMode, Text, useDisclosure } from '@chakra-ui/react';
import { FiUsers } from 'react-icons/fi';
import { AddFriendModal } from '../../../modals/AddFriendModal';
import { homeStore } from '../../../../lib/stores/homeStore';
import { PingIcon } from '../../../common/NotificationIcon';

export const FriendsListHeader: React.FC = () => {

  const { isOpen, onOpen, onClose } = useDisclosure();
  const toggle = homeStore(state => state.toggleDisplay);
  const isPending = homeStore(state => state.isPending);
  const requests = homeStore(state => state.requestCount);

  return (
    <GridItem
      gridColumn={3}
      gridRow={'1'}
      bg='brandGray.light'
      padding='10px'
      zIndex='2'
      boxShadow='md'
    >
      <Flex align='center' justify='space-between'>
        <Flex align='center' ml={2} fontSize='14px'>
          <Icon as={FiUsers} fontSize='20px' />
          <Text ml='2' fontWeight='semibold'>
            Friends
          </Text>
          <Button
            fontSize='14px'
            ml={'4'}
            size={"xs"}
            colorScheme={"gray"}
            onClick={() => {
              if (isPending) toggle();
            }}
            variant={!isPending ? "solid" : "ghost"}
          >
            Friends
          </Button>
          <Button
            fontSize='14px'
            size={"xs"}
            ml={'2'}
            colorScheme={"gray"}
            variant={isPending ? "solid" : "ghost"}
            onClick={() => {
              if (!isPending) toggle();
            }}
          >
            Pending
            { requests > 0 && <PingIcon count={requests} /> }
          </Button>
        </Flex>
        <LightMode>
          <Button
            fontSize='14px'
            size={"xs"}
            bg={'brandGreen'}
            _hover={{ bg: 'brandGreen' }}
            _active={{ bg: 'brandGreen' }}
            onClick={onOpen}
          >
            Add Friend
          </Button>
        </LightMode>
      </Flex>
      {isOpen &&
      <AddFriendModal isOpen={isOpen} onClose={onClose} />
      }
    </GridItem>
  );
}
