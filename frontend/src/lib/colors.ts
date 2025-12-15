import type { DeploymentStatus, EventLevel } from '@/api/api';

export function colorForStatus(status: DeploymentStatus): string {
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

export function logColor(level: EventLevel): string {
  switch (level) {
    case 'ERROR':
      return 'text-red-700 dark:text-red-300 ';
    case 'WARN':
      return 'text-yellow-700 dark:text-yellow-300';
    case 'DEBUG':
      return 'text-gray-700 dark:text-gray-300';
    case 'INFO':
    default:
      return '';
  }
}
