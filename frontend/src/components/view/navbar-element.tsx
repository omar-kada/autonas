import { cn } from '@/lib/utils';
import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { useMatch, useNavigate } from 'react-router-dom';
import { Button } from '../ui/button';

export type NavbarElementProps = {
  label: string;
  path: string;
  icon?: React.ReactNode;
  className?: string;
};

export function NavbarElement({ label, path, icon, className }: NavbarElementProps) {
  const { t } = useTranslation();
  const matched = useMatch({
    path,
    end: path === '/',
  });
  const navigate = useNavigate();
  const onNavigate = useCallback((path: string) => () => navigate(path), [navigate]);

  return (
    <Button
      variant="ghost"
      className={cn(
        `text-sm font-medium w-full rounded-none px-4 ${matched ? 'bg-accent' : ''}`,
        className,
      )}
      onClick={onNavigate(path)}
    >
      {icon}
      <span className="hidden md:inline-flex">{t(label)}</span>
    </Button>
  );
}
