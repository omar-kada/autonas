import { useLogout, useUser } from '@/hooks';
import { useTheme } from '@/hooks/theme-provider';
import { Bell, LogOutIcon, Moon, Settings, User } from 'lucide-react';
import { useCallback, useEffect, useState, type ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import { SettingsSheet } from './settings';
import { Avatar, AvatarFallback } from './ui/avatar';
import { Button } from './ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from './ui/dropdown-menu';
import { Field, FieldLabel } from './ui/field';
import { Switch } from './ui/switch';

export function Topbar({ children }: { children?: ReactNode }) {
  const { t } = useTranslation();
  const { data: user } = useUser();

  const [initial, setInitial] = useState(user?.username ?? '');

  useEffect(() => {
    if (user) {
      setInitial(user.username.charAt(0).toUpperCase());
    }
  }, [user]);

  const { theme, setTheme } = useTheme();

  const toggleTheme = useCallback(() => {
    setTheme(theme === 'dark' ? 'light' : 'dark');
  }, [theme, setTheme]);

  const [openSettings, setOpenSettings] = useState(false);
  const openSettingsSheet = useCallback(() => setOpenSettings(true), [setOpenSettings]);

  const { logout } = useLogout();

  return (
    <header className="h-14 min-h-14 border-b w-full flex items-center justify-between px-4 bg-sidebar sticky top-0 z-50">
      {/* Logo */}
      <div className="text-xl font-semibold mr-5">AirCompose</div>
      <div className="flex-1 w-1 max-w-10">{/*gap*/}</div>
      <div className="flex flex-2 justify-between">{children}</div>

      {/* Settings + Notifications */}
      {user && (
        <div className="flex items-center gap-2">
          <Button variant="ghost" size="icon" className="rounded-full" disabled>
            <Bell className="h-5 w-5" />
          </Button>

          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="rounded-full cursor-pointer">
                <Avatar>
                  <AvatarFallback className="select-none">{initial}</AvatarFallback>
                </Avatar>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuGroup>
                <DropdownMenuLabel>
                  {t('MENU.LOGGED_AS')} : {user?.username ?? ''}
                </DropdownMenuLabel>

                <DropdownMenuItem onClick={toggleTheme}>
                  <Field className="" orientation="horizontal">
                    <Moon></Moon>
                    <FieldLabel className="font-normal pe-2">{t('MENU.DARK_MODE')}</FieldLabel>
                    <Switch checked={theme === 'dark'} onCheckedChange={toggleTheme} />
                  </Field>
                </DropdownMenuItem>

                <DropdownMenuItem disabled>
                  <User />
                  {t('MENU.ACCOUNT')}
                </DropdownMenuItem>
                <DropdownMenuItem onSelect={openSettingsSheet}>
                  <Settings />
                  {t('MENU.SETTINGS')}
                </DropdownMenuItem>
              </DropdownMenuGroup>
              <DropdownMenuSeparator />
              <DropdownMenuItem onSelect={logout}>
                <LogOutIcon />
                {t('ACTION.SIGN_OUT')}
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
          <SettingsSheet open={openSettings} setOpen={setOpenSettings}></SettingsSheet>
        </div>
      )}
    </header>
  );
}
