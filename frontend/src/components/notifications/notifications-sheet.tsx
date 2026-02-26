import { EventType, type Event } from '@/api/api';
import { Button } from '@/components/ui/button';
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet';
import { getNotificationsQueryOptions, useRelativeTime } from '@/hooks';
import { useInfiniteQuery } from '@tanstack/react-query';
import { useCallback, useEffect, useState, type ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import { useInView } from 'react-intersection-observer';
import { useNavigate } from 'react-router-dom';
import { ScrollArea } from '../ui/scroll-area';
import { Spinner } from '../ui/spinner';
import { ErrorAlert } from '../view';
import { NotificationBadge } from './notification-badge';

export function NotificationSheet({ children }: { children: ReactNode }) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const {
    data: notifications,
    isPending,
    error,
    isFetchingNextPage,
    hasNextPage,
    fetchNextPage,
  } = useInfiniteQuery(getNotificationsQueryOptions());
  const { ref, inView } = useInView();
  useEffect(() => {
    if (inView && hasNextPage && !isFetchingNextPage && notifications?.length) {
      fetchNextPage();
    }
  }, [inView, hasNextPage, isFetchingNextPage, fetchNextPage, notifications?.length]);

  const [open, setOpen] = useState(false);
  const handleNotificationClick = useCallback(
    (event: Event) => {
      setOpen(false);
      navigate(event.type);
    },
    [setOpen, navigate],
  );

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>{children}</SheetTrigger>
      <SheetContent className="w-full md:w-none flex flex-col h-full">
        <SheetHeader>
          <SheetTitle>{t('NOTIFICATIONS.NOTIFICATIONS')}</SheetTitle>
        </SheetHeader>
        <ScrollArea className="h-1 flex-1 gap-2">
          {isPending && <Spinner />}
          {error && (
            <ErrorAlert title={t('ALERT.LOAD_NOTIFICATIONS_ERROR')} details={error.message} />
          )}
          {notifications && (
            <div className="grid auto-rows-min px-4 mb-10">
              {notifications.map((notification) => (
                <Notification
                  key={notification.ID}
                  notification={notification}
                  onClick={handleNotificationClick}
                />
              ))}
            </div>
          )}
          <div ref={ref} className="flex justify-around w-full min-h-1">
            {(isFetchingNextPage || (hasNextPage && notifications?.length)) && <Spinner />}
          </div>
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

function Notification({
  notification,
  onClick,
}: {
  notification: Event;
  onClick: (notif: Event) => void;
}) {
  const { t } = useTranslation();
  const relativeTime = useRelativeTime(notification.time);
  return (
    <div
      className="p-3 cursor-pointer hover:bg-accent border-b"
      onClick={() => onClick(notification)}
    >
      <div className="flex items-center gap-2">
        <NotificationBadge type={notification.type} />
        <span className="self-center">
          {t(`NOTIFICATIONS.TITLE.${notification.type}`, { ...notification })}
        </span>
      </div>
      <div className="text-sm text-muted-foreground mt-1">
        {getNotificaitonDetails(notification)}
      </div>
      <div className="text-xs text-muted-foreground mt-2">{relativeTime}</div>
    </div>
  );
}

function getNotificaitonDetails(notif: Event): string {
  switch (notif.type) {
    case EventType.ERROR:
    case EventType.MISC:
      if (notif.objectId && notif.objectName) {
        return `#${notif.objectId} \"${notif.objectName}\" - ${notif.msg}`;
      }
      return notif.msg;
    default:
      return '';
  }
}
