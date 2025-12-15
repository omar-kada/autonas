import type { Deployment } from '@/api/api';
import { colorForStatus, iconForStatus } from '@/lib';
import { ChevronRight } from 'lucide-react';
import { useCallback } from 'react';
import { Badge } from '../ui/badge';
import { Item, ItemActions, ItemContent, ItemDescription, ItemTitle } from '../ui/item';
import { HumanTime } from './human-time';

export function DeploymentList({
  deployments,
  selectedDeployment,
  OnSelect,
}: {
  deployments: Deployment[];
  selectedDeployment?: string;
  OnSelect: (item: Deployment) => void;
}) {
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
