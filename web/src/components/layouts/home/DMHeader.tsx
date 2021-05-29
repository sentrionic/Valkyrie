import React from 'react';
import { Box, Flex, GridItem, Icon, Text } from '@chakra-ui/react';
import { FaAt } from 'react-icons/fa';
import { useParams } from 'react-router-dom';
import { RouterProps } from '../../../routes/Routes';
import { useGetCurrentDM } from '../../../lib/utils/hooks/useGetCurrentDM';

export const DMHeader: React.FC = () => {
  const { channelId } = useParams<RouterProps>();
  const channel = useGetCurrentDM(channelId);

  return (
    <GridItem gridColumn={3} gridRow={'1'} bg="brandGray.light" padding="10px" zIndex="2" boxShadow="md">
      <Flex align="center" ml={2}>
        <Icon as={FaAt} fontSize={'20px'} color={'brandGray.accent'} />
        <Text ml="2" fontWeight="semibold">
          {channel?.user.username}
        </Text>
        <Box
          ml={'2'}
          borderRadius={'50%'}
          h={'10px'}
          w={'10px'}
          bg={channel?.user.isOnline ? 'green.500' : 'gray.500'}
        />
      </Flex>
    </GridItem>
  );
};
