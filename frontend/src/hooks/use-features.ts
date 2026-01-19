import { getFeaturesAPIGetQueryOptions, type Features } from '@/api/api';
import type { AxiosResponse } from 'axios';

export const getFeaturesQueryOptions = () => {
  return getFeaturesAPIGetQueryOptions({
    query: {
      select: (data) => data?.data ?? {},
      placeholderData: { data: {} as Features } as AxiosResponse,
      staleTime: Infinity,
      gcTime: 10 * 60 * 1000,
    },
  });
};
