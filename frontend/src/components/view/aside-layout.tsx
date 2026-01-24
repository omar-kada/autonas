import { type ReactNode } from 'react';
import { HeaderLayout } from './header-layout';

export function AsideLayout({
  header,
  aside,
  children,
  focusMain,
}: {
  header: ReactNode;
  aside?: ReactNode;
  children: ReactNode;
  focusMain?: boolean;
}) {
  return (
    <HeaderLayout header={header}>
      {/* Sidebar (hidden on mobile if main is focused) */}
      {aside && (
        <aside
          className={`w-full h-full max-h-full flex flex-col sm:w-75 sm:shrink-0 pb-4 ${focusMain ? 'hidden sm:flex' : ''}`}
        >
          {aside}
        </aside>
      )}

      {/* Main content (hidden on mobile until main is focused) */}
      <main className={`${!focusMain ? 'hidden sm:block sm:flex-1' : 'w-full'}`}>{children}</main>
    </HeaderLayout>
  );
}
