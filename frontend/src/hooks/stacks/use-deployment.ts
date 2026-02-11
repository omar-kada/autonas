import { DeploymentStatus, getDeployementAPIReadQueryOptions } from '@/api/api';

export const getDeploymentOptions = (id: string) => {
  return getDeployementAPIReadQueryOptions(id, {
    query: {
      select: (data) => data?.data,
      staleTime: (query) => {
        switch (query.state?.data?.data?.status) {
          case DeploymentStatus.running:
            return 500;
          case DeploymentStatus.error:
          case DeploymentStatus.success:
            return Infinity;
          default:
            return 10 * 1000;
        }
      },
      gcTime: 10 * 60 * 1000,
    },
  });
};
