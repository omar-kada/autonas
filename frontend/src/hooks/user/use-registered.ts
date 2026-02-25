import { getAuthAPIRegisteredQueryOptions } from '@/api/api';
import { useQuery } from '@tanstack/react-query';

export function useRegistered() {
  return useQuery(
    getAuthAPIRegisteredQueryOptions({
      query: {
        select: (data) => {
          return data.data?.registered;
        },
      },
    }),
  );
}
