import { useDeployementAPIRead, type Deployment, type Error } from '@/api/api';
import type { AxiosError } from 'axios';

export const useDeployment = (
  id?: string,
): {
  deployment?: Deployment;
  isLoading?: boolean;
  error?: AxiosError<Error, unknown> | null;
} => {
  if (id == null) {
    return {};
  }
  const { data, isLoading, error } = useDeployementAPIRead(id);

  return {
    deployment: data?.data,
    isLoading,
    error,
  };
};
