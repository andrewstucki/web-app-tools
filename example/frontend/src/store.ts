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
export const initialState: RootState = {
  profile,
}
export const reducer = combineReducers<RootState>(initialState);

export default (client: AxiosInstance) =>
  createStore(
    reducer,
    composeWithDevTools(
      applyMiddleware(
        authenticatedMiddleware(client, tokenKey, headerKey, authenticatedUrl)
      )
    )
  );
