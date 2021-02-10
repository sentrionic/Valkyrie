import React from "react";
import {
  Redirect,
  Route,
  RouteComponentProps,
  RouteProps,
} from "react-router-dom";
import { getCurrent } from "../lib/stores/userStore";

interface IProps extends RouteProps {
  component: React.ComponentType<RouteComponentProps<any>>;
}

export const AuthRoute: React.FC<IProps> = ({
  component: Component,
  ...rest
}) => {
  const storage = JSON.parse(sessionStorage.getItem("user-storage")!!);
  const current = getCurrent();
  return (
    <Route
      {...rest}
      render={(props) =>
        current || storage?.state?.current ? (
          <Component {...props} />
        ) : (
          <Redirect to="/login" />
        )
      }
    />
  );
};
