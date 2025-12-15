import type { DeploymentStatus } from '@/api/api';
import { Check, CircleQuestionMark, Clock, LoaderCircle, X } from 'lucide-react';

export function iconForStatus(status: DeploymentStatus) {
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
