import { Separator } from '@radix-ui/react-separator';
import { Layers, ScrollText, ServerCog, SlidersHorizontal } from 'lucide-react';
import { NavbarElement } from './view';

const navigationElements = [
  { label: 'STATUS', path: '/', icon: <Layers></Layers> },
  { label: 'DEPLOYMENTS', path: '/deployments', icon: <ServerCog></ServerCog> },
  { label: 'LOGS', path: '/logs', icon: <ScrollText></ScrollText> },
  { label: 'CONFIG', path: '/config', icon: <SlidersHorizontal></SlidersHorizontal> },
];

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
      {navigationElements.map((element, index) => (
        <>
          {index > 0 && <Separator orientation="vertical" className="bg-accent w-px h-8" />}
          <NavbarElement
            label={element.label}
            icon={element.icon}
            path={element.path}
            className={'flex-1'}
          />
        </>
      ))}
    </>
  );
}
