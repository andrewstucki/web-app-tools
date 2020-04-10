import { applyMiddleware, combineReducers, createStore } from "redux";
import { composeWithDevTools } from "redux-devtools-extension";
import { AxiosInstance } from "axios";
import { authenticatedMiddleware } from "@andrewstucki/web-app-tools-middleware";

import profile, { ProfileState } from "./reducer";

const tokenKey = "__google_id";
const headerKey = "x-google-id";
const authenticatedUrl = "/oauth";

export interface RootState {
  profile: ProfileState;
}

export default (client: AxiosInstance) =>
  createStore(
    combineReducers<RootState>({ profile }),
    composeWithDevTools(
      applyMiddleware(
        authenticatedMiddleware(client, tokenKey, headerKey, authenticatedUrl)
      )
    )
  );
