import { Flex, Text } from '@chakra-ui/react';
import React from 'react';

interface NotificationIconProps {
  count: number
}

export const NotificationIcon: React.FC<NotificationIconProps> = ({ count }) =>
  <Flex
    borderRadius={'50%'}
    bg={'menuRed'}
    position={'absolute'}
    bottom={0} right={0}
    transform={'translate(25%, 25%)'}
    border={'0.3em solid'}
    borderColor={'brandBorder'}
    w={'1.4em'}
    h={'1.4em'}
    justify={'center'}
    align={'center'}
  >
    <Text fontSize={'12px'} fontWeight={'bold'} color={'white'}>{count}</Text>
  </Flex>;

export const PingIcon: React.FC<NotificationIconProps> = ({ count }) =>
  <Flex
    borderRadius={'50%'}
    bg={'menuRed'}
    w={'1.2em'}
    h={'1.2em'}
    justify={'center'}
    align={'center'}
    ml={2}
  >
    <Text fontSize={'11px'} fontWeight={'bold'} color={'white'}>{count}</Text>
  </Flex>;
