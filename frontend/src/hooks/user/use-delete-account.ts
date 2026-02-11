import {
  getRegisterAPIRegisteredQueryOptions,
  getUserAPIDeleteMutationOptions,
  getUserAPIGetQueryOptions,
} from '@/api/api';
import { useMutation } from '@tanstack/react-query';
import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { toast } from 'sonner';

export const getDeleteAccountOptions = () => {
  return getUserAPIDeleteMutationOptions({
    mutation: {
      onSuccess: (data, _, __, context) => {
        if (data.data.success) {
          context.client.refetchQueries(getUserAPIGetQueryOptions());
          context.client.refetchQueries(getRegisterAPIRegisteredQueryOptions());
        }
      },
    },
  });
};

export const useDeleteAccount = () => {
  const { t } = useTranslation();

  const deleteAccountMutation = useMutation(getDeleteAccountOptions());

  const handledeleteAccount = useCallback(() => {
    return toast
      .promise(
        () =>
          deleteAccountMutation.mutateAsync().then((res) => {
            return res.data?.success;
          }),
        {
          error: t('ALERT.DELETE_ACCOUNT_ERROR'),
        },
      )
      .unwrap();
  }, [deleteAccountMutation.mutateAsync, t]);

  return {
    ...deleteAccountMutation,
    deleteAccount: handledeleteAccount,
  };
};
