import { useDeployementAPIRead, type Deployment, type Error } from '@/api/api';
import type { AxiosError } from 'axios';

export const useDeployment = (
  id: string,
): {
  deployment?: Deployment;
  isPending?: boolean;
  error?: AxiosError<Error, unknown> | null;
} => {
  const { data, isPending, error } = useDeployementAPIRead(id);

  return {
    deployment: data?.data,
    isPending,
    error,
  };
};
