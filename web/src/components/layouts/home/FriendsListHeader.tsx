import React from 'react';
import { Button, Flex, GridItem, Icon, LightMode, Text, useDisclosure } from '@chakra-ui/react';
import { FiUsers } from 'react-icons/fi';
import { AddFriendModal } from '../../modals/AddFriendModal';
import { friendStore } from '../../../lib/stores/friendStore';

export const FriendsListHeader: React.FC = () => {

  const { isOpen, onOpen, onClose } = useDisclosure();
  const toggle = friendStore(state => state.toggleDisplay);
  const isPending = friendStore(state => state.isPending);

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
        <Flex align='center' ml={2}>
          <Icon as={FiUsers} fontSize='20px' />
          <Text ml='2' fontWeight='semibold'>
            Friends
          </Text>
          <Button
            ml={'4'}
            size={"sm"}
            colorScheme={"gray"}
            onClick={() => {
              if (isPending) toggle();
            }}
            variant={!isPending ? "solid" : "ghost"}
          >
            Friends
          </Button>
          <Button
            ml={'2'}
            size={"sm"}
            colorScheme={"gray"}
            variant={isPending ? "solid" : "ghost"}
            onClick={() => {
              if (!isPending) toggle();
            }}
          >
            Pending
          </Button>
        </Flex>
        <LightMode>
          <Button
            size='sm'
            bg={'#43b581'}
            _hover={{ bg: '#43b581' }}
            _active={{ bg: '#43b581' }}
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
