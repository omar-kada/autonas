import { useStats } from '@/hooks';
import { cn } from '@/lib';
import { FileDiff, History, RefreshCcw } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { Button } from '../ui/button';
import { HumanTime } from '../view';
import { DeploymentStatusBadge } from './deployment-status-badge';

export function DeploymentStats({ className }: { className?: string }) {
  const { t } = useTranslation();
  const { data: stats, isLoading, error } = useStats(30);
  if (isLoading) {
    return <div>Loading stats...</div>;
  }

  if (error || stats == null) {
    return <div>Error fetching stats: {error?.message}</div>;
  }

  return (
    <div className={cn('flex flex-wrap gap-4 m-4', className)}>
      <div className="flex flex-col flex-wrap text-muted-foreground self-start justify-between gap-1">
        <h2 className="text-2xl font-bold text-primary">{t('DEPLOYMENTS')}</h2>
        <div className="flex items-center ">
          <History className="size-4"></History>
          <span className="text-sm font-light mx-1">{t('LAST_X_DAYS', { days: 7 })} :</span>
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
        </div>
      </div>

      <div className="flex flex-col items-end justify-end-safe gap-1 justify-self-end-safe grow">
        <span className="text-sm font-light text-muted-foreground">
          {t('AUTO_SYNC')} : <HumanTime time={stats.nextDeploy}></HumanTime>
        </span>
        <div className="flex flex-row items-center gap-1">
          <Button variant="outline">
            <FileDiff />
            {t('DIFF')}
          </Button>
          <Button variant="outline">
            <RefreshCcw />
            {t('SYNC_NOW')}
          </Button>
        </div>
      </div>
    </div>
  );
}
