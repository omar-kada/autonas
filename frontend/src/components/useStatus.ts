import { useQuery } from '@tanstack/react-query';

export type ContainerState =
  | 'created'
  | 'running'
  | 'paused'
  | 'restarting'
  | 'removing'
  | 'exited'
  | 'dead';

export type HealthStatus = 'healthy' | 'unhealthy' | 'starting' | 'none';

// Define the type for the status response
type StatusResponse = {
  [key: string]: Array<{
    ID: string;
    State: ContainerState;
    Name: string;
    Health: HealthStatus;
  }>;
};

// Custom hook to fetch the status
export const useStatus = () => {
  return useQuery<StatusResponse, Error>({
    queryKey: ['status'],
    queryFn: async () => {
      const response = await fetch('/api/status');
      return response.json();
    },
  });
};
