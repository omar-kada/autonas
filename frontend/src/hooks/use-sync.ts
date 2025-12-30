import { getDeployementAPISyncMutationOptions } from '@/api/api';
import { getDeploymentsQueryOptions } from './use-deployments';

export const getSyncOptions = () => {
  return getDeployementAPISyncMutationOptions({
    mutation: {
      onSuccess: (data, _, __, context) => {
        if (data.data?.id && data.data?.id !== '0') {
          context.client.refetchQueries(getDeploymentsQueryOptions());
        }
      },
    },
  });
};
