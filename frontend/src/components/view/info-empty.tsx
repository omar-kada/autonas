import { Info } from 'lucide-react';
import type { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import {
  Empty,
  EmptyContent,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from '../ui/empty';

export function InfoEmpty({
  title,
  details,
  className,
  children,
}: {
  title: string | null;
  details?: string | null;
  className?: string;
  children?: ReactNode;
}) {
  const { t } = useTranslation();
  return (
    title && (
      <Empty className={className}>
        <EmptyHeader>
          <EmptyMedia variant="icon">
            <Info />
          </EmptyMedia>
          <EmptyTitle>{t(title)}</EmptyTitle>
          {details && <EmptyDescription>{t(details)}</EmptyDescription>}
        </EmptyHeader>
        <EmptyContent>{children}</EmptyContent>
      </Empty>
    )
  );
}
