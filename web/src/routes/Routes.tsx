import React from 'react';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import { Login } from './Login';
import { Register } from './Register';
import { ForgotPassword } from './ForgotPassword';
import { ResetPassword } from './ResetPassword';
import { Home } from './Home';
import { ViewGuild } from './ViewGuild';
import { AuthRoute } from './AuthRoute';
import { Account } from './Account';
import { Landing } from './Landing';

export const Routes: React.FC = () => {
  return (
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
        <AuthRoute exact path="/channels/me" component={Home} />
        <AuthRoute path="/channels/:guildId" component={ViewGuild} />
        <AuthRoute path="/account" component={Account} />
        <Route path="/">
          <Landing />
        </Route>
      </Switch>
    </BrowserRouter>
  );
}
