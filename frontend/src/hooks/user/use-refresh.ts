import { getAuthAPIRefreshMutationOptions, getUserAPIGetQueryOptions } from '@/api/api';
import { useMutation } from '@tanstack/react-query';
import { useCallback } from 'react';

export const getRefreshOptions = () => {
  return getAuthAPIRefreshMutationOptions({
    mutation: {
      onSuccess: (data, _, __, context) => {
        if (data.data.success) {
          context.client.refetchQueries(getUserAPIGetQueryOptions());
        }
      },
    },
  });
};

export const useRefresh = () => {
  const refreshMutation = useMutation(getRefreshOptions());

  const handleRefresh = useCallback(() => {
    return refreshMutation
      .mutateAsync()
      .then((res) => res.data?.success)
      .catch(() => false);
  }, [refreshMutation.mutateAsync]);

  return {
    ...refreshMutation,
    refresh: handleRefresh,
  };
};
