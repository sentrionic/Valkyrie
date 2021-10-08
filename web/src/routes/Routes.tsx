import React from 'react';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import { Login } from './Login';
import { Register } from './Register';
import { ForgotPassword } from './ForgotPassword';
import { ResetPassword } from './ResetPassword';
import { Home } from './Home';
import { ViewGuild } from './ViewGuild';
import { AuthRoute } from './AuthRoute';
import { Settings } from './Settings';
import { Landing } from './Landing';
import { Invite } from './Invite';

export const Routes: React.FC = () => (
  <BrowserRouter>
    <Switch>
      <Route path="/login">
        <Login />
      </Route>
      <Route path="/register">
        <Register />
      </Route>
      <Route path="/forgot-password">
        <ForgotPassword />
      </Route>
      <Route path="/reset-password/:token">
        <ResetPassword />
      </Route>
      <Route exact path="/">
        <Landing />
      </Route>
      <AuthRoute exact path="/channels/me" component={Home} />
      <AuthRoute exact path="/channels/me/:channelId" component={Home} />
      <AuthRoute exact path="/channels/:guildId/:channelId" component={ViewGuild} />
      <AuthRoute exact path="/account" component={Settings} />
      <AuthRoute exact path="/:link" component={Invite} />
    </Switch>
  </BrowserRouter>
);
