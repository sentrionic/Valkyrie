import { Grid } from "@chakra-ui/react";
import React from "react";
import { Guilds } from "./guild/Guilds";
import { DMSidebar } from "./home/DMSidebar";
import { FriendList } from "./home/FriendList";

export const Home: React.FC = () => {
  return (
    <Grid
      height="100vh"
      templateColumns="75px 240px auto"
      templateRows="auto 1fr auto"
      bg="brandGray.light"
    >
      <Guilds />
      <DMSidebar />
      <FriendList />
    </Grid>
  );
};
