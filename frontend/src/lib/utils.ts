import { ErrorCode } from '@/api/api';
import { AxiosError } from 'axios';
import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatTime(time: string, locale: string): string {
  const date = new Date(time);
  return date.toLocaleTimeString(locale, {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
  });
}

export function isInvalidToken(error: unknown): boolean {
  if (error instanceof AxiosError && error.status) {
    return error.status === 401 && error.response?.data?.code === ErrorCode.INVALID_TOKEN;
  }
  return false;
}
