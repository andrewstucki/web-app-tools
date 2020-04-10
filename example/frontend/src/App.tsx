import React, { useEffect } from "react";
import moment from 'moment';
import { useDispatch, useSelector } from "react-redux";
import { RootState } from "./store";
import { User } from "./models";
import { getProfile } from "./actions";
import { authenticatedLogOut } from "@andrewstucki/web-app-tools-middleware";

export const App = () => {
  const user = useSelector<RootState, User | null>(
    (state) => state.profile.user
  );
  const dispatch = useDispatch();

  useEffect(() => {
    dispatch(getProfile());
  }, [dispatch]);

  if (user) {
    return (
      <div>
        <h1>Logged in as:</h1>
        <table>
          <tr>
            <td>ID</td>
            <td>{user.id}</td>
          </tr>
          <tr>
            <td>Email</td>
            <td>{user.email}</td>
          </tr>
          <tr>
            <td>Created At</td>
            <td>{moment(user.createdAt).fromNow()}</td>
          </tr>
          <tr>
            <td>Last Updated</td>
            <td>{moment(user.updatedAt).fromNow()}</td>
          </tr>
        </table>
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
