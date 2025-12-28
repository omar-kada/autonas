type RelativeTimeUnit = 'second' | 'minute' | 'hour' | 'day' | 'week' | 'month' | 'year';

interface Division {
  readonly amount: number;
  readonly unit: RelativeTimeUnit;
  readonly short?: string;
}

const divisions: Array<Division> = [
  { amount: 1, unit: 'second' },
  { amount: 60, unit: 'minute' },
  { amount: 60, unit: 'hour' },
  { amount: 24, unit: 'day' },
  { amount: 7, unit: 'week' },
  { amount: 30, unit: 'month', short: 'M' }, // average month
  { amount: 12, unit: 'year' }, // average year
] as const;

export function humanizeDurationMs(diffMs: number, locale = 'en'): string {
  const rtf = new Intl.RelativeTimeFormat(locale, { numeric: 'auto' });
  let duration = diffMs / 1000; // seconds
  let unit: RelativeTimeUnit = 'second';

  for (const division of divisions) {
    if (Math.abs(duration) < division.amount) break;
    duration /= division.amount;
    unit = division.unit;
  }

  return rtf.format(Math.round(duration), unit);
}

export function humanizeDuration(startDate: string | Date, endDate: string | Date) {
  const start = new Date(startDate);
  const end = new Date(endDate);
  const diffMs = end.getTime() - start.getTime();
  return humanizeDurationMs(diffMs);
}

export function formatElapsedMs(diffMs: number): string {
  let duration = diffMs / 1000; // seconds
  let unit: RelativeTimeUnit = 'second';

  for (const division of divisions) {
    if (Math.abs(duration) < division.amount) break;
    duration /= division.amount;
    unit = division.unit;
  }

  const shortUnit = divisions.find((d) => d.unit === unit)?.short || unit.charAt(0);
  return `${Math.round(duration)}${shortUnit}`;
}

export function formatElapsed(startDate?: string | Date, endDate?: string | Date): string {
  if (!isDateValid(startDate) || !isDateValid(endDate)) {
    return '-';
  }
  const start = new Date(startDate);
  const end = new Date(endDate);
  const diffMs = end.getTime() - start.getTime();
  return formatElapsedMs(diffMs);
}

export function isDateValid(date?: string | Date): date is string | Date {
  return date != null && new Date(date).getFullYear() > 1970;
}
