import { Grid } from '@chakra-ui/react';
import React from 'react';

interface AppLayoutProps {
  showLastColumn?: boolean | null
}

export const AppLayout: React.FC<AppLayoutProps> = ({ showLastColumn = false, children }) => {
  // Col Chat: GuildList ChannelList Chat [MemberList]
  // Col Home: GuildList DMList Chat/FriendsList
  return (
    <Grid
      height='100vh'
      templateColumns={`75px 240px 1fr ${showLastColumn ? "240px" : ""} `}
      templateRows='auto 1fr auto'
      bg='brandGray.light'
    >
      {children}
    </Grid>
  );
};
