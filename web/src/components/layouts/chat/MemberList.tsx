import React from "react";
import {
  GridItem,
  UnorderedList,
  Text,
} from "@chakra-ui/react";
import { MemberListItem } from '../../items/MemberListItem';

export const MemberList: React.FC = () => {
  return (
    <GridItem gridColumn={4} gridRow={"1 / 4"} bg="#2f3136">
      <UnorderedList listStyleType="none" ml="0">
        <Text fontSize="14" p="5px" m="5px 10px">
          Online
        </Text>
        <MemberListItem />
      </UnorderedList>
    </GridItem>
  );
};
