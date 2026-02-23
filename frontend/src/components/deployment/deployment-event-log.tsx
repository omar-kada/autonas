import type { Event } from '@/api/api';
import { formatTime, logColor } from '@/lib';
import { useTranslation } from 'react-i18next';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../ui/card';

export function DeploymentEventLog({ events }: { events: Event[] }) {
  const { t, i18n } = useTranslation();
  return (
    <Card>
      <CardHeader>
        <CardTitle>{t('DEPLOYMENTS.EVENTS_LOG')}</CardTitle>
        <CardDescription></CardDescription>
      </CardHeader>
      <CardContent>
        {events.map((event) => (
          <pre key={event.ID} className={`whitespace-pre-wrap ${logColor(event.type)}`}>
            {formatTime(event.time, i18n.language)} : [{t('EVENT_TYPE.' + event.type)}] {event.msg}
          </pre>
        ))}
      </CardContent>
    </Card>
  );
}
