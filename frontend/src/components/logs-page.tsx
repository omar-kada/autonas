import { useLogs } from '@/hooks';
import { useTranslation } from 'react-i18next';
import { LogEntries, LogFilterBar } from './logs';
import { useLogFilter as useFilterLogs } from './logs/use-filter-logs';
import { Skeleton } from './ui/skeleton';
import { ErrorAlert, HeaderLayout, InfoEmpty } from './view';

export function LogsPage() {
  const { t } = useTranslation();
  const { data: logs = [], isPending, error } = useLogs(50);
  const { text, setText, activeLevels, setActiveLevels, filtered } = useFilterLogs(logs);

  return (
    <HeaderLayout
      header={
        <div className="flex items-center justify-between gap-5 flex-wrap">
          <h2 className="text-2xl font-bold">{t('LOGS.LOGS')}</h2>
          <LogFilterBar
            text={text}
            onTextChange={setText}
            activeLevels={activeLevels}
            onLevelsChange={setActiveLevels}
          />
        </div>
      }
    >
      <div className="flex flex-col gap-2 min-h-0 flex-1 h-full">
        <ErrorAlert title="ALERT.LOAD_LOGS_ERROR" error={error} />

        {isPending ? (
          Array(10)
            .fill({})
            .map((_, index) => (
              <Skeleton className="mt-2 w-10/12 h-4" key={`status-skeleton-${index}`} />
            ))
        ) : logs && Object.keys(logs).length > 0 ? (
          <LogEntries logs={filtered} />
        ) : (
          <InfoEmpty title="LOGS.NO_LOGS" details="LOGS.NO_LOGS_DESCRIPTION"></InfoEmpty>
        )}
      </div>
    </HeaderLayout>
  );
}
