import { getStatsAPIGetQueryOptions, type Error as ApiError, type Stats } from '@/api/api';
import type { UseQueryOptions } from '@tanstack/react-query';
import type { AxiosError, AxiosResponse } from 'axios';

export const getStatsQueryOptions = (
  queryOptions?: Partial<
    UseQueryOptions<AxiosResponse<Stats, unknown>, AxiosError<ApiError>, Stats>
  >,
) => {
  return getStatsAPIGetQueryOptions(30, {
    query: {
      select: (data) => data?.data,
      refetchInterval: 20 * 1000,
      refetchIntervalInBackground: false,
      staleTime: 0,
      gcTime: 10 * 60 * 1000,
      ...queryOptions,
    },
  });
};
