import React from 'react';
import { useStatus } from './useStatus';
import dockerLogo from '@app/assets/docker.svg';
import { useTranslation } from 'react-i18next';

const StatusDisplay: React.FC = () => {
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
    <>
      <h2>{t('STATUS')}</h2>
      {Object.entries(data).map(([serviceName, serviceState]) => (
        <div key={serviceName}>
          <h3>{serviceName}</h3>
          <img
            src={`https://raw.githubusercontent.com/walkxcode/dashboard-icons/main/png/${serviceName}.png`}
            onError={(e: React.SyntheticEvent<HTMLImageElement, Event>) => {
              e.currentTarget.src = dockerLogo;
            }}
            className="logo"
          ></img>
          <ul>
            {serviceState.map((item, index) => (
              <li key={`${serviceName}-${index}`}>
                <strong>{item.Name}:</strong> {item.State || 'No state available'} &nbsp;
                {item.Health || 'No health available'}
              </li>
            ))}
          </ul>
        </div>
      ))}
    </>
  );
};

export default StatusDisplay;
