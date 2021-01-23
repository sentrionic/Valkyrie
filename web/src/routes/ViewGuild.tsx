import React from 'react';
import { Channels } from '../components/layouts/chat/Channels';
import { GuildList } from '../components/layouts/chat/GuildList';
import { ChannelHeader } from '../components/layouts/chat/ChannelHeader';
import { MemberList } from '../components/layouts/chat/MemberList';
import { MessageInput } from '../components/layouts/chat/MessageInput';
import { ChatScreen } from '../components/layouts/chat/ChatScreen';
import { AppLayout } from '../components/layouts/AppLayout';

export const ViewGuild: React.FC = () => {
  return (
    <AppLayout>
      <GuildList />
      <Channels />
      <ChannelHeader />
      <ChatScreen />
      <MessageInput />
      <MemberList />
    </AppLayout>
  );
};
