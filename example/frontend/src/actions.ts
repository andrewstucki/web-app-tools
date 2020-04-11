import { authenticatedRequest } from '@andrewstucki/web-app-tools-middleware';
import { AxiosError, AxiosResponse } from 'axios';

import { ProfileResponse, APIError, createProfileResponseFrom } from './models';

export const PROFILE_INIT = 'PROFILE_INIT';
export const PROFILE_SUCCESS = 'PROFILE_SUCCESS';
export const PROFILE_ERROR = 'PROFILE_ERROR';

type ProfileSuccess = {
  type: typeof PROFILE_SUCCESS;
  payload: ProfileResponse;
};
export const profileSuccess = (profile: ProfileResponse): ProfileSuccess => ({
  type: PROFILE_SUCCESS,
  payload: profile,
});

type ProfileError = {
  type: typeof PROFILE_ERROR;
  payload: Error;
};
export const profileError = (error: Error): ProfileError => ({
  type: PROFILE_ERROR,
  payload: error,
});

type ProfileInit = {
  type: typeof PROFILE_INIT;
};
export const profileInit = (): ProfileInit => ({
  type: PROFILE_INIT,
});

export const getProfile = () => {
  return authenticatedRequest<ProfileResponse, APIError>({
    config: {
      url: '/api/v1/me',
    },
    onStart: () => profileInit(),
    onError: (error: AxiosError<APIError>) => {
      const { response } = error;
      if (response) return profileError(new Error(response.data.reason));
      if (error.request) {
        return profileError(new Error('no response received'));
      }
      return profileError(error);
    },
    convertData: createProfileResponseFrom,
    onResponse: (response: AxiosResponse<ProfileResponse>) =>
      profileSuccess(response.data),
  });
};

export type RootAction = ProfileInit | ProfileSuccess | ProfileError;
