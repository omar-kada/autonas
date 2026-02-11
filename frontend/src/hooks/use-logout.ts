import { getLogoutAPILogoutMutationOptions, getUserAPIGetQueryOptions } from '@/api/api';
import { ROUTES } from '@/lib';
import { useMutation } from '@tanstack/react-query';
import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
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
  const navigate = useNavigate();

  const logoutMutation = useMutation(getLogoutOptions());

  const handlelogout = useCallback(() => {
    toast.promise(
      () =>
        logoutMutation.mutateAsync().then((res) => {
          if (res.data?.success) {
            navigate(ROUTES.LOGIN);
          }
          return res.data?.success;
        }),
      {
        error: t('ALERT.LOGOUT_ERROR'),
      },
    );
  }, [logoutMutation.mutateAsync, t]);

  return {
    ...logoutMutation,
    logout: handlelogout,
  };
};
