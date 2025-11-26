import type { Deployment } from '@/models/deployment';
import { useTranslation } from 'react-i18next';
import { Card, CardContent, CardHeader, CardTitle } from '../ui/card';

export function DeploymentList(props: {
  deployments: Deployment[];
  OnSelect: (item: Deployment) => void;
}) {
  const { t } = useTranslation();

  return (
    <div className="space-y-2">
      {props.deployments.map((deployment) => (
        <Card
          key={deployment.id}
          className="cursor-pointer"
          onClick={() => props.OnSelect(deployment)}
        >
          <CardHeader>
            <CardTitle className="text-sm">{deployment.name}</CardTitle>
          </CardHeader>
          <CardContent className="text-sm text-muted-foreground">
            <p>
              {t('Time')}: {new Date(deployment.time).toLocaleString()}
            </p>
            <pre>{deployment.diff}</pre>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
