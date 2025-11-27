import { Button } from '@/components/ui/button'; // shadcn
import { useTheme } from '@/hooks/theme-provider';
import { Moon, Sun } from 'lucide-react';
import { useTranslation } from 'react-i18next';

export function ThemeToggle() {
  const { t } = useTranslation();
  const { theme, setTheme } = useTheme();

  return (
    <Button
      variant="ghost"
      size="icon"
      onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
    >
      <Sun className="h-5 w-5 transform rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
      <Moon className="absolute h-5 w-5 transform rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />

      <span className="sr-only">{t('TOGGLE_THEME')}</span>
    </Button>
  );
}
