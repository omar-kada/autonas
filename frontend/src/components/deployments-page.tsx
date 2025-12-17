import { type Deployment } from '@/api/api';
import { useDeployments } from '@/hooks';
import { useDeploymentNavigate } from '@/lib';
import { ArrowLeft } from 'lucide-react';
import { useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Navigate, useParams } from 'react-router-dom';
import { DeploymentDetail, DeploymentList } from './deployment';
import { Button } from './ui/button';
import { ScrollArea } from './ui/scroll-area';

export function DeploymentsPage() {
  const { t } = useTranslation();
  const deploymentNavigate = useDeploymentNavigate();
  const { deployments, isLoading, error } = useDeployments();
  const { id } = useParams();

  const [showItem, setShowItem] = useState(false);
  const handleSelect = useCallback((item: Deployment) => {
    if (item.id !== id) {
      deploymentNavigate(item.id);
      setShowItem(true);
    }
  }, []);

  const handleBack = useCallback(() => {
    setShowItem(false);
  }, []);

  if (isLoading) {
    return <div>Loading deployments...</div>;
  }

  if (error) {
    return <div>Error fetching deployments: {error?.message}</div>;
  }
  // Check if data exists and is an object
  if (!deployments || typeof deployments !== 'object' || !deployments.length) {
    return <div>No deployments data available</div>;
  }

  if (id == null) {
    return <Navigate to={deployments[0].id}></Navigate>;
  }
  return (
    <div className="flex flex-1 overflow-hidden">
      {/* Sidebar (hidden on mobile if an item is selected) */}

      <aside
        className={`w-full sm:w-75 sm:shrink-0 border-r bg-muted/30 ${showItem ? 'hidden sm:block' : ''}`}
      >
        <ScrollArea className="h-full m-2">
          <DeploymentList deployments={deployments} OnSelect={handleSelect} />
        </ScrollArea>
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
