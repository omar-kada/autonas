import type { Deployment } from '@/api/api';
import { ChevronRight } from 'lucide-react';
import { useCallback } from 'react';
import { Item, ItemActions, ItemContent, ItemDescription, ItemTitle } from '../ui/item';
import { HumanTime } from '../view';
import { DeploymentStatusBadge } from './deployment-status-badge';

export function DeploymentListItem({
  deployment,
  isSelected,
  onSelect,
}: {
  deployment: Deployment;
  isSelected: boolean;
  onSelect: (deployment: Deployment) => void;
}) {
  const onDeploymentClick = useCallback(
    (deployment: Deployment) => () => onSelect(deployment),
    [onSelect],
  );

  return (
    <Item
      key={deployment.id}
      className={`cursor-pointer ${isSelected ? 'bg-accent' : ''}`}
      onClick={onDeploymentClick(deployment)}
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
