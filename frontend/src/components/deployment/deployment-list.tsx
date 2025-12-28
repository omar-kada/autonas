import type { Deployment } from '@/api/api';
import { ChevronRight } from 'lucide-react';
import { useCallback } from 'react';
import { Item, ItemActions, ItemContent, ItemDescription, ItemTitle } from '../ui/item';
import { HumanTime } from '../view';
import { DeploymentStatusBadge } from './deployment-status-badge';

export function DeploymentList({
  deployments,
  selectedDeployment,
  OnSelect,
}: {
  deployments: Array<Deployment>;
  selectedDeployment?: string;
  OnSelect: (item: Deployment) => void;
}) {
  const onDeploymentClick = useCallback(
    (deployment: Deployment) => () => OnSelect(deployment),
    [OnSelect],
  );

  return (
    <div className="space-y-2 p-3">
      {deployments.map((deployment) => (
        <DeploymentItem
          key={deployment.id}
          deployment={deployment}
          isSelected={deployment.id === selectedDeployment}
          onSelect={onDeploymentClick}
        ></DeploymentItem>
      ))}
    </div>
  );
}

function DeploymentItem({
  deployment,
  isSelected,
  onSelect,
}: {
  deployment: Deployment;
  isSelected: boolean;
  onSelect: (deployment: Deployment) => () => void;
}) {
  return (
    <Item
      key={deployment.id}
      className={`cursor-pointer ${isSelected ? 'bg-accent' : ''}`}
      onClick={onSelect(deployment)}
      variant="outline"
    >
      <ItemContent>
        <ItemTitle>
          <DeploymentStatusBadge status={deployment.status} iconOnly />
          {deployment.title}
        </ItemTitle>
        <ItemDescription className="text-xs">
          #{deployment.id} - <HumanTime time={deployment.time} />
        </ItemDescription>
      </ItemContent>
      <ItemActions className="flex-col justify-between h-full">
        <ChevronRight />
      </ItemActions>
    </Item>
  );
}
