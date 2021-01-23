import { GridItem } from "@chakra-ui/react";
import React from "react";
import { Message } from "../../items/Message";
import { StartMessages } from '../../sections/StartMessages';
import { scrollbarCss } from '../../../lib/utils/theme';

export const ChatScreen: React.FC = () => {
  return (
    <GridItem
      gridColumn={3}
      gridRow={"2"}
      bg="brandGray.light"
      mr="5px"
      display="flex"
      flexDirection="column-reverse"
      overflowY="auto"
      css={scrollbarCss}
    >
      <Message />
      <Message />
      <Message />
      <Message />
      <Message />
      <Message />
      <Message />
      <Message />
      <Message />
      <Message />
      <Message />
      <Message />
      <Message />
      <StartMessages />
    </GridItem>
  );
};
