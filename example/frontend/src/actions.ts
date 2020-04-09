import { User, ApiError } from "./models";

import { authenticatedRequest } from "@andrewstucki/google-oauth-tools-middleware";
import { AxiosError, AxiosResponse } from "axios";

export const CURRENT_USER_INIT = "CURRENT_USER_INIT";
export const CURRENT_USER_SUCCESS = "CURRENT_USER_SUCCESS";
export const CURRENT_USER_ERROR = "CURRENT_USER_ERROR";

type CurrentUserSuccess = {
  type: typeof CURRENT_USER_SUCCESS;
  payload: User;
};
export const currentUserSuccess = (user: User): CurrentUserSuccess => ({
  type: CURRENT_USER_SUCCESS,
  payload: user
});

type CurrentUserError = {
  type: typeof CURRENT_USER_ERROR;
  payload: Error;
};
export const currentUserError = (error: Error): CurrentUserError => ({
  type: CURRENT_USER_ERROR,
  payload: error
});

type CurrentUserInit = {
  type: typeof CURRENT_USER_INIT;
};
export const currentUserInit = (): CurrentUserInit => ({
  type: CURRENT_USER_INIT
});

export const getCurrentUser = () => {
  return authenticatedRequest<User, ApiError>({
    config: {
      url: "/api/v1/me"
    },
    onStart: () => currentUserInit(),
    onError: (error: AxiosError<ApiError>) => {
      const response = error.response;
      if (response) return currentUserError(new Error(response.data.error));
      if (error.request) {
        return currentUserError(new Error("no response received"));
      }
      return currentUserError(error);
    },
    onResponse: (response: AxiosResponse<User>) =>
      currentUserSuccess(response.data)
  });
};

export type RootAction =
  | CurrentUserInit
  | CurrentUserSuccess
  | CurrentUserError;
