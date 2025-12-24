import { useNavigate } from 'react-router-dom';

export const ROUTES = {
  DEPLOYMENTS: '/deployments',
  DEPLOYMENT: (id: string) => `/deployments/${id}`,
  STATUS: '/status',
  CONFIG: '/config',
  LOGS: '/logs',
};

export function useDeploymentNavigate() {
  const navigate = useNavigate();
  return (id?: string | null) => {
    if (id != null && id != '0') {
      navigate(ROUTES.DEPLOYMENT(id));
    } else {
      navigate(ROUTES.DEPLOYMENTS);
    }
  };
}
