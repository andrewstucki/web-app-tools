import {
  RootAction,
  PROFILE_INIT,
  PROFILE_SUCCESS,
  PROFILE_ERROR,
} from "./actions";
import { User, Policy } from "./models";

export type ProfileState = {
  readonly user: User | null;
  readonly policies: Policy[];
  readonly error: Error | null;
  readonly loading: boolean;
};

export const initialState = {
  user: null,
  policies: [],
  error: null,
  loading: false,
};

export default function (
  state: ProfileState = initialState,
  action: RootAction
): ProfileState {
  switch (action.type) {
    case PROFILE_INIT:
      return { ...initialState, loading: true };

    case PROFILE_SUCCESS:
      return {
        ...state,
        loading: false,
        ...action.payload,
      };

    case PROFILE_ERROR:
      return { ...state, loading: false, error: action.payload };

    default:
      return state;
  }
}
