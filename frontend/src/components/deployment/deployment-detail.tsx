import { useDeployment } from '@/hooks';
import { formatElapsed, ROUTES } from '@/lib';
import { Loader, Timer, User } from 'lucide-react';
import type { ReactElement } from 'react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import { DeploymentDiff, DeploymentEventLog, DeploymentStatusBadge } from '.';
import { Item, ItemContent, ItemMedia, ItemTitle } from '../ui/item';
import { ScrollArea } from '../ui/scroll-area';
import { HumanTime } from '../view';

export function DeploymentDetail({ id }: { id?: string }) {
  const { t } = useTranslation();
  const { deployment, error, isPending } = useDeployment(id);

  if (error != null) {
    return <div>{t('ERROR_WHILE_LOADING_DEPLOYMENT')}</div>;
  }

  if (isPending) {
    return <Loader></Loader>;
  }

  if (deployment == null) {
    return <div>{t('SELECT_DEPLOYMENT_FOR_DETAILS')}</div>;
  }

  return (
    <div className="flex flex-col h-full">
      <div className="text-2xl font-semibold m-4">
        <Link to={ROUTES.DEPLOYMENT(id)}>#{id} - </Link>
        {deployment.title}
        <DeploymentStatusBadge status={deployment.status} className="mx-3"></DeploymentStatusBadge>
        <HumanTime className="text-sm" time={deployment.time}></HumanTime>
      </div>
      <ScrollArea className="gap-4 h-1 flex-1">
        <div className="flex flex-col gap-4 m-4">
          <div className="grid grid-cols-2 gap-4 self-start">
            <InfoItem
              icon={<User className="size-5" />}
              label={deployment.author !== '' ? deployment.author : t('AUTOMATIC')}
            />
            <InfoItem
              icon={<Timer className="size-5" />}
              label={formatElapsed(deployment.time, deployment.endTime)}
            />
          </div>
          <DeploymentDiff fileDiffs={deployment.files ?? []} />
          <DeploymentEventLog events={deployment.events ?? []} />
        </div>
      </ScrollArea>
    </div>
  );
}

function InfoItem({ icon, label }: { icon: ReactElement; label: string }) {
  return (
    <Item variant="muted" size="sm">
      <ItemMedia>{icon}</ItemMedia>
      <ItemContent>
        <ItemTitle>{label}</ItemTitle>
      </ItemContent>
    </Item>
  );
}
