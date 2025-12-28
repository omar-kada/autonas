import { useStats } from '@/hooks';
import { cn, ROUTES } from '@/lib';
import { Link } from 'react-router-dom';
import { ContainerStatusBadge } from '.';

export function EnvironementHealth({ className }: { className?: string }) {
  const { data: stats, isLoading, error } = useStats(30);

  if (isLoading) {
    return <div>Loading stats...</div>;
  }

  if (error || stats == null) {
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
