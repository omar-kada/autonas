import type { ContainerStatusHealth } from '@/api/api';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';

function colorForStatus(status: ContainerStatusHealth): string {
  switch (status) {
    case 'healthy':
      return 'bg-green-400';
    case 'unhealthy':
      return 'bg-red-400';
    case 'starting':
      return 'bg-yellow-400';
    default:
      return '';
  }
}

export function StatusBadge(props: {
  status: ContainerStatusHealth;
  label: string;
  className?: string;
}) {
  return (
    <Badge className={cn(colorForStatus(props.status), props.className)}>
      {props.label ?? 'unknown'}
    </Badge>
  );
}
