import { Separator } from '@radix-ui/react-separator';
import { Layers, ScrollText, ServerCog, SlidersHorizontal } from 'lucide-react';
import { NavbarElement } from './view';

export function Navbar() {
  return (
    <>
      {/* Top navigation bar, on big screens */}
      <nav className="hidden sm:flex bg-sidebar h-14 items-center">
        <NavBarElementList />
      </nav>
      {/* Bottom navigation bar, on small screens */}
      <nav className="flex sm:hidden bg-sidebar py-2 h-14 border-t w-full fixed items-center justify-around bottom-0 left-0 right-0 z-50">
        <NavBarElementList />
      </nav>
    </>
  );
}

function NavBarElementList() {
  return (
    <>
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
    </>
  );
}
