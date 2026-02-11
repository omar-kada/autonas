import {
  getRegisterAPIRegisteredQueryOptions,
  getRegisterAPIRegisterMutationOptions,
  getUserAPIGetQueryOptions,
  type Credentials,
} from '@/api/api';
import { useMutation } from '@tanstack/react-query';
import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { toast } from 'sonner';

export const getRegisterOptions = () => {
  return getRegisterAPIRegisterMutationOptions({
    mutation: {
      onSuccess: (data, _, __, context) => {
        if (data.data.success) {
          context.client.refetchQueries(getRegisterAPIRegisteredQueryOptions());
          context.client.refetchQueries(getUserAPIGetQueryOptions());
        }
      },
    },
  });
};

export const useRegister = () => {
  const { t } = useTranslation();

  const registerMutation = useMutation(getRegisterOptions());

  const handleRegister = useCallback(
    (user: Credentials) => {
      toast.promise(
        () =>
          registerMutation.mutateAsync({ data: user }).then((res) => {
            return res.data?.success;
          }),
        {
          error: t('ALERT.USER_CREATION_ERROR'),
        },
      );
    },
    [registerMutation.mutateAsync, t],
  );

  return {
    ...registerMutation,
    register: handleRegister,
  };
};
