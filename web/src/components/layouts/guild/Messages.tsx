import { GridItem } from "@chakra-ui/react";
import React from "react";
import { Message } from "../../common/Message";

export const Messages: React.FC = () => {
  return (
    <GridItem
      gridColumn={3}
      gridRow={"2"}
      bg="brandGray.light"
      mr="5px"
      display="flex"
      flexDirection="column-reverse"
      overflowY="auto"
      css={{
        "&::-webkit-scrollbar": {
          width: "8px",
        },
        "&::-webkit-scrollbar-track": {
          background: "#2f3136",
          width: "10px",
        },
        "&::-webkit-scrollbar-thumb": {
          background: "#202225",
          borderRadius: "18px",
        },
      }}
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
    </GridItem>
  );
};
