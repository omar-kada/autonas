import { useDeployementAPIList } from '@/api/api';

export const useDeployments = () => {
  const { data, isLoading, error } = useDeployementAPIList({
    query: {
      refetchInterval: 50000,
    },
  });

  return {
    data: data?.data ?? [],
    isLoading,
    error,
  };
};
