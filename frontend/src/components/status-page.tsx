import { useStatus } from '@/hooks';
import { useTranslation } from 'react-i18next';
import { ServiceStatus } from './view';

export function StatusPage() {
  const { t } = useTranslation();
  const { data, isLoading, error } = useStatus();

  if (isLoading) {
    return <div>Loading status...</div>;
  }

  if (error) {
    return <div>Error fetching status: {error.message}</div>;
  }

  // Check if data exists and is an object
  if (!data || typeof data !== 'object') {
    return <div>No status data available</div>;
  }

  return (
    <div className="p-4 space-y-4">
      <h2 className="text-2xl font-bold">{t('STATUS')}</h2>
      <div className="space-y-2">
        {data.map((stackStatus) => (
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
