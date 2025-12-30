import { getStatusAPIGetQueryOptions } from '@/api/api';

export const getStatusQueryOptions = () => {
  return getStatusAPIGetQueryOptions({
    query: {
      select: (data) => data?.data,
      staleTime: 10 * 1000,
      gcTime: 60 * 10 * 1000,
    },
  });
};
