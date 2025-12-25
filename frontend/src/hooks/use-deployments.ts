import { getDeployementAPIListQueryOptions, useDeployementAPIList } from '@/api/api';
import type { QueryClient } from '@tanstack/react-query';

export const useDeployments = () => {
  const { data, isLoading, error } = useDeployementAPIList({
    query: {
      refetchInterval: 50000,
    },
  });

  return {
    deployments: data?.data ?? [],
    isLoading,
    error,
  };
};

export function refetchDeployments(queryClient: QueryClient) {
  queryClient.refetchQueries(getDeployementAPIListQueryOptions());
}
