import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import type { HealthStatus } from '@/models/stack-status';

function colorForStatus(status: HealthStatus): string {
  switch (status) {
    case 'healthy':
      return 'bg-green-700';
    case 'unhealthy':
      return 'bg-red-700';
    case 'starting':
      return 'bg-yellow-700';
    default:
      return 'bg-gray-700';
  }
}

export function StatusBadge(props: { status: HealthStatus; label: string; className?: string }) {
  return (
    <Badge variant="secondary" className={cn(colorForStatus(props.status), props.className)}>
      {props.label}
    </Badge>
  );
}
