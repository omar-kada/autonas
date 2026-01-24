import dockerLogo from '@/assets/docker.svg';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { useMemo } from 'react';
import { cn } from './utils';

const logoPrefix = 'https://raw.githubusercontent.com/walkxcode/dashboard-icons/main/png/';

export function serviceLogoURL(service: string): string {
  return `${logoPrefix}${service}.png`;
}

export function ServiceLogo({ service, className }: { service: string; className?: string }) {
  const serviceLogo = useMemo(() => serviceLogoURL(service), [service]);

  return (
    <Avatar className={cn('size-10 rounded-none', className)}>
      <AvatarImage src={serviceLogo} />
      <AvatarFallback>
        <img src={dockerLogo} alt={`${service} logo`} />
      </AvatarFallback>
    </Avatar>
  );
}
