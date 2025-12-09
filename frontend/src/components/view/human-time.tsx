import { useRelativeTime } from '@/hooks/use-relative-time';
import { Tooltip, TooltipContent, TooltipTrigger } from '../ui/tooltip';

export function HumanTime({ time }: { time: Date | string }) {
  const relativeTime = useRelativeTime(new Date(time));

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <span className="text-sm font-light">{time ? relativeTime : ''}</span>
      </TooltipTrigger>
      <TooltipContent side="bottom">{new Date(time).toLocaleString()}</TooltipContent>
    </Tooltip>
  );
}
