import React from 'react';
import { GuildList } from '../components/layouts/chat/GuildList';
import { DMSidebar } from '../components/layouts/home/DMSidebar';
import { FriendList } from '../components/layouts/home/FriendList';
import { AppLayout } from '../components/layouts/AppLayout';
import { useParams } from 'react-router-dom';
import { RouterProps } from './Routes';
import { ChatScreen } from '../components/layouts/chat/ChatScreen';
import { DMHeader } from '../components/layouts/home/DMHeader';

export const Home: React.FC = () => {

  const { channelId } = useParams<RouterProps>();

  return (
    <AppLayout>
      <GuildList />
      <DMSidebar />
      {channelId === undefined ?
        <FriendList />
        :
        <>
          <DMHeader />
          <ChatScreen />
          {/* DMMessageInput */}
        </>
      }
    </AppLayout>
  );
}
