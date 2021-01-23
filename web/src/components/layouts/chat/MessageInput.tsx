import {
  GridItem,
  InputGroup,
  Input,
  InputLeftElement,
  Icon,
} from "@chakra-ui/react";
import React from "react";
import { MdAddCircle } from "react-icons/md";

export const MessageInput: React.FC = () => {
  return (
    <GridItem gridColumn={3} gridRow={3} p="5px 20px 20px 20px" bg="#36393f">
      <InputGroup size="md" bg="#40444b" alignItems="center" borderRadius="8px">
        <Input pl="3rem" placeholder="Message #general" border="0" />
        <InputLeftElement _hover={{ cursor: "pointer" }}>
          <Icon as={MdAddCircle} boxSize={"20px"} />
        </InputLeftElement>
      </InputGroup>
    </GridItem>
  );
};
