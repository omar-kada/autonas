import { getConfigAPIGetQueryKey, getConfigAPISetMutationOptions, type Config } from '@/api/api';
import { useMutation } from '@tanstack/react-query';
import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { toast } from 'sonner';

export const getUpdateConfigOptions = () => {
  return getConfigAPISetMutationOptions({
    mutation: {
      onSuccess: (data, _, __, context) => {
        context.client.setQueryData(getConfigAPIGetQueryKey(), data);
      },
    },
  });
};

export const useUpdateConfig = () => {
  const { t } = useTranslation();

  const updateMutation = useMutation(getUpdateConfigOptions());

  const handleSync = useCallback(
    (config: Config) => {
      toast.promise(() => updateMutation.mutateAsync({ data: config }), {
        success: t('ALERT.SUCCESS'),
        error: t('ALERT.ERROR'),
      });
    },
    [updateMutation.mutateAsync, t],
  );

  return {
    ...updateMutation,
    updateConfig: handleSync,
  };
};
