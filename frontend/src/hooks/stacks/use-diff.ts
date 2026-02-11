import { getDiffAPIGetQueryOptions } from '@/api/api';

export const getDiffQueryOptions = () => {
  return getDiffAPIGetQueryOptions({
    query: {
      select: (data) => data?.data,
      gcTime: 10 * 60 * 1000,
    },
  });
};
