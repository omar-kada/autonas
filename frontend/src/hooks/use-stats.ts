import { getStatsAPIGetQueryOptions } from '@/api/api';

export const getStatsQueryOptions = (days = 30) => {
  return getStatsAPIGetQueryOptions(days, {
    query: {
      select: (data) => data?.data,
      refetchInterval: 20 * 1000,
      staleTime: 0,
      gcTime: 10 * 60 * 1000,
    },
  });
};
