import type { DeploymentStatus } from '@/api/api';
import { Badge } from '@/components/ui/badge';
import { colorForStatus, iconForStatus } from '@/lib';
import { cn } from '@/lib/utils';

export function DeploymentStatusBadge(props: {
  status: DeploymentStatus;
  iconOnly?: boolean;
  label?: string;
  className?: string;
}) {
  return (
    <Badge className={cn(colorForStatus(props.status), props.className)}>
      {iconForStatus(props.status)} {!props.iconOnly && (props.label ?? props.status ?? 'unknown')}
    </Badge>
  );
}
