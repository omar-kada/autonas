import { getRegisterAPIRegisteredQueryOptions } from '@/api/api';
import { useQuery } from '@tanstack/react-query';

export function useRegistered() {
  return useQuery(
    getRegisterAPIRegisteredQueryOptions({
      query: {
        select: (data) => {
          return data.data?.registered;
        },
      },
    }),
  );
}
