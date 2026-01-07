import { useRelativeTime } from '@/hooks/use-relative-time';
import { cn } from '@/lib';
import { Tooltip, TooltipContent, TooltipTrigger } from '../ui/tooltip';

export function HumanTime({
  time,
  className,
  defaultValue,
}: {
  time?: Date | string;
  className?: string;
  defaultValue?: string;
}) {
  const relativeTime = useRelativeTime(time);

  if (!time || !relativeTime) {
    return <span className={cn('font-light text-nowrap', className)}>{defaultValue}</span>;
  }

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <span className={cn('font-light text-nowrap', className)}>{relativeTime}</span>
      </TooltipTrigger>
      <TooltipContent side="bottom">{new Date(time).toLocaleString()}</TooltipContent>
    </Tooltip>
  );
}
