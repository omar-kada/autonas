import { Bell, Settings } from 'lucide-react';
import { Navbar } from './navbar';
import { Button } from './ui/button';
import { ThemeToggle } from './view/theme-toggle';

export function Topbar() {
  return (
    <header className="h-14 min-h-14 border-b w-full flex items-center justify-between px-4 bg-sidebar sticky top-0 z-50">
      {/* Logo */}
      <div className="text-xl font-semibold">MyLogo</div>

      <Navbar />
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
