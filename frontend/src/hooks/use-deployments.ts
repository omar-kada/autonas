import type { Deployment } from '@/models/deployment';
import { useQuery } from '@tanstack/react-query';

export const useDeployments = () => {
  return useQuery<Array<Deployment>, Error>({
    queryKey: ['deployments'],
    queryFn: async () => {
      //await fetch('/api/deployments');
      return [
        {
          id: '1',
          name: 'Deployment 1',
          time: '2024-01-02T12:00:00Z',
          status: 'running',
          diff: 'Diff details for Deployment 1',
        },
        {
          id: '2',
          name: 'Deployment 2',
          time: '2024-01-01T12:00:00Z',
          status: 'success',
          diff: 'Diff details for Deployment 2',
        },
      ];
    },
  });
};
