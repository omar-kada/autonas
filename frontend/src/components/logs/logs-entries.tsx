import type { LogLine, LogMessages } from '@/api';
import { Badge } from '@/components/ui/badge';
import { useTranslation } from 'react-i18next';
import { AutoScrollArea } from './auto-scroll-area';
import { LEVEL_STYLES } from './log-levels';

export function LogEntries({ logs }: { logs: LogMessages }) {
  const { t } = useTranslation();
  return (
    <AutoScrollArea className="min-h-0 flex-1">
      <ul className="font-mono text-xs">
        {logs.length === 0 ? (
          <li className="py-4 text-center text-muted-foreground text-sm">
            {t('LOGS.NO_LOGS_FILTERED')}
          </li>
        ) : (
          logs.map((line: LogLine) => (
            <li
              key={line.time + line.msg}
              className="flex items-start gap-2 py-1 border-b border-border/50"
            >
              <span className="text-muted-foreground shrink-0 tabular-nums">
                {new Date(line.time).toLocaleTimeString()}
              </span>
              <Badge
                variant="outline"
                className={`shrink-0 min-w-10 h-4 justify-center rounded px-1.5 text-[0.6rem] font-semibold uppercase ${
                  LEVEL_STYLES[line.level] ?? LEVEL_STYLES.DEBUG
                }`}
              >
                {line.level}
              </Badge>
              <span className="break-all">{line.msg}</span>
            </li>
          ))
        )}
      </ul>

      <LiveIndicator />
    </AutoScrollArea>
  );
}
function LiveIndicator() {
  const { t } = useTranslation();
  return (
    <div className="flex items-center gap-1.5 py-2 text-[10px] text-muted-foreground">
      <span className="relative flex h-1.5 w-1.5">
        <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-emerald-500 opacity-75" />
        <span className="relative inline-flex h-1.5 w-1.5 rounded-full bg-emerald-500" />
      </span>
      {t('LOGS.STREAMING')}
    </div>
  );
}
