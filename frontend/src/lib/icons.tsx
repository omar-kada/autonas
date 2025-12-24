import type { ContainerStatusHealth, DeploymentStatus } from '@/api/api';
import {
  Check,
  CircleQuestionMark,
  Clock,
  Heart,
  HeartCrack,
  HeartOff,
  HeartPulse,
  LoaderCircle,
  X,
} from 'lucide-react';
import { textColorForStatus } from './colors';

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

export function iconForHealth(health: ContainerStatusHealth) {
  const color = textColorForStatus(health);
  switch (health) {
    case 'healthy':
      return <Heart className={color} />;
    case 'unhealthy':
      return <HeartCrack className={color} />;
    case 'starting':
      return <HeartPulse className={color} />;
    default:
      return <HeartOff className={color} />;
  }
}
