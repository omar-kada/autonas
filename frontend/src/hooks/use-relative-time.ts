import type { TFunction } from 'i18next';
import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

// ---- TYPES ----------------------------------------------------

type RelativeTimeUnit = 'second' | 'minute' | 'hour' | 'day' | 'week' | 'month' | 'year';

interface Division {
  readonly amount: number;
  readonly unit: RelativeTimeUnit;
}

// ---- HUMANIZER ------------------------------------------------

const divisions: Array<Division> = [
  { amount: 1, unit: 'second' },
  { amount: 60, unit: 'minute' },
  { amount: 60, unit: 'hour' },
  { amount: 24, unit: 'day' },
  { amount: 7, unit: 'week' },
  { amount: 30, unit: 'month' }, // average month
  { amount: 12, unit: 'year' }, // average year
] as const;

export function humanize(
  date: Date | number,
  t: TFunction<'translation', undefined>,
  locale = 'en',
): string {
  const target = typeof date === 'number' ? date : date.getTime();
  const diffMs = target - Date.now();

  const rtf = new Intl.RelativeTimeFormat(locale, { numeric: 'auto' });

  let duration = diffMs / 1000; // seconds
  let unit: RelativeTimeUnit = 'second';

  for (const division of divisions) {
    if (Math.abs(duration) < division.amount) break;
    duration /= division.amount;
    unit = division.unit;
  }
  if (unit === 'second') {
    return t('JUST_NOW');
  }

  return rtf.format(Math.round(duration), unit);
}

// ---- REACT HOOK -----------------------------------------------

export function useRelativeTime(target: Date | number, locale = 'en'): string {
  const { t } = useTranslation();
  const [formatted, setFormatted] = useState(() => humanize(target, t, locale));
  useEffect(() => {
    const targetMs = typeof target === 'number' ? target : target.getTime();

    function update() {
      setFormatted(humanize(targetMs, t, locale));
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
