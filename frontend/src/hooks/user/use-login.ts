import {
  getAuthAPILoginMutationOptions,
  getUserAPIGetQueryOptions,
  type Credentials,
} from '@/api/api';
import { useMutation } from '@tanstack/react-query';
import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { toast } from 'sonner';

export const getLoginOptions = () => {
  return getAuthAPILoginMutationOptions({
    mutation: {
      onSuccess: (data, _, __, context) => {
        if (data.data.success) {
          context.client.refetchQueries(getUserAPIGetQueryOptions());
        }
      },
    },
  });
};

export const useLogin = () => {
  const { t } = useTranslation();

  const loginMutation = useMutation(getLoginOptions());

  const handleLogin = useCallback(
    (user: Credentials) => {
      toast.promise(
        () =>
          loginMutation.mutateAsync({ data: user }).then((res) => {
            return res.data?.success;
          }),
        {
          error: t('ALERT.LOGIN_ERROR'),
        },
      );
    },
    [loginMutation.mutateAsync, t],
  );

  return {
    ...loginMutation,
    login: handleLogin,
  };
};
