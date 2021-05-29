import React from 'react';
import { Channels } from '../components/layouts/guild/Channels';
import { GuildList } from '../components/layouts/guild/GuildList';
import { ChannelHeader } from '../components/layouts/guild/ChannelHeader';
import { MemberList } from '../components/layouts/guild/MemberList';
import { MessageInput } from '../components/layouts/guild/chat/MessageInput';
import { ChatScreen } from '../components/layouts/guild/chat/ChatScreen';
import { AppLayout } from '../components/layouts/AppLayout';
import { settingsStore } from '../lib/stores/settingsStore';

export const ViewGuild: React.FC = () => {
  const showMemberList = settingsStore((state) => state.showMembers);

  return (
    <AppLayout showLastColumn={showMemberList}>
      <GuildList />
      <Channels />
      <ChannelHeader />
      <ChatScreen />
      <MessageInput />
      {showMemberList && <MemberList />}
    </AppLayout>
  );
};
