import type { Deployment } from '@/api/api';
import { useDeployments } from '@/hooks';
import { colorForStatus, iconForStatus } from '@/lib';
import { ChevronRight } from 'lucide-react';
import { useCallback } from 'react';
import { useParams } from 'react-router-dom';
import { Badge } from '../ui/badge';
import { Item, ItemActions, ItemContent, ItemDescription, ItemTitle } from '../ui/item';
import { HumanTime } from '../view';

export function DeploymentList({ OnSelect }: { OnSelect: (item: Deployment) => void }) {
  const { deployments } = useDeployments();
  const { id: selectedDeployment } = useParams();

  const onDeploymentClick = useCallback(
    (deployment: Deployment) => () => OnSelect(deployment),
    [OnSelect],
  );

  return (
    <div className="space-y-2">
      {deployments.map((deployment) =>
        DeploymentItem(deployment, deployment.id === selectedDeployment, onDeploymentClick),
      )}
    </div>
  );
}

function DeploymentItem(
  deployment: Deployment,
  isSelected: boolean,
  onDeploymentClick: (deployment: Deployment) => () => void,
) {
  return (
    <Item
      key={deployment.id}
      className={`cursor-pointer ${isSelected ? 'bg-accent' : ''}`}
      onClick={onDeploymentClick(deployment)}
      variant="outline"
    >
      <ItemContent>
        <ItemTitle>
          <Badge className={colorForStatus(deployment.status)}>
            {iconForStatus(deployment.status)}
          </Badge>
          {deployment.title}
        </ItemTitle>
        <ItemDescription className="text-xs">
          <HumanTime time={deployment.time} />
        </ItemDescription>
      </ItemContent>
      <ItemActions className="flex-col justify-between h-full">
        <ChevronRight />
      </ItemActions>
    </Item>
  );
}
