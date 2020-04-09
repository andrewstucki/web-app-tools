import {
  RootAction,
  CURRENT_USER_INIT,
  CURRENT_USER_SUCCESS,
  CURRENT_USER_ERROR
} from "./actions";
import { User } from "./models";

export type ProfileState = {
  readonly user: User | null;
  readonly error: Error | null;
  readonly loading: boolean;
};

const initialState = {
  user: null,
  error: null,
  loading: false
};

export default function(
  state: ProfileState = initialState,
  action: RootAction
): ProfileState {
  switch (action.type) {
    case CURRENT_USER_INIT:
      return { ...initialState, loading: true };

    case CURRENT_USER_SUCCESS:
      return { ...state, loading: false, user: action.payload };

    case CURRENT_USER_ERROR:
      return { ...state, loading: false, error: action.payload };

    default:
      return state;
  }
}
