import React from 'react';
import { Box, Flex, Image, Text } from '@chakra-ui/react';
import { Message as MessageResponse } from '../../../lib/api/models';

interface MessageProps {
  message: MessageResponse;
}

export const MessageContent: React.FC<MessageProps> = ({ message: { url, text, filetype, createdAt, updatedAt } }) => {
  if (url && filetype) {
    if (filetype.startsWith('image/')) {
      return (
        <Box boxSize='sm' my={'2'} h={'full'}>
          <Image
            fit={'contain'}
            src={url}
            alt={''}
            borderRadius='md'
          />
        </Box>
      );
    } else if (filetype.startsWith('audio/')) {
      return (
        <Box my={'2'}>
          <audio controls>
            <source src={url} type={filetype} />
          </audio>
        </Box>
      );
    }
  }
  return (
    <Flex alignItems={'center'}>
      <Text>{text}</Text>
      {createdAt !== updatedAt &&
        <Text fontSize={'10px'} ml={'1'} color={'labelGray'}>(edited)</Text>
      }
    </Flex>
  );
};
