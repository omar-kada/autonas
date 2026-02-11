import { getLogoutAPILogoutMutationOptions, getUserAPIGetQueryOptions } from '@/api/api';
import { useMutation } from '@tanstack/react-query';
import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { toast } from 'sonner';

export const getLogoutOptions = () => {
  return getLogoutAPILogoutMutationOptions({
    mutation: {
      onSuccess: (data, _, __, context) => {
        if (data.data.success) {
          context.client.refetchQueries(getUserAPIGetQueryOptions());
        }
      },
    },
  });
};

export const useLogout = () => {
  const { t } = useTranslation();

  const logoutMutation = useMutation(getLogoutOptions());

  const handlelogout = useCallback(() => {
    return toast
      .promise(
        () =>
          logoutMutation.mutateAsync().then((res) => {
            return res.data?.success;
          }),
        {
          error: t('ALERT.LOGOUT_ERROR'),
        },
      )
      ?.unwrap();
  }, [logoutMutation.mutateAsync, t]);

  return {
    ...logoutMutation,
    logout: handlelogout,
  };
};
