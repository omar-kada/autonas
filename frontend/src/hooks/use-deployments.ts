import { getDeployementAPIListQueryOptions, useDeployementAPIList } from '@/api/api';
import type { QueryClient } from '@tanstack/react-query';

export const useDeployments = () => {
  const { data, isLoading, error } = useDeployementAPIList(
    {
      after: '',
      first: 15,
    },
    {
      query: {
        refetchInterval: 50000,
      },
    },
  );

  return {
    deployments: data?.data ?? { data: [], pageInfo: { hasNextPage: false, endCursor: '' } },
    isLoading,
    error,
  };
};

export function refetchDeployments(queryClient: QueryClient) {
  queryClient.refetchQueries(
    getDeployementAPIListQueryOptions({
      after: '',
      first: 15,
    }),
  );
}
