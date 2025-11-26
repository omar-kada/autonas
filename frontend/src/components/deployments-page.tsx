import { useDeployments } from '@/hooks';
import type { Deployment } from '@/models/deployment';
import { ArrowLeft } from 'lucide-react';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Button } from './ui/button';
import { ScrollArea } from './ui/scroll-area';
import { DeploymentList } from './view';

export function DeploymentsPage() {
  const { t } = useTranslation();
  const { data, isLoading, error } = useDeployments();
  const [selectedItem, setSelectedItem] = useState<Deployment>(null as any as Deployment);
  const [showItem, setShowItem] = useState(false);
  const handleSelect = (item: Deployment) => {
    setSelectedItem(item);
    setShowItem(true);
  };

  const handleBack = () => {
    setShowItem(false);
  };

  if (isLoading) {
    return <div>Loading deployments...</div>;
  }

  if (error) {
    return <div>Error fetching deployments: {error.message}</div>;
  }

  // Check if data exists and is an object
  if (!data || typeof data !== 'object') {
    return <div>No deployments data available</div>;
  }

  return (
    <>
      <div className="flex flex-1 overflow-hidden">
        {/* Sidebar (hidden on mobile if an item is selected) */}

        <aside
          className={`w-full sm:w-75 sm:shrink-0 border-r bg-muted/30 ${showItem ? 'hidden sm:block' : ''}`}
        >
          <ScrollArea className="h-full m-2">
            <DeploymentList deployments={data} OnSelect={handleSelect} />
          </ScrollArea>
        </aside>

        {/* Main content (hidden on mobile until item is selected) */}
        <main className={`flex-col ${!showItem ? 'hidden sm:block sm:flex-1' : 'w-full'}`}>
          {/* Mobile back button (only shows when content is open) */}
          <Button variant="ghost" className="sm:hidden mx-4 mt-2" onClick={handleBack}>
            <ArrowLeft className="h-5 w-5" /> {t('BACK')}
          </Button>

          <ScrollArea className=" p-6 overflow-auto w-full h-full">
            {selectedItem != null ? (
              <>
                <h3 className="text-2xl font-semibold mb-4">{selectedItem.name}</h3>
                <pre>{selectedItem.diff}</pre>
              </>
            ) : (
              <div>{t('SELECT_DEPLOYMENT_FOR_DETAILS')}</div>
            )}
          </ScrollArea>
        </main>
      </div>
    </>
  );
}
