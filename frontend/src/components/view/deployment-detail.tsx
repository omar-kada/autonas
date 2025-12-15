import { useDeployment } from '@/hooks';
import { Loader } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { DeploymentDiff } from './deployment-diff';
import { DeploymentEventLog } from './deployment-event-log';

export function DeploymentDetail({ id }: { id?: string }) {
  const { t } = useTranslation();

  const { deployment, error, isLoading } = useDeployment(id);

  if (error != null) {
    return <div>{t('ERROR_WHILE_LOADING_DEPLOYMENT')}</div>;
  }

  if (isLoading) {
    return <Loader></Loader>;
  }

  if (deployment == null) {
    return <div>{t('SELECT_DEPLOYMENT_FOR_DETAILS')}</div>;
  }

  return (
    <div className=" flex flex-col gap-4">
      <h2 className="text-2xl font-semibold mb-4">{deployment.title}</h2>
      <DeploymentDiff fileDiffs={deployment.files ?? []} />
      <DeploymentEventLog events={deployment.events ?? []} />
    </div>
  );
}
