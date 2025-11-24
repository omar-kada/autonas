import { useCallback } from 'react';
import { useStatus } from './useStatus';
import dockerLogo from '@app/assets/docker.svg';
import { useTranslation } from 'react-i18next';

function StatusDisplay() {
  const { t } = useTranslation();
  const { data, isLoading, error } = useStatus();

  const defatulIconOnErrorCallback = useCallback(
    (e: React.SyntheticEvent<HTMLImageElement, Event>) => {
      e.currentTarget.src = dockerLogo;
    },
    [],
  );
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
    <>
      <h2>{t('STATUS')}</h2>
      {Object.entries(data).map(([serviceName, serviceState]) => (
        <div key={serviceName}>
          <h3>{serviceName}</h3>
          <img
            src={`https://raw.githubusercontent.com/walkxcode/dashboard-icons/main/png/${serviceName}.png`}
            onError={defatulIconOnErrorCallback}
            alt={'service logo'}
            className="logo"
          ></img>
          <ul>
            {serviceState.map((item) => (
              <li key={`${serviceName}-${item.Name}`}>
                <strong>{item.Name}:</strong> {item.State || 'No state available'} &nbsp;
                {item.Health || 'No health available'}
              </li>
            ))}
          </ul>
        </div>
      ))}
    </>
  );
}

export default StatusDisplay;
