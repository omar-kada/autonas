import { getDiffQueryOptions } from '@/hooks';
import { useQuery } from '@tanstack/react-query';
import { type ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import { DeploymentDiff, DeploymentDiffSkeleton } from '.';
import { Dialog, DialogContent, DialogDescription, DialogTitle, DialogTrigger } from '../ui/dialog';
import { ScrollArea } from '../ui/scroll-area';
import { ErrorAlert } from '../view';

export function DeploymentDiffDialog({ children }: { children?: ReactNode }) {
  const { t } = useTranslation();
  const {
    data: diffs,
    isPending,
    error,
  } = useQuery({ ...getDiffQueryOptions(), refetchOnMount: true });

  return (
    <Dialog>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent className="w-full max-w-none sm:max-w-[90vw] lg:max-w-4xl">
        <DialogTitle>{t('DIFF.DIFF_BETWEEN_DEPLOYED_AND_REMOTE')}</DialogTitle>
        <DialogDescription>
          {t('DIFF.DIFF_BETWEEN_DEPLOYED_AND_REMOTE_DESCRIPTION')}
        </DialogDescription>
        <ScrollArea className="max-h-[80vh] max-w-[90vw]">
          {error ? (
            <ErrorAlert title={error && 'ALERT.DIFF_ERROR'} details={error?.message} />
          ) : isPending ? (
            <DeploymentDiffSkeleton />
          ) : (
            <DeploymentDiff fileDiffs={diffs ?? []} autoOpen></DeploymentDiff>
          )}
        </ScrollArea>
      </DialogContent>
    </Dialog>
  );
}
