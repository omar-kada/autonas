import { useDeployementAPISync } from '@/api/api';
import { refetchDeployments, useDiff, useStats } from '@/hooks';
import { cn, useDeploymentNavigate } from '@/lib';
import { useQueryClient } from '@tanstack/react-query';
import { FileDiff, RefreshCcw, TriangleAlert } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { toast } from 'sonner';
import { DeploymentDiffDialog } from './deployment/deployment-diff-dialog';
import { DeploymentStatusBadge } from './deployment/deployment-status-badge';
import { ContainerStatusBadge } from './status';
import { Badge } from './ui/badge';
import { Button } from './ui/button';
import { HumanTime } from './view';

export function EnvironementStats({ className }: { className?: string }) {
  const { t } = useTranslation();
  const queryClient = useQueryClient();

  const depNavigate = useDeploymentNavigate();
  const { data: stats, isLoading, error } = useStats(30);
  const { data: diffs } = useDiff();

  const { mutateAsync: sync, isPending: isSyncLoading, error: syncError } = useDeployementAPISync();

  if (isLoading) {
    return <div>Loading stats...</div>;
  }

  if (error || stats == null) {
    return <div>Error fetching stats: {error?.message}</div>;
  }

  function handleSync() {
    toast.promise(
      () =>
        sync().then((res) => {
          if (res.data?.id && res.data.id != '0') {
            refetchDeployments(queryClient);
            depNavigate(res.data.id);
            return true;
          } else {
            return false;
          }
        }),
      {
        loading: t('SYNCHRONIZING'),
        success: (synced) => t(synced ? 'SYNC_SUCCESS' : 'SYNC_NO_CHANGES'),
        error: t('SYNC_ERROR'),
      },
    );
  }

  return (
    <div className={cn('flex flex-wrap items-center align-bottom gap-4 m-4', className)}>
      <span>
        {t('Deployment')} : <DeploymentStatusBadge status={stats.status}></DeploymentStatusBadge>
      </span>
      <span>
        {t('Containers')} : <ContainerStatusBadge status={stats.health}></ContainerStatusBadge>
      </span>
      <div className="flex flex-row items-center gap-1 justify-end-safe flex-1">
        <span className="text-sm font-light text-muted-foreground">
          {syncError ? (
            syncError.message
          ) : (
            <>
              {t('AUTO_SYNC')} : <HumanTime time={stats.nextDeploy}></HumanTime>
            </>
          )}
        </span>
        <DeploymentDiffDialog>
          <Button variant="outline">
            <FileDiff />
            {t('DIFF')}
            {diffs && (
              <Badge
                className="h-5 min-w-5 rounded-full px-1 font-mono tabular-nums"
                variant={diffs.length > 0 ? 'default' : 'outline'}
              >
                {diffs.length}
              </Badge>
            )}
          </Button>
        </DeploymentDiffDialog>
        <Button variant="outline" onClick={handleSync} disabled={isSyncLoading}>
          <RefreshCcw className={isSyncLoading ? 'animate-spin' : ''} />
          {t('SYNC_NOW')}
          {syncError ? <TriangleAlert className="text-destructive" /> : null}
        </Button>
      </div>
    </div>
  );
}
