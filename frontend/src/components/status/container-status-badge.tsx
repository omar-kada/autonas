import type { ContainerStatusHealth } from '@/api/api';
import { Badge } from '@/components/ui/badge';
import { borderForStatus, iconForHealth } from '@/lib';
import { cn } from '@/lib/utils';

export function ContainerStatusBadge(props: {
  status: ContainerStatusHealth;
  label?: string;
  className?: string;
  iconOnly?: string;
}) {
  return (
    <Badge variant="outline" className={cn(borderForStatus(props.status), props.className)}>
      {iconForHealth(props.status)} {!props.iconOnly && (props.label ?? props.status ?? 'unknown')}
    </Badge>
  );
}
