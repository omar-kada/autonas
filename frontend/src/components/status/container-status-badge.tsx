import type { ContainerStatusHealth } from '@/api/api';
import { Badge } from '@/components/ui/badge';
import { borderForStatus, iconForHealth, textColorForStatus } from '@/lib';
import { cn } from '@/lib/utils';

export function ContainerStatusBadge(props: {
  status: ContainerStatusHealth;
  label?: string;
  className?: string;
  iconOnly?: boolean;
}) {
  const Icon = iconForHealth(props.status);
  return (
    <Badge variant="outline" className={cn(borderForStatus(props.status), props.className)}>
      <Icon className={textColorForStatus(props.status)} />
      {!props.iconOnly && (props.label ?? props.status ?? 'unknown')}
    </Badge>
  );
}
