import React from "react";
import { GuildList } from "../components/layouts/chat/GuildList";
import { DMSidebar } from "../components/layouts/home/DMSidebar";
import { FriendList } from "../components/layouts/home/FriendList";
import { AppLayout } from '../components/layouts/AppLayout';

export const Home: React.FC = () => {
  return (
    <AppLayout>
      <GuildList />
      <DMSidebar />
      <FriendList />
    </AppLayout>
  );
};
