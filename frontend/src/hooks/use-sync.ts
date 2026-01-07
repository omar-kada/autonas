import { getDeployementAPISyncMutationOptions } from '@/api/api';
import { useDeploymentNavigate } from '@/lib';
import { useMutation } from '@tanstack/react-query';
import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { toast } from 'sonner';
import { getDeploymentsQueryOptions } from './use-deployments';

export const getSyncOptions = () => {
  return getDeployementAPISyncMutationOptions({
    mutation: {
      onSuccess: (data, _, __, context) => {
        if (data.data?.id && data.data?.id !== '0') {
          context.client.refetchQueries(getDeploymentsQueryOptions());
        }
      },
    },
  });
};

export const useSync = (navigateOnSuccess = true) => {
  const depNavigate = useDeploymentNavigate();
  const { t } = useTranslation();

  const syncMutation = useMutation(getSyncOptions());

  const handleSync = useCallback(() => {
    toast.promise(
      () =>
        syncMutation.mutateAsync().then((res) => {
          if (res.data?.id && res.data.id !== '0') {
            if (navigateOnSuccess) {
              depNavigate(res.data.id);
            }
            return true;
          } else {
            return false;
          }
        }),
      {
        loading: t('ALERT.SYNCHRONIZING'),
        success: (synced) => t(synced ? 'ALERT.SYNC_SUCCESS' : 'ALERT.SYNC_NO_CHANGES'),
        error: t('ALERT.SYNC_ERROR'),
      },
    );
  }, [syncMutation.mutateAsync, t, depNavigate, navigateOnSuccess]);

  return {
    ...syncMutation,
    sync: handleSync,
  };
};
