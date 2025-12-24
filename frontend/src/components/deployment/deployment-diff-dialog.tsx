import { useDiffAPIGet } from '@/api/api';
import { type ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import { DeploymentDiff } from '.';
import { Dialog, DialogContent, DialogTitle, DialogTrigger } from '../ui/dialog';
import { ScrollArea } from '../ui/scroll-area';

export function DeploymentDiffDialog({ children }: { children?: ReactNode }) {
  const { t } = useTranslation();
  const {
    data: diffs,
    isLoading,
    error,
  } = useDiffAPIGet({
    query: {
      refetchOnMount: true,
    },
  });

  return (
    <Dialog>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent className="w-full max-w-none sm:max-w-[90vw] lg:max-w-4xl">
        <DialogTitle>{t('DIFF_BETWEEN_DEPLOYED_AND_REMOTE')}</DialogTitle>
        <ScrollArea className="max-h-[80vh] max-w-[90vw] spa">
          {error ? (
            error.message
          ) : isLoading ? (
            'Loading '
          ) : (
            <DeploymentDiff fileDiffs={diffs?.data ?? []} autoOpen={true}></DeploymentDiff>
          )}
        </ScrollArea>
      </DialogContent>
    </Dialog>
  );
}
