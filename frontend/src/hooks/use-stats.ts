import { useStatsAPIGet } from '@/api/api';

export const useStats = (days: number) => {
  const { data, isLoading, error } = useStatsAPIGet(days, {
    query: {
      refetchInterval: 50000,
    },
  });

  return {
    data: data?.data,
    isLoading,
    error,
  };
};
