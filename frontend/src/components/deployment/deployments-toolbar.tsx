import { getDiffQueryOptions, getStatsQueryOptions, useIsMobile, useSync } from '@/hooks';
import { cn } from '@/lib';
import { useQuery } from '@tanstack/react-query';
import { AlertCircleIcon, FileDiff, History, RefreshCcw, TriangleAlert } from 'lucide-react';
import { useTranslation } from 'react-i18next';
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
  const { sync, error: syncError, isPending: isSyncLoading } = useSync();

  const { data: stats, isPending, error } = useQuery(getStatsQueryOptions());
  const {
    data: diffs,
    isFetching: isDiffsLoading,
    error: diffError,
  } = useQuery(getDiffQueryOptions());

  return (
    <div className={cn('flex flex-wrap items-center align-bottom gap-4', className)}>
      <div className="flex items-center p-2 gap-2">
        <span className="text-sm font-light mx-1 flex-1 flex gap-1 items-center">
          <History className="size-4"></History>
          {t('TIME.LAST_X_DAYS', { days: 30 })} :
        </span>
        {error && (
          <>
            <AlertCircleIcon className="size-4 text-destructive" />
            <span className="text-sm text-destructive">{t('ALERT.LOAD_STATS_ERROR')}</span>
          </>
        )}
        {isPending ? (
          <StatsSkeleton />
        ) : (
          stats && (
            <>
              <DeploymentStatusBadge
                status="success"
                label={String(stats.success)}
              ></DeploymentStatusBadge>
              {stats.error ? (
                <DeploymentStatusBadge
                  status="error"
                  label={String(stats.error)}
                ></DeploymentStatusBadge>
              ) : null}
            </>
          )
        )}
      </div>

      <div className="flex flex-row items-center gap-1 justify-end-safe flex-1">
        <span className="text-sm font-light text-muted-foreground mr-2">
          {syncError
            ? syncError.message
            : !isMobile && (
                <>
                  {t('DEPLOYMENTS.AUTO_SYNC')} :&nbsp;
                  {error ? (
                    <AlertCircleIcon className="size-4 text-destructive inline" />
                  ) : isPending ? (
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
            {!isMobile && t('DIFF.DIFF')}
            {diffError ? (
              <AlertCircleIcon className="size-4 text-destructive inline" />
            ) : isDiffsLoading ? (
              <Spinner></Spinner>
            ) : (
              diffs != null && (
                <Badge
                  className="h-5 min-w-5 rounded-full px-1 font-mono tabular-nums"
                  variant={diffs.length > 0 ? 'default' : 'outline'}
                >
                  {diffs.length}
                </Badge>
              )
            )}
          </Button>
        </DeploymentDiffDialog>
        <Button variant="outline" onClick={sync} disabled={isSyncLoading}>
          {isSyncLoading ? <Spinner /> : <RefreshCcw />}
          {!isMobile && t('ACTION.SYNC_NOW')}
          {syncError ? <TriangleAlert className="text-destructive" /> : null}
        </Button>
      </div>
    </div>
  );
}

function StatsSkeleton() {
  return <Skeleton className="h-4 w-20"></Skeleton>;
}
