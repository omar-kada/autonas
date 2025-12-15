import { type Deployment } from '@/api/api';
import { useDeployments } from '@/hooks';
import { useDeploymentNavigate } from '@/lib';
import { ArrowLeft } from 'lucide-react';
import { useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useParams } from 'react-router-dom';
import { Button } from './ui/button';
import { ScrollArea } from './ui/scroll-area';
import { DeploymentDetail, DeploymentList } from './view';

export function DeploymentsPage() {
  const { t } = useTranslation();
  const deploymentNavigate = useDeploymentNavigate();
  const { deployments, isLoading, error } = useDeployments();
  const { id } = useParams();

  if (id == null && deployments != null) {
    deploymentNavigate(deployments[0].id);
  }

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

  if (error) {
    return <div>Error fetching deployments: {error?.message}</div>;
  }
  // Check if data exists and is an object
  if (!deployments || typeof deployments !== 'object') {
    return <div>No deployments data available</div>;
  }

  return (
    <div className="flex flex-1 overflow-hidden">
      {/* Sidebar (hidden on mobile if an item is selected) */}

      <aside
        className={`w-full sm:w-75 sm:shrink-0 border-r bg-muted/30 ${showItem ? 'hidden sm:block' : ''}`}
      >
        <ScrollArea className="h-full m-2">
          <DeploymentList
            deployments={deployments}
            selectedDeployment={id}
            OnSelect={handleSelect}
          />
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
