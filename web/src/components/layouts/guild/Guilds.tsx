import {
  Avatar,
  Text,
  Divider,
  Flex,
  GridItem,
  ListItem,
  UnorderedList,
  useDisclosure,
} from "@chakra-ui/react";
import React from "react";
import { AddGuildModal } from "../../modals/AddGuildModal";
import { AiOutlinePlus } from "react-icons/ai";

export const Guilds: React.FC = () => {
  const { isOpen, onOpen, onClose } = useDisclosure();

  return (
    <GridItem gridColumn={1} gridRow={"1 / 4"} bg="#202225" overflowY="hidden">
      <Flex direction="column" my="2" align="center">
        <Avatar
          src={`${process.env.PUBLIC_URL}/icon.png`}
          size="md"
          _hover={{ cursor: "pointer" }}
        />
        <Divider mt="2" w="40px" />
      </Flex>
      <UnorderedList listStyleType="none" ml="0">
        <ListItem
          h="50px"
          w="50px"
          bg="#36393f"
          m="auto"
          mb="10px"
          fontSize="24px"
          borderRadius="50%"
          alignItems="center"
          justifyContent="center"
          display="flex"
          _hover={{
            borderStyle: "solid",
            borderWidth: "thick",
            borderColor: "#707070",
            cursor: "pointer",
            borderRadius: "25%",
          }}
        >
          H
        </ListItem>
      </UnorderedList>
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
      <AddGuildModal isOpen={isOpen} onClose={onClose} />
    </GridItem>
  );
};
