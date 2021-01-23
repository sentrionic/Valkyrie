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
      bg="#36393f"
      borderRadius="50%"
      h="50px"
      w="50px"
      _hover={{
        borderStyle: "solid",
        borderWidth: "thick",
        borderColor: "#707070",
        cursor: "pointer",
        borderRadius: "25%",
      }}
      onClick={onOpen}
    >
      <AiOutlinePlus fontSize="25px" />
    </Flex>
  );
}
