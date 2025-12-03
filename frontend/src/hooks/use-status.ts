import type { StackState } from '@/models/stack-status';
import { useQuery } from '@tanstack/react-query';

export const useStatus = () => {
  return useQuery<StackState, Error>({
    queryKey: ['status'],
    queryFn: async () => {
      const response = await fetch('/api/status');
      return response.json();
    },
  });
};
