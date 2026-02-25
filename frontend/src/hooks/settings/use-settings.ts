import { getSettingsAPIGetQueryOptions, type Error, type Settings } from '@/api/api';
import type { AxiosError } from 'axios';

export const getSettingsQueryOptions = () => {
  return getSettingsAPIGetQueryOptions<Settings, AxiosError<Error>>({
    query: {
      select: (data) => data?.data,
      gcTime: 10 * 60 * 1000,
    },
  });
};
