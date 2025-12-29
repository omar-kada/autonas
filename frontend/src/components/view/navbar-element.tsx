import { useIsMobile } from '@/hooks';
import { cn } from '@/lib/utils';
import type { LucideProps } from 'lucide-react';
import type { ComponentType } from 'react';
import { useTranslation } from 'react-i18next';
import { Link, useMatch } from 'react-router-dom';

export type NavbarElementProps = {
  label: string;
  path: string;
  Icon: ComponentType<LucideProps>;
  className?: string;
};

export function NavbarElement({ label, path, Icon, className }: NavbarElementProps) {
  const { t } = useTranslation();
  const isMobile = useIsMobile();
  const matched = useMatch({
    path,
    end: path === '/',
  });

  return (
    <Link
      className={cn(
        `flex flex-col md:flex-row text-sm font-medium md:gap-2 h-full justify-around items-center px-2 ${isMobile ? 'border-t-2' : 'border-b-2'} ${matched ? 'border-primary box-border' : 'opacity-75 border-accent'}`,
        className,
      )}
      to={path}
    >
      <Icon className="size-5 mt-1 md:size-4 md:mt-0" />
      <span className="inline-flex text-xs font-light md:text-sm md:font-bold">{t(label)}</span>
    </Link>
  );
}
