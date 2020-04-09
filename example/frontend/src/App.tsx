import React, { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { RootState } from "./store";
import { User } from "./models";
import { getCurrentUser } from "./actions";
import { authenticatedLogOut } from "@andrewstucki/google-oauth-tools-middleware";

export const App = () => {
  const user = useSelector<RootState, User | null>(state => state.profile.user);
  const dispatch = useDispatch();

  useEffect(() => {
    dispatch(getCurrentUser());
  }, [dispatch]);

  if (user) {
    return (
      <div>
        <span>Hi: {user.email}</span>
        <button
          onClick={() => {
            dispatch(authenticatedLogOut());
          }}
        >
          Log Out
        </button>
      </div>
    );
  }
  return <div>Loading</div>;
};
