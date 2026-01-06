import { getStatusQueryOptions } from '@/hooks';
import { useQuery } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { ServiceStatus, ServiceStatusSkeleton } from './status';

export function StatusPage() {
  const { t } = useTranslation();
  const { data, isPending, error } = useQuery(getStatusQueryOptions());

  if (error) {
    return <div>Error fetching status: {error.message}</div>;
  }

  // Check if data exists and is an object
  if (!data && !isPending) {
    return <div>No status data available</div>;
  }

  return (
    <div className="p-4 space-y-4">
      <h2 className="text-2xl font-bold">{t('STATUS')}</h2>
      <div className="space-y-2">
        {isPending
          ? Array(3)
              .fill({})
              .map((_, index) => <ServiceStatusSkeleton key={index} />)
          : data.map((stackStatus) => (
              <div key={stackStatus.stackId}>
                <ServiceStatus
                  serviceName={stackStatus.name}
                  serviceContainers={stackStatus.services}
                />
              </div>
            ))}
      </div>
    </div>
  );
}
