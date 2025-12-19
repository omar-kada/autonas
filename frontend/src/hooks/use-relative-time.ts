import { humanizeDurationMs } from '@/lib';
import type { TFunction } from 'i18next';
import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

export function humanizeFromNow(
  date: Date | number,
  t: TFunction<'translation', undefined>,
  locale = 'en',
): string {
  const target = typeof date === 'number' ? date : date.getTime();
  const diffMs = target - Date.now();

  if (Math.abs(diffMs) < 60_000) {
    if (diffMs < 0) {
      return t('JUST_NOW');
    } else {
      return t('IN_FEW_SECONDS');
    }
  }

  return humanizeDurationMs(diffMs, locale);
}

// ---- REACT HOOK -----------------------------------------------

export function useRelativeTime(target: Date | number, locale = 'en'): string {
  const { t } = useTranslation();
  const [formatted, setFormatted] = useState(() => humanizeFromNow(target, t, locale));
  useEffect(() => {
    const targetMs = typeof target === 'number' ? target : target.getTime();

    function update() {
      setFormatted(humanizeFromNow(targetMs, t, locale));
    }

    update(); // run immediately

    const diff = Math.abs(Date.now() - targetMs);

    // Smart interval selection
    let interval = 60_000; // 1 minute
    if (diff > 3_600_000) interval = 3_600_000; // 1 hour

    const id = setInterval(update, interval);
    return () => clearInterval(id);
  }, [target, locale]);

  return formatted;
}
