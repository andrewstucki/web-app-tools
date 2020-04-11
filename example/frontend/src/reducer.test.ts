import reducer, { initialState } from './reducer';
import { profileInit, profileSuccess, profileError } from './actions';

describe('reducer', () => {
  it('tracks loading state', () => {
    expect(reducer(initialState, profileInit()).loading).toEqual(true);
  });

  it('handles profile data', () => {
    const payload = {
      user: {
        id: 'foo',
        email: 'bar@baz.com',
        createdAt: new Date(),
        updatedAt: new Date(),
      },
      policies: [{ resource: '*', action: '*' }],
    };
    const state = reducer(
      { ...initialState, loading: true },
      profileSuccess(payload)
    );
    expect(state).toEqual(
      expect.objectContaining({ loading: false, ...payload })
    );
  });

  it('handles profile errors', () => {
    const error = new Error('foo bar');
    expect(reducer(initialState, profileError(error))).toEqual(
      expect.objectContaining({ error })
    );
  });
});
