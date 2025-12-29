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
import { AsideLayout } from './view/aside-layout';

export function DeploymentsPage() {
  const { t } = useTranslation();
  const isMobile = useIsMobile();
  const deploymentNavigate = useDeploymentNavigate();
  const { id } = useParams();
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
    <AsideLayout
      focusMain={selectedItem != null}
      header={(!selectedItem || !isMobile) && <DeploymentToolbar />}
      aside={
        <DeploymentList
          selectedDeployment={isMobile ? undefined : selectedItemOrDefault}
          onSelect={handleSelect}
          className="border rounded-lg h-full max-h-full bg-muted/30"
        />
      }
    >
      {/* Mobile back button (only shows when content is open) */}
      <Button variant="ghost" className="sm:hidden mx-4 my-2" onClick={handleBack}>
        <ArrowLeft className="h-5 w-5" /> {t('BACK')}
      </Button>
      <Separator className="sm:hidden"></Separator>

      <DeploymentDetail id={selectedItemOrDefault} />
    </AsideLayout>
  );
}
