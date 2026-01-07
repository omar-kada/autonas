import type { Deployment } from '@/api/api';
import { getDeploymentsQueryOptions } from '@/hooks';
import { cn } from '@/lib';
import { useInfiniteQuery } from '@tanstack/react-query';
import { Loader } from 'lucide-react';
import { useEffect } from 'react';
import { useInView } from 'react-intersection-observer';
import { DeploymentItemSkeleton, DeploymentListItem } from '.';
import { ScrollArea } from '../ui/scroll-area';

export function DeploymentList({
  selectedDeployment,
  onSelect,
  className,
}: {
  selectedDeployment?: string;
  onSelect: (item: Deployment) => void;
  className?: string;
}) {
  const { ref, inView } = useInView();
  const {
    data: deployments,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isPending,
    error,
  } = useInfiniteQuery(getDeploymentsQueryOptions());

  useEffect(() => {
    if (inView && hasNextPage && !isFetchingNextPage) {
      fetchNextPage();
    }
  }, [inView, hasNextPage, isFetchingNextPage, fetchNextPage]);

  if (error) {
    return <div>Error fetching deployments: {error?.message}</div>;
  }

  if (isPending) {
    return DeploymentListSkeleton(className);
  }

  // Check if data exists and is an object
  if (!deployments || typeof deployments !== 'object' || !deployments.length) {
    return <div>No deployments data available</div>;
  }

  return (
    <ScrollArea className={cn('p-3', className)}>
      <div className="flex flex-col gap-2">
        {deployments.map((deployment) => (
          <DeploymentListItem
            key={deployment.id}
            deployment={deployment}
            isSelected={deployment.id === selectedDeployment}
            onSelect={onSelect}
          ></DeploymentListItem>
        ))}
        <div ref={ref} className="flex justify-around">
          {(isFetchingNextPage || hasNextPage) && <Loader className="animate-spin my-2" />}
        </div>
      </div>
    </ScrollArea>
  );
}
function DeploymentListSkeleton(className?: string) {
  return (
    <ScrollArea className={cn('p-3', className)}>
      <div className="flex flex-col gap-2">
        {Array.from({ length: 5 }, (_, index) => (
          <DeploymentItemSkeleton key={index} />
        ))}
      </div>
    </ScrollArea>
  );
}
