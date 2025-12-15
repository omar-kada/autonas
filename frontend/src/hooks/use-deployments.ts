import { useDeployementAPIList } from '@/api/api';

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
