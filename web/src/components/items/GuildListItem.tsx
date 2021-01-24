import React from 'react';
import { ListItem } from '@chakra-ui/react';

export const GuildListItem: React.FC = () => {
  return (
    <ListItem
      h='48px'
      w='48px'
      bg='brandGray.light'
      m='auto'
      mb='3'
      fontSize='24px'
      borderRadius='50%'
      alignItems='center'
      justifyContent='center'
      display='flex'
      _hover={{
        cursor: 'pointer',
        borderRadius: '35%',
        bg: 'highlight.standard',
        color: 'white'
      }}
    >
      H
    </ListItem>
  );
};
