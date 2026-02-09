import { getUserAPIGetQueryOptions } from '@/api/api';
import { useQuery } from '@tanstack/react-query';

export function useUser() {
  return useQuery(
    getUserAPIGetQueryOptions({
      query: {
        select: (data) => data.data,
      },
    }),
  );
}
