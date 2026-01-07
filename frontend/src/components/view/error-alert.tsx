import { cn } from '@/lib';
import { AlertCircleIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { Alert, AlertDescription, AlertTitle } from '../ui/alert';

export function ErrorAlert({
  title,
  details,
  className,
}: {
  title: string | null;
  details?: string | null;
  className?: string;
}) {
  const { t } = useTranslation();
  return (
    title && (
      <Alert variant="destructive" className={cn('w-auto', className)}>
        <AlertCircleIcon />
        <AlertTitle>{t(title)}</AlertTitle>
        {details && <AlertDescription>{t(details)}</AlertDescription>}
      </Alert>
    )
  );
}
