import React from 'react';
import { Button, Flex, GridItem, Icon, LightMode, Text } from '@chakra-ui/react';
import { FiUsers } from 'react-icons/fi';

export const FriendsListHeader: React.FC = () => {
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
        </Flex>
        <LightMode>
          <Button size='sm' colorScheme='blue'>
            Add Friend
          </Button>
        </LightMode>
      </Flex>
    </GridItem>
  );
};
