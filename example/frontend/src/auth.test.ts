import { Policy } from './models';

const injectState = (mockPolicies: Policy[] = []) => {
  jest.mock('react-redux', () => ({
    ...jest.requireActual('react-redux'),
    useSelector: (): Policy[] => mockPolicies,
  }));
  // eslint-disable-next-line global-require
  const { can } = require('./auth');
  return can;
};

describe('auth', () => {
  beforeEach(() => {
    jest.resetModules();
  });

  it('handles unauthorized users', () => {
    const can = injectState();
    expect(can('edit', 'post')).toEqual(false);
  });

  it('handles admins', () => {
    const can = injectState([{ action: '*', resource: '*' }]);
    expect(can('edit', 'post')).toEqual(true);
  });

  it('handles authorized users', () => {
    const can = injectState([{ action: 'edit', resource: 'post' }]);
    expect(can('edit', 'post')).toEqual(true);
  });
});
