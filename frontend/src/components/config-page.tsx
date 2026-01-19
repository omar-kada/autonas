import { getConfigQueryOptions, getFeaturesQueryOptions } from '@/hooks';
import { useQuery } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { Skeleton } from './ui/skeleton';
import { ErrorAlert, InfoEmpty } from './view';

export function ConfigPage() {
  const { t } = useTranslation();
  const { data: features } = useQuery(getFeaturesQueryOptions());
  const {
    data: config,
    isPending,
    error,
  } = useQuery(
    getConfigQueryOptions({
      enabled: !!features?.displayConfig,
    }),
  );
  if (!features?.displayConfig) {
    return (
      <InfoEmpty
        title="CONFIGURATION.DISABLED_TITLE"
        details="CONFIGURATION.DISABLED_DESCRIPTION"
      ></InfoEmpty>
    );
  }
  return (
    <div className="p-4 space-y-4">
      <h2 className="text-2xl font-bold">{t('CONFIGURATION.CONFIGURATION')}</h2>
      <ErrorAlert title={error && 'ALERT.LOAD_CONFIGURATION_ERROR'} details={error?.message} />

      {isPending ? (
        <div className="flex flex-col gap-10">
          <EnvVarSkeleton repeat={2} />
          <EnvVarSkeleton repeat={1} />
        </div>
      ) : (
        <>
          {config?.globalVariables &&
            Object.entries(config.globalVariables).map(([key, value]) => (
              <div key={key}>
                <strong>{key}:</strong> {value}
              </div>
            ))}
          {config?.services &&
            Object.entries(config.services)?.map(([serviceKey, serviceValue]) => (
              <div key={serviceKey}>
                <strong>{serviceKey}:</strong>
                {serviceValue &&
                  Object.entries(serviceValue)?.map(([key, value]) => (
                    <div key={key} className="ml-4">
                      <strong>{key}:</strong> {value}
                    </div>
                  ))}
              </div>
            ))}
        </>
      )}
    </div>
  );
}

function EnvVarSkeleton({ repeat }: { repeat: number }) {
  return (
    <div className="flex flex-col gap-4">
      <Skeleton className="h-6 w-25"></Skeleton>
      {Array(repeat)
        .fill({})
        .map((_, index) => (
          <span className="flex gap-4" key={index}>
            <Skeleton className="h-6 w-50" />
            <Skeleton className="h-6 w-50" />
          </span>
        ))}
    </div>
  );
}
