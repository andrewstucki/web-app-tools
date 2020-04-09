import { AnyAction, Dispatch, Middleware, MiddlewareAPI } from "redux";
import {
  AxiosInstance,
  AxiosError,
  AxiosResponse,
  AxiosRequestConfig
} from "axios";

import {
  AUTHENTICATED_REQUEST,
  AuthenticatedRequest,
  AUTHENTICATED_LOG_OUT,
  AuthenticatedRequestConfig
} from "./action";

const callbacks = ["onUploadProgress", "onDownloadProgress"];

const defaultRequestConfig: AxiosRequestConfig = {};

function mergeDefaultRequestConfig(
  config?: AuthenticatedRequestConfig
): AxiosRequestConfig {
  if (config) return { ...defaultRequestConfig, ...config };
  return defaultRequestConfig;
}

function clearTokenAndRedirect(tokenKey: string, authenticateUrl: string) {
  localStorage.removeItem(tokenKey);
  if (authenticateUrl !== "") {
    window.location.href = authenticateUrl;
  }
}

function getToken(tokenKey: string): string | null {
  const token = localStorage.getItem(tokenKey);
  if (!token) return null;
  return token;
}

function setToken(tokenKey: string, token: string): Error | void {
  try {
    localStorage.setItem(tokenKey, token);
  } catch (e) {
    return e;
  }
}

function setAuthorizationHeader(
  config: AxiosRequestConfig,
  token: string | null
): AxiosRequestConfig {
  if (!token || token === "") return config;
  return {
    ...config,
    headers: {
      ...config.headers,
      Authorization: "bearer " + token
    }
  };
}

/* tslint:disable:no-any*/
function setupCallbacks(dispatch: any, config: any, payload: any) {
  callbacks.forEach(attribute => {
    const callback = payload[attribute];
    if (callback) {
      config[attribute] = (arg: any) => {
        const actionToDispatch = callback(arg);
        actionToDispatch && dispatch(actionToDispatch);
      };
    }
  });
}

export function authenticatedMiddleware(
  instance: AxiosInstance,
  tokenKey: string,
  headerKey: string,
  authenticateUrl: string
): Middleware {
  let token = getToken(tokenKey);

  const tryRefreshTokenFromHeaders = (headers: any) => {
    const refreshedToken: string | null = headers[headerKey];
    if (refreshedToken && refreshedToken !== "") {
      const error = setToken(tokenKey, refreshedToken);
      if (error) {
        console.log(error);
      } else {
        token = refreshedToken;
      }
    }
  };

  return ({ dispatch }: MiddlewareAPI<any>) => (next: Dispatch<AnyAction>) => (
    action: any
  ) => {
    if (action.type === AUTHENTICATED_LOG_OUT)
      clearTokenAndRedirect(tokenKey, authenticateUrl);

    if (!isAuthenticatedRequest(action)) {
      return next(action);
    }

    const { payload } = action;

    if (payload.onStart) {
      const actionToDispatch = payload.onStart();
      actionToDispatch && dispatch(actionToDispatch);
    }

    const config = mergeDefaultRequestConfig(payload.config);
    setupCallbacks(dispatch, config, payload);

    instance
      .request(setAuthorizationHeader(config, token))
      .then((response: AxiosResponse<any>) => {
        tryRefreshTokenFromHeaders(response.headers);
        const actionToDispatch = payload.onResponse(response);
        return actionToDispatch && dispatch(actionToDispatch);
      })
      .catch((error: AxiosError<any>) => {
        if (error.response) {
          tryRefreshTokenFromHeaders(error.response.headers);
          if (error.response.status === 401)
            return clearTokenAndRedirect(tokenKey, authenticateUrl);
        }
        const actionToDispatch = payload.onError(error);
        return actionToDispatch && dispatch(actionToDispatch);
      });

    return next(action);
  };
}

function isAuthenticatedRequest(
  action: any
): action is AuthenticatedRequest<any, any> {
  return (
    action &&
    action.type &&
    action.type === AUTHENTICATED_REQUEST &&
    isAuthenticatedRequestPayload(action)
  );
}

function isAuthenticatedRequestPayload(action: any): boolean {
  return (
    action &&
    action.payload &&
    action.payload.onResponse &&
    action.payload.onError
  );
}
/* tslint:enable:no-any*/
