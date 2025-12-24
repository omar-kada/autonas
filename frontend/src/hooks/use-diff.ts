import { useDiffAPIGet } from '@/api/api';

export const useDiff = () => {
  const { data, isLoading, error } = useDiffAPIGet();

  return {
    data: data?.data,
    isLoading,
    error,
  };
};
