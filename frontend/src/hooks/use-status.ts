import { useStatusAPIGet } from '@/api/api';

export const useStatus = () => {
  const { data, isLoading, error } = useStatusAPIGet();

  return {
    data: data?.data,
    isLoading,
    error,
  };
};
