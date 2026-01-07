import { getStatsQueryOptions } from '@/hooks';
import { cn, ROUTES } from '@/lib';
import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import { ContainerStatusBadge } from '.';
import { Skeleton } from '../ui/skeleton';

export function EnvironementHealth({ className }: { className?: string }) {
  const { data: stats, isPending, error } = useQuery(getStatsQueryOptions());

  if (isPending) {
    return <Skeleton className="h-4 w-20 my-auto" />;
  }

  if (error) {
    return <div>Error fetching stats: {error?.message}</div>;
  }

  return (
    <Link
      to={ROUTES.STATUS}
      className={cn('flex flex-wrap items-center align-bottom gap-4 m-4 ', className)}
    >
      <ContainerStatusBadge status={stats.health}></ContainerStatusBadge>
    </Link>
  );
}
