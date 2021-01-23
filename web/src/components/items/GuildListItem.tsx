import React from 'react';
import { ListItem } from '@chakra-ui/react';

export const GuildListItem: React.FC = () => {
  return (
    <ListItem
      h='50px'
      w='50px'
      bg='#36393f'
      m='auto'
      mb='10px'
      fontSize='24px'
      borderRadius='50%'
      alignItems='center'
      justifyContent='center'
      display='flex'
      _hover={{
        borderStyle: 'solid',
        borderWidth: 'thick',
        borderColor: '#707070',
        cursor: 'pointer',
        borderRadius: '25%',
      }}
    >
      H
    </ListItem>
  );
};
