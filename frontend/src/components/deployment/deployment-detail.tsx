import { useDeployment } from '@/hooks';
import { formatElapsed, ROUTES } from '@/lib';
import { Loader, Timer, User } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import { DeploymentDiff, DeploymentEventLog, DeploymentStatusBadge } from '.';
import { Item, ItemContent, ItemMedia, ItemTitle } from '../ui/item';
import { HumanTime } from '../view';

export function DeploymentDetail({ id }: { id: string }) {
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
    <>
      <div className="text-2xl font-semibold mb-4">
        <Link to={ROUTES.DEPLOYMENT(id)}>#{id} - </Link>
        {deployment.title}
        <DeploymentStatusBadge status={deployment.status} className="mx-3"></DeploymentStatusBadge>
        <HumanTime time={deployment.time}></HumanTime>
      </div>

      <div className=" flex flex-col gap-4">
        <div className="grid grid-cols-2 gap-4 self-start">
          <Item variant="muted" size="sm">
            <ItemMedia>
              <User className="size-5" />
            </ItemMedia>
            <ItemContent>
              <ItemTitle>{deployment.author !== '' ? deployment.author : t('AUTOMATIC')}</ItemTitle>
            </ItemContent>
          </Item>
          <Item variant="muted" size="sm">
            <ItemMedia>
              <Timer className="size-5" />
            </ItemMedia>
            <ItemContent>
              <ItemTitle> {formatElapsed(deployment.time, deployment.endTime)}</ItemTitle>
            </ItemContent>
          </Item>
        </div>
        <DeploymentDiff fileDiffs={deployment.files ?? []} />
        <DeploymentEventLog events={deployment.events ?? []} />
      </div>
    </>
  );
}
