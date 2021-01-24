import React from 'react';
import { Flex } from '@chakra-ui/react';
import { AiOutlinePlus } from 'react-icons/ai';

interface AddGuildIconProps {
  onOpen: () => void;
}

export const AddGuildIcon: React.FC<AddGuildIconProps> = ({ onOpen }) => {
  return (
    <Flex
      direction="column"
      m="auto"
      align="center"
      justify="center"
      bg='brandGray.light'
      borderRadius="50%"
      h='48px'
      w='48px'
      _hover={{
        cursor: "pointer",
        borderRadius: "35%",
        bg: '#43b581',
        color: 'white'
      }}
      onClick={onOpen}
    >
      <AiOutlinePlus fontSize="25px" />
    </Flex>
  );
}
