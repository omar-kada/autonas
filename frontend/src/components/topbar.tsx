import { Bell, Settings } from 'lucide-react';
import type { ReactNode } from 'react';
import { Button } from './ui/button';
import { ThemeToggle } from './view/theme-toggle';

export function Topbar({ children }: { children?: ReactNode }) {
  return (
    <header className="h-14 min-h-14 border-b w-full flex items-center justify-between px-4 bg-sidebar sticky top-0 z-50">
      {/* Logo */}
      <div className="text-xl font-semibold mr-5">AirCompose</div>
      <div className="flex-1 w-1 max-w-10">{/*gap*/}</div>
      <div className="flex flex-2 justify-between">{children}</div>

      {/* Settings + Notifications */}
      <div className="flex items-center gap-2">
        <Button variant="ghost" size="icon">
          <Settings className="h-5 w-5" />
        </Button>
        <Button variant="ghost" size="icon">
          <Bell className="h-5 w-5" />
        </Button>
        <ThemeToggle />
      </div>
    </header>
  );
}
