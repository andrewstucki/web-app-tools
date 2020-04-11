import { useSelector } from 'react-redux';

import { Policy } from './models';
import { RootState } from './store';

export const can = (requestAction: string, requestResource: string) =>
  useSelector<RootState, Policy[]>((state) => state.profile.policies).some(
    ({ action, resource }) => {
      const matchedAction = action === '*' || action === requestAction;
      const matchedResource = resource === '*' || resource === requestResource;
      return matchedAction && matchedResource;
    }
  );
