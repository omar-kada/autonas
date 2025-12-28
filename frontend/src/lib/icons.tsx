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

export function iconForStatus(status: DeploymentStatus) {
  switch (status) {
    case 'success':
      return Check;
    case 'error':
      return X;
    case 'running':
      return LoaderCircle;
    case 'planned':
      return Clock;
    default:
      return CircleQuestionMark;
  }
}

export function iconForHealth(health: ContainerStatusHealth) {
  switch (health) {
    case 'healthy':
      return Heart;
    case 'unhealthy':
      return HeartCrack;
    case 'starting':
      return HeartPulse;
    default:
      return HeartOff;
  }
}
