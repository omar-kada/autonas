import { getDeployementAPIReadQueryOptions } from '@/api/api';
import { useQueryClient } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';

export const ROUTES = {
  DEPLOYMENTS: '/deployments',
  DEPLOYMENT: (id?: string) => `${ROUTES.DEPLOYMENTS}/${id}`,
  STATUS: '/status',
  CONFIG: '/config',
  LOGS: '/logs',
};

export function useDeploymentNavigate() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  return (id?: string | null) => {
    if (id != null && id !== '0') {
      navigate(ROUTES.DEPLOYMENT(id));
      queryClient.prefetchQuery(getDeployementAPIReadQueryOptions(id));
    } else {
      navigate(ROUTES.DEPLOYMENTS);
    }
  };
}
