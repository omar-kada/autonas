import { Level } from '@/api';

export const LEVELS = [Level.DEBUG, Level.INFO, Level.WARN, Level.ERROR];

export const LEVEL_STYLES: Record<Level, string> = {
  DEBUG:
    'text-zinc-400 border-zinc-300 bg-zinc-50 hover:bg-zinc-100 dark:border-zinc-700 dark:bg-zinc-900/40 dark:hover:bg-zinc-900/70',
  INFO: 'text-sky-600 border-sky-300 bg-sky-50 hover:bg-sky-100 dark:text-sky-400 dark:border-sky-800 dark:bg-sky-950/40 dark:hover:bg-sky-950/70',
  WARN: 'text-amber-600 border-amber-300 bg-amber-50 hover:bg-amber-100 dark:text-amber-400 dark:border-amber-800 dark:bg-amber-950/40 dark:hover:bg-amber-950/70',
  ERROR:
    'text-red-600 border-red-300 bg-red-50 hover:bg-red-100 dark:text-red-400 dark:border-red-800 dark:bg-red-950/40 dark:hover:bg-red-950/70',
};
