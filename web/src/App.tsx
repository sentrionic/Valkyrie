import * as React from "react";
import { QueryClient, QueryClientProvider } from "react-query";
import { BrowserRouter, Switch, Route } from "react-router-dom";
import { Landing } from "./pages/Landing";
import { Login } from "./pages/auth/Login";
import { Register } from "./pages/auth/Register";
import { Account } from "./pages/Account";
import { AuthRoute } from "./components/layouts/AuthRoute";
import { ForgotPassword } from "./pages/ForgotPassword";
import { ResetPassword } from "./pages/ResetPassword";
import { ViewGuild } from "./pages/ViewGuild";
import { Home } from "./components/layouts/Home";

export const App = () => (
  <QueryClientProvider client={new QueryClient()}>
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
        <Route exact path="/channels/me">
          <Home />
        </Route>
        <Route path="/channels">
          <ViewGuild />
        </Route>
        <AuthRoute path="/account" component={Account} />
        <Route path="/">
          <Landing />
        </Route>
      </Switch>
    </BrowserRouter>
  </QueryClientProvider>
);
