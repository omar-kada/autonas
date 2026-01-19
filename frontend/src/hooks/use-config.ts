import { getConfigAPIGetQueryOptions, type Config } from '@/api/api';
import type { AxiosError } from 'axios';

export const getConfigQueryOptions = ({ enabled }: { enabled: boolean }) => {
  return getConfigAPIGetQueryOptions<Config, AxiosError<Error>>({
    query: {
      select: (data) => data?.data,
      gcTime: 10 * 60 * 1000,
      enabled,
    },
  });
};
