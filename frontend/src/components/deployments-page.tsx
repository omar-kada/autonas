import { type Deployment } from '@/api/api';
import { getDeploymentsQueryOptions, useIsMobile } from '@/hooks';
import { ArrowLeft, Loader } from 'lucide-react';
import { useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useInView } from 'react-intersection-observer';
import { useParams } from 'react-router-dom';

import { useDeploymentNavigate } from '@/lib';
import { useInfiniteQuery } from '@tanstack/react-query';
import { DeploymentDetail, DeploymentList, DeploymentToolbar } from './deployment';
import { Button } from './ui/button';
import { ScrollArea } from './ui/scroll-area';
import { Separator } from './ui/separator';

export function DeploymentsPage() {
  const { t } = useTranslation();
  const { ref, inView } = useInView();
  const isMobile = useIsMobile();
  const deploymentNavigate = useDeploymentNavigate();
  const {
    data: deployments,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
    error,
  } = useInfiniteQuery(getDeploymentsQueryOptions());
  const { id } = useParams();

  const [showItem, setShowItem] = useState(!!id);
  const handleSelect = useCallback(
    (item: Deployment) => {
      deploymentNavigate(item.id);
      setShowItem(true);
    },
    [deploymentNavigate, setShowItem],
  );

  const handleBack = useCallback(() => {
    deploymentNavigate();

    setShowItem(false);
  }, [deploymentNavigate, setShowItem]);

  useEffect(() => {
    if (inView && hasNextPage && !isFetchingNextPage) {
      fetchNextPage();
    }
  }, [inView, hasNextPage, isFetchingNextPage, fetchNextPage]);

  if (isLoading) {
    return <div>Loading deployments...</div>;
  }

  if (error || !deployments) {
    return <div>Error fetching deployments: {error?.message}</div>;
  }
  // Check if data exists and is an object
  if (!deployments || typeof deployments !== 'object' || !deployments.length) {
    return <div>No deployments data available</div>;
  }

  return (
    <>
      {(!showItem || !isMobile) && (
        <>
          <DeploymentToolbar />
          <Separator />
        </>
      )}
      <div className="flex flex-1 overflow-hidden">
        {/* Sidebar (hidden on mobile if an item is selected) */}

        <aside
          className={`w-full h-full max-h-full flex flex-col sm:w-75 sm:shrink-0 m-2 pb-4 ${showItem ? 'hidden sm:flex' : ''}`}
        >
          <div className="flex-1 h-1">
            <ScrollArea className="mb-5 border h-full rounded-lg max-h-full bg-muted/30">
              <DeploymentList
                deployments={deployments ?? []}
                selectedDeployment={id}
                OnSelect={handleSelect}
              />
              <div ref={ref} className="flex justify-around">
                {(isFetchingNextPage || hasNextPage) && <Loader className="animate-spin my-2" />}
              </div>
            </ScrollArea>
          </div>
        </aside>

        {/* Main content (hidden on mobile until item is selected) */}
        <main className={`flex-col ${!showItem ? 'hidden sm:block sm:flex-1' : 'w-full'}`}>
          {/* Mobile back button (only shows when content is open) */}
          <Button variant="ghost" className="sm:hidden mx-4 my-2" onClick={handleBack}>
            <ArrowLeft className="h-5 w-5" /> {t('BACK')}
          </Button>
          <Separator className="sm:hidden"></Separator>

          <DeploymentDetail id={id ?? deployments[0].id}></DeploymentDetail>
        </main>
      </div>
    </>
  );
}
