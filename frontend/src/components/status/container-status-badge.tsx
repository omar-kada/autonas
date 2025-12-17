import type { ContainerStatusHealth } from '@/api/api';
import { Badge } from '@/components/ui/badge';
import { colorForStatus } from '@/lib';
import { cn } from '@/lib/utils';

export function ContainerStatusBadge(props: {
  status: ContainerStatusHealth;
  label?: string;
  className?: string;
}) {
  return (
    <Badge className={cn(colorForStatus(props.status), props.className)}>
      {props.label ?? props.status ?? 'unknown'}
    </Badge>
  );
}
