import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import type { HealthStatus } from '@/models/stack-status';

function colorForStatus(status: HealthStatus): string {
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

export function StatusBadge(props: { status: HealthStatus; label: string; className?: string }) {
  return <Badge className={cn(colorForStatus(props.status), props.className)}>{props.label}</Badge>;
}
