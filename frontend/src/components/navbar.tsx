import { Separator } from '@radix-ui/react-separator';
import { Layers, ScrollText, ServerCog, SlidersHorizontal } from 'lucide-react';
import { NavbarElement } from './view';

export function NavBar({ className }: { className?: string }) {
  return (
    <nav className={className}>
      <NavbarElement label="STATUS" icon={<Layers />} path={'/'} className={'flex-1'} />
      <Separator orientation="vertical" className="bg-accent w-px h-8" />
      <NavbarElement
        label="DEPLOYMENTS"
        icon={<ServerCog />}
        path={'/deployments'}
        className={'flex-1'}
      />
      <Separator orientation="vertical" className="bg-accent w-px h-8" />
      <NavbarElement label="LOGS" icon={<ScrollText />} path={'/logs'} className={'flex-1'} />
      <Separator orientation="vertical" className="bg-accent w-px h-8" />
      <NavbarElement
        label="CONFIG"
        icon={<SlidersHorizontal />}
        path={'/config'}
        className={'flex-1'}
      />
    </nav>
  );
}
