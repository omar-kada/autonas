import { humanizeDurationMs } from '@/lib';
import type { TFunction } from 'i18next';
import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

export function humanizeFromNow(
  date: Date | string,
  t: TFunction<'translation', undefined>,
  locale = 'en',
): string | null {
  if (!date) {
    return null;
  }
  const target = new Date(date).getTime();
  const diffMs = target - Date.now();
  if (target < Date.UTC(1970, 1, 1)) {
    return null;
  }
  if (Math.abs(diffMs) < 60_000) {
    if (diffMs < 0) {
      return t('TIME.JUST_NOW');
    } else {
      return t('TIME.IN_FEW_SECONDS');
    }
  }

  return humanizeDurationMs(diffMs, locale);
}

// ---- REACT HOOK -----------------------------------------------

export function useRelativeTime(target: Date | string = '', locale = 'en'): string | null {
  const { t } = useTranslation();
  const [formatted, setFormatted] = useState(() => humanizeFromNow(target, t, locale));
  useEffect(() => {
    const targetMs = new Date(target).getTime();

    function update() {
      setFormatted(humanizeFromNow(target, t, locale));
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
