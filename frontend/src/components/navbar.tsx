import { ROUTES } from '@/lib';
import { Layers, ScrollText, ServerCog, SlidersHorizontal } from 'lucide-react';
import { NavbarElement, type NavbarElementProps } from './view';

const LINKS: Array<NavbarElementProps> = [
  {
    label: 'DEPLOYMENTS.DEPLOYMENTS',
    Icon: ServerCog,
    path: ROUTES.DEPLOYMENTS,
  },
  {
    label: 'STATUS.STATUS',
    Icon: Layers,
    path: ROUTES.STATUS,
  },
  {
    label: 'LOGS',
    Icon: ScrollText,
    path: ROUTES.LOGS,
  },
  {
    label: 'CONFIGURATION.CONFIGURATION',
    Icon: SlidersHorizontal,
    path: ROUTES.CONFIG,
  },
];

export function NavBar({ className }: { className?: string }) {
  return (
    <nav className={className}>
      {LINKS.map((link) => (
        <NavbarElement key={link.label} {...link} className={'flex-1'} />
      ))}
    </nav>
  );
}
