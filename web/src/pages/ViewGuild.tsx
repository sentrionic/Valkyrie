import {
  Button,
  Grid,
  GridItem,
  Icon,
  IconButton,
  Input,
  InputGroup,
  InputLeftElement,
  InputRightElement,
} from "@chakra-ui/react";
import React from "react";
import { MdAddCircle } from "react-icons/md";
import { Channels } from "../components/layouts/guild/Channels";
import { Guilds } from "../components/layouts/guild/Guilds";
import { Header } from "../components/layouts/guild/Header";
import { MemberList } from "../components/layouts/guild/MemberList";
import { MessageInput } from "../components/layouts/guild/MessageInput";
import { Messages } from "../components/layouts/guild/Messages";

export const ViewGuild: React.FC = () => {
  return (
    <Grid
      height="100vh"
      templateColumns="75px 240px 1fr 240px"
      templateRows="auto 1fr auto"
      bg="brandGray.light"
    >
      <Guilds />
      <Channels />
      <Header />
      <Messages />
      <MessageInput />
      <MemberList />
    </Grid>
  );
};
