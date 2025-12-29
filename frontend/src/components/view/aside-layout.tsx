import { type ReactNode } from 'react';
import { Separator } from '../ui/separator';

export function AsideLayout({
  header,
  aside,
  children,
  focusMain,
}: {
  header: ReactNode;
  aside: ReactNode;
  children: ReactNode;
  focusMain: boolean;
}) {
  return (
    <>
      {header}
      {header && <Separator orientation="horizontal" />}
      <div className="flex flex-1 overflow-hidden">
        {/* Sidebar (hidden on mobile if main is focused) */}
        <aside
          className={`w-full h-full max-h-full flex flex-col sm:w-75 sm:shrink-0 m-2 pb-4 ${focusMain ? 'hidden sm:flex' : ''}`}
        >
          {aside}
        </aside>

        {/* Main content (hidden on mobile until main is focused) */}
        <main className={`flex-col ${!focusMain ? 'hidden sm:block sm:flex-1' : 'w-full'}`}>
          {children}
        </main>
      </div>
    </>
  );
}
