import React from "react";
import { GridItem, Flex, Text, Icon } from "@chakra-ui/react";
import { FaHashtag } from "react-icons/fa";
import { BsPeopleFill } from "react-icons/bs";

export const ChannelHeader: React.FC = () => {
  return (
    <GridItem
      gridColumn={3}
      gridRow={"1"}
      bg="brandGray.light"
      padding="10px"
      zIndex="2"
      boxShadow="md"
    >
      <Flex align="center" justify="space-between">
        <Flex align="center">
          <FaHashtag />
          <Text ml="2" fontWeight="semibold">
            general
          </Text>
        </Flex>
        <Icon
          as={BsPeopleFill}
          fontSize="20px"
          mr="2"
          _hover={{ cursor: "pointer" }}
        />
      </Flex>
    </GridItem>
  );
};
