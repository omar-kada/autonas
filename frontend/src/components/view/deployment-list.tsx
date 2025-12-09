import type { Deployment, DeploymentStatus } from '@/api/api';
import { Check, ChevronRight, CircleQuestionMark, Clock, LoaderCircle, X } from 'lucide-react';
import { useCallback } from 'react';
import { Badge } from '../ui/badge';
import { Item, ItemActions, ItemContent, ItemDescription, ItemTitle } from '../ui/item';
import { HumanTime } from './human-time';
function colorForStatus(status: DeploymentStatus): string {
  switch (status) {
    case 'success':
      return 'bg-green-400';
    case 'error':
      return 'bg-red-400';
    case 'running':
      return 'bg-blue-400';
    case 'planned':
      return 'bg-slate-400';
    default:
      return '';
  }
}

function iconForStatus(status: DeploymentStatus) {
  switch (status) {
    case 'success':
      return <Check className="h-4 w-4" />;
    case 'error':
      return <X className="h-4 w-4" />;
    case 'running':
      return <LoaderCircle className="h-4 w-4 animate-spin" />;

    case 'planned':
      return <Clock className="h-4 w-4" />;
    default:
      return <CircleQuestionMark className="h-4 w-4"></CircleQuestionMark>;
  }
}
export function DeploymentList(props: {
  deployments: Deployment[];
  OnSelect: (item: Deployment) => void;
}) {
  const onDeploymentClick = useCallback(
    (deployment: Deployment) => () => props.OnSelect(deployment),
    [props],
  );

  return (
    <div className="space-y-2">
      {props.deployments.map((deployment) => DeploymentItem(deployment, onDeploymentClick))}
    </div>
  );
}
function DeploymentItem(
  deployment: Deployment,
  onDeploymentClick: (deployment: Deployment) => () => void,
) {
  return (
    <Item
      key={deployment.id}
      className="cursor-pointer"
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
