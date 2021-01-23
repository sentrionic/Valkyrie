import { Grid } from '@chakra-ui/react';
import React from 'react';

export const AppLayout: React.FC = ({ children }) => {
  // Col: GuildList ChannelList Chat MemberList
  return (
    <Grid
      height='100vh'
      templateColumns='75px 240px 1fr 240px'
      templateRows='auto 1fr auto'
      bg='brandGray.light'
    >
      {children}
    </Grid>
  );
};
