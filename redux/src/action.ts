import { Action } from "redux";
import {
  AxiosTransformer,
  AxiosBasicCredentials,
  AxiosAdapter,
  Method,
  AxiosProxyConfig,
  CancelToken,
  AxiosError,
  AxiosResponse,
  ResponseType,
} from "axios";

export type AuthenticatedRequestConfig = {
  // taken from Axios' typings
  url?: string;
  method?: Method;
  baseURL?: string;
  transformRequest?: AxiosTransformer | AxiosTransformer[];
  transformResponse?: AxiosTransformer | AxiosTransformer[];
  headers?: any;
  params?: any;
  paramsSerializer?: (params: any) => string;
  data?: any;
  timeout?: number;
  timeoutErrorMessage?: string;
  withCredentials?: boolean;
  adapter?: AxiosAdapter;
  auth?: AxiosBasicCredentials;
  responseType?: ResponseType;
  xsrfCookieName?: string;
  xsrfHeaderName?: string;
  maxContentLength?: number;
  maxBodyLength?: number;
  validateStatus?: (status: number) => boolean;
  maxRedirects?: number;
  socketPath?: string | null;
  httpAgent?: any;
  httpsAgent?: any;
  proxy?: AxiosProxyConfig | false;
  cancelToken?: CancelToken;
  decompress?: boolean;
};

// Descriptor of a authenticated request payload
export type AuthenticatedRequestPayload<T = any, U = any> = {
  config?: AuthenticatedRequestConfig;
  // Called immediately before the request is started, useful for toggling a loading status
  onStart?: () => Action | void;
  // Called when Axios receives upload progress updates
  onUploadProgress?: (progressEvent: ProgressEvent) => Action | void;
  // Called when Axios receives download progress updates
  onDownloadProgress?: (progressEvent: ProgressEvent) => Action | void;
  // Called on a response
  convertData?: (data: any) => T;
  // Called on a response
  onResponse: (response: AxiosResponse<T>) => Action | void;
  // Called on an error
  onError: (error: AxiosError<U>) => Action | void;
};
export const AUTHENTICATED_REQUEST = "__AUTHENTICATED_REQUEST";
export const AUTHENTICATED_LOG_OUT = "__AUTHENTICATED_LOG_OUT";

// Basic type for a log out action
export type AuthenticatedLogOut = {
  type: typeof AUTHENTICATED_LOG_OUT;
};

// Action creator, Use it to create a new log out action
export const authenticatedLogOut = (): AuthenticatedLogOut => ({
  type: AUTHENTICATED_LOG_OUT,
});

// Basic type for a request action
export type AuthenticatedRequest<T, U> = {
  type: typeof AUTHENTICATED_REQUEST;
  payload: AuthenticatedRequestPayload<T, U>;
};

// Action creator, Use it to create a new authenticated request action
export function authenticatedRequest<T, U>(
  payload: AuthenticatedRequestPayload<T, U>
): AuthenticatedRequest<T, U> {
  return {
    type: AUTHENTICATED_REQUEST,
    payload,
  };
}

export type AuthenticatedActionTypes<T, U> =
  | AuthenticatedRequest<T, U>
  | AuthenticatedLogOut;
