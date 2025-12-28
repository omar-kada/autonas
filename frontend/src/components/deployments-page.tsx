import { type Deployment } from '@/api/api';
import { getDeploymentsQueryOptions, useIsMobile } from '@/hooks';
import { ArrowLeft } from 'lucide-react';
import { useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useParams } from 'react-router-dom';

import { useDeploymentNavigate } from '@/lib';
import { useInfiniteQuery } from '@tanstack/react-query';
import { DeploymentDetail, DeploymentList, DeploymentToolbar } from './deployment';
import { Button } from './ui/button';
import { Separator } from './ui/separator';

export function DeploymentsPage() {
  const { t } = useTranslation();
  const isMobile = useIsMobile();
  const deploymentNavigate = useDeploymentNavigate();
  var { id } = useParams();
  const { data: deployments } = useInfiniteQuery(getDeploymentsQueryOptions());

  const [selectedItem, setSelectedItem] = useState(id);
  const handleSelect = useCallback(
    (item: Deployment) => {
      deploymentNavigate(item.id);
      setSelectedItem(item.id);
    },
    [deploymentNavigate, setSelectedItem],
  );

  const handleBack = useCallback(() => {
    deploymentNavigate();

    setSelectedItem(undefined);
  }, [deploymentNavigate, setSelectedItem]);
  // Check if data exists and is an object
  if (!deployments || typeof deployments !== 'object' || !deployments.length) {
    return <div>No deployments data available</div>;
  }

  const selectedItemOrDefault = selectedItem ?? deployments[0].id;
  return (
    <>
      {(!selectedItem || !isMobile) && (
        <>
          <DeploymentToolbar />
          <Separator />
        </>
      )}
      <div className="flex flex-1 overflow-hidden">
        {/* Sidebar (hidden on mobile if an item is selected) */}

        <aside
          className={`w-full h-full max-h-full flex flex-col sm:w-75 sm:shrink-0 m-2 pb-4 ${selectedItem ? 'hidden sm:flex' : ''}`}
        >
          <DeploymentList
            selectedDeployment={isMobile ? undefined : selectedItemOrDefault}
            onSelect={handleSelect}
            className="border rounded-lg h-full max-h-full bg-muted/30"
          />
        </aside>

        {/* Main content (hidden on mobile until item is selected) */}
        <main className={`flex-col ${!selectedItem ? 'hidden sm:block sm:flex-1' : 'w-full'}`}>
          {/* Mobile back button (only shows when content is open) */}
          <Button variant="ghost" className="sm:hidden mx-4 my-2" onClick={handleBack}>
            <ArrowLeft className="h-5 w-5" /> {t('BACK')}
          </Button>
          <Separator className="sm:hidden"></Separator>

          <DeploymentDetail id={selectedItemOrDefault}></DeploymentDetail>
        </main>
      </div>
    </>
  );
}
