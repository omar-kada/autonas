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
        `flex text-sm font-medium gap-2 h-full justify-around items-center px-4 ${matched ? `${isMobile ? 'border-t-2' : 'border-b-2'} border-primary box-border` : 'opacity-75'}`,
        className,
      )}
      to={path}
    >
      <Icon className="size-5" />
      {!isMobile && <span className="hidden md:inline-flex">{t(label)}</span>}
    </Link>
  );
}
