import { EventType, type Event } from '@/api/api';
import { Button } from '@/components/ui/button';
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet';
import { useDeploymentNavigate } from '@/lib';
import { useCallback, useState, type ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import { ScrollArea } from '../ui/scroll-area';
import { Separator } from '../ui/separator';
import { NotificationFilter } from './notifications-filters';
import { NotificationList } from './notifications-list';

export function NotificationSheet({ children }: { children: ReactNode }) {
  const { t } = useTranslation();
  const depNavigate = useDeploymentNavigate();
  const [open, setOpen] = useState(false);
  const handleNotificationClick = useCallback(
    (event: Event) => {
      setOpen(false);
      depNavigate(event.objectId + '');
    },
    [setOpen, depNavigate],
  );
  const [selectedTypes, setSelectedTypes] = useState<Array<EventType>>(Object.values(EventType));

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>{children}</SheetTrigger>
      <SheetContent
        className="w-full md:w-none flex flex-col h-full"
        aria-describedby={t('NOTIFICATIONS.DESCRIPTION')}
      >
        <SheetHeader>
          <SheetTitle>{t('NOTIFICATIONS.NOTIFICATIONS')}</SheetTitle>
          <div className="flex flex-nowrap justify-between items-center-safe">
            <SheetDescription>{t('NOTIFICATIONS.DESCRIPTION')}</SheetDescription>

            <NotificationFilter onFilterChanged={setSelectedTypes} />
          </div>
        </SheetHeader>
        <Separator></Separator>
        <ScrollArea className="h-1 flex-1 gap-2">
          <NotificationList
            onNotificationClick={handleNotificationClick}
            selectedTypes={selectedTypes}
          />
        </ScrollArea>
        <SheetFooter>
          <SheetClose asChild>
            <Button variant="outline">{t('ACTION.CLOSE')}</Button>
          </SheetClose>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
