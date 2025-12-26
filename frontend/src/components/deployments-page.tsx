import { type Deployment } from '@/api/api';
import { useDeployments, useStats } from '@/hooks';
import { useDeploymentNavigate } from '@/lib';
import { ArrowLeft, History } from 'lucide-react';
import { useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Navigate, useParams } from 'react-router-dom';
import { DeploymentDetail, DeploymentList, DeploymentStatusBadge } from './deployment';
import { Button } from './ui/button';
import { ScrollArea } from './ui/scroll-area';

export function DeploymentsPage() {
  const { t } = useTranslation();
  const deploymentNavigate = useDeploymentNavigate();
  const { deployments, isLoading, error } = useDeployments();
  const { id } = useParams();
  const { data: stats } = useStats(30);

  const [showItem, setShowItem] = useState(false);
  const handleSelect = useCallback((item: Deployment) => {
    deploymentNavigate(item.id);
    setShowItem(true);
  }, []);

  const handleBack = useCallback(() => {
    setShowItem(false);
  }, []);

  if (isLoading) {
    return <div>Loading deployments...</div>;
  }

  if (error || !deployments) {
    return <div>Error fetching deployments: {error?.message}</div>;
  }
  // Check if data exists and is an object
  if (!deployments || typeof deployments !== 'object' || !deployments.data.length) {
    return (
      <>
        <div>No deployments data available</div>;
      </>
    );
  }

  if (id == null) {
    return <Navigate to={deployments.data[0].id}></Navigate>;
  }
  return (
    <div className="flex flex-1 overflow-hidden">
      {/* Sidebar (hidden on mobile if an item is selected) */}

      <aside
        className={`w-full h-full max-h-full flex flex-col sm:w-75 sm:shrink-0 m-2 pb-4 ${showItem ? 'hidden sm:flex' : ''}`}
      >
        {stats && (
          <div className="flex items-center p-2 mb-2 gap-2">
            <span className="text-sm font-light mx-1 flex-1 flex gap-1 items-center">
              <History className="size-4"></History>
              {t('LAST_X_DAYS', { days: 30 })} :
            </span>
            <DeploymentStatusBadge
              status="success"
              label={String(stats.success)}
            ></DeploymentStatusBadge>
            {stats.error ? (
              <DeploymentStatusBadge
                status="error"
                label={String(stats.error)}
              ></DeploymentStatusBadge>
            ) : null}
          </div>
        )}
        <div className="flex-1 h-1">
          <ScrollArea className="p-2 mb-5 border h-full rounded-lg max-h-full bg-muted/30">
            <DeploymentList OnSelect={handleSelect} />
          </ScrollArea>
        </div>
      </aside>

      {/* Main content (hidden on mobile until item is selected) */}
      <main className={`flex-col ${!showItem ? 'hidden sm:block sm:flex-1' : 'w-full'}`}>
        {/* Mobile back button (only shows when content is open) */}
        <Button variant="ghost" className="sm:hidden mx-4 mt-2" onClick={handleBack}>
          <ArrowLeft className="h-5 w-5" /> {t('BACK')}
        </Button>

        <ScrollArea className=" p-6 overflow-auto w-full h-full">
          <DeploymentDetail id={id}></DeploymentDetail>
        </ScrollArea>
      </main>
    </div>
  );
}
