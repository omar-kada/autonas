import { getDiffQueryOptions, getStatsQueryOptions, getSyncOptions, useIsMobile } from '@/hooks';
import { cn, useDeploymentNavigate } from '@/lib';
import { useMutation, useQuery } from '@tanstack/react-query';
import { FileDiff, History, RefreshCcw, TriangleAlert } from 'lucide-react';
import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { toast } from 'sonner';
import { Badge } from '../ui/badge';
import { Button } from '../ui/button';
import { Skeleton } from '../ui/skeleton';
import { Spinner } from '../ui/spinner';
import { HumanTime } from '../view';
import { DeploymentDiffDialog } from './deployment-diff-dialog';
import { DeploymentStatusBadge } from './deployment-status-badge';

export function DeploymentToolbar({ className }: { className?: string }) {
  const { t } = useTranslation();
  const isMobile = useIsMobile();
  const depNavigate = useDeploymentNavigate();

  const { data: stats, isPending, error } = useQuery(getStatsQueryOptions());
  const { data: diffs, isFetching: isDiffsLoading } = useQuery(getDiffQueryOptions());

  const {
    mutateAsync: sync,
    isPending: isSyncLoading,
    error: syncError,
  } = useMutation(getSyncOptions());

  const handleSync = useCallback(() => {
    toast.promise(
      () =>
        sync().then((res) => {
          if (res.data?.id && res.data.id !== '0') {
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
  }, [sync, t, depNavigate]);

  if (error) {
    return <div>Error fetching stats: {error?.message}</div>;
  }

  return (
    <div className={cn('flex flex-wrap items-center align-bottom gap-4 m-2', className)}>
      <div className="flex items-center p-2 gap-2">
        <span className="text-sm font-light mx-1 flex-1 flex gap-1 items-center">
          <History className="size-4"></History>
          {t('LAST_X_DAYS', { days: 30 })} :
        </span>
        {isPending ? (
          <StatsSkeleton />
        ) : (
          <>
            <DeploymentStatusBadge
              status="success"
              label={String(stats?.success ?? 0)}
            ></DeploymentStatusBadge>
            {stats?.error ? (
              <DeploymentStatusBadge
                status="error"
                label={String(stats.error)}
              ></DeploymentStatusBadge>
            ) : null}
          </>
        )}
      </div>

      <div className="flex flex-row items-center gap-1 justify-end-safe flex-1">
        <span className="text-sm font-light text-muted-foreground mr-2">
          {syncError
            ? syncError.message
            : !isMobile && (
                <>
                  {t('AUTO_SYNC')} :&nbsp;
                  {isPending ? (
                    <Spinner className="inline"></Spinner>
                  ) : (
                    <HumanTime time={stats?.nextDeploy} defaultValue={t('DISABLED')}></HumanTime>
                  )}
                </>
              )}
        </span>
        <DeploymentDiffDialog>
          <Button variant="outline">
            <FileDiff />
            {!isMobile && t('DIFF')}
            {diffs != null && (
              <Badge
                className="h-5 min-w-5 rounded-full px-1 font-mono tabular-nums"
                variant={diffs.length > 0 ? 'default' : 'outline'}
              >
                {diffs.length}
              </Badge>
            )}
            {isDiffsLoading ? <Spinner></Spinner> : null}
          </Button>
        </DeploymentDiffDialog>
        <Button variant="outline" onClick={handleSync} disabled={isSyncLoading}>
          <RefreshCcw className={isSyncLoading ? 'animate-spin' : ''} />
          {!isMobile && t('SYNC_NOW')}
          {syncError ? <TriangleAlert className="text-destructive" /> : null}
        </Button>
      </div>
    </div>
  );
}

function StatsSkeleton() {
  return <Skeleton className="h-4 w-20"></Skeleton>;
}
