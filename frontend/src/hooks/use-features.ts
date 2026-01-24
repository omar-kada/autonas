import { getFeaturesAPIGetQueryOptions } from '@/api/api';

export const getFeaturesQueryOptions = () => {
  return getFeaturesAPIGetQueryOptions({
    query: {
      select: (data) => data?.data ?? {},
      staleTime: Infinity,
      gcTime: 10 * 60 * 1000,
    },
  });
};
