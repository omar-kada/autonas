import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { useLocation, useNavigate } from 'react-router-dom';
import { NavbarElement } from './view';

const navigationElements = [
  { label: 'STATUS', path: '/' },
  { label: 'DEPLOYMENTS', path: '/deployments' },
  { label: 'LOGS', path: '/logs' },
  { label: 'CONFIG', path: '/config' },
];

function NavBarElementList() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const onNavigate = useCallback((path: string) => () => navigate(path), [navigate]);

  const location = useLocation();

  const isMatched = useCallback(
    (path: string) => {
      return location.pathname === path;
    },
    [location.pathname],
  );
  return (
    <>
      {navigationElements.map((element) => (
        <NavbarElement
          key={element.path}
          label={t(element.label)}
          navigate={onNavigate(element.path)}
          className={isMatched(element.path) ? '' : 'opacity-75'}
        />
      ))}
    </>
  );
}

export function Navbar() {
  return (
    <>
      {/* Top navigation bar, on big screens */}
      <nav className="hidden sm:flex items-center gap-6">
        <NavBarElementList />
      </nav>
      {/* Bottom navigation bar, on small screens */}
      <nav className="sm:hidden h-14 border-t w-full bg-background fixed flex items-center justify-around bottom-0 left-0 right-0 z-50">
        <NavBarElementList />
      </nav>
    </>
  );
}
