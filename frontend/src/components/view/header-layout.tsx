import { type ReactNode } from 'react';
import { Separator } from '../ui/separator';

export function HeaderLayout({ header, children }: { header: ReactNode; children: ReactNode }) {
  return (
    <>
      {header && <div className="p-2 px-4">{header}</div>}
      {header && <Separator orientation="horizontal" />}
      <div className="flex flex-col flex-1 overflow-hidden">{children}</div>
    </>
  );
}
