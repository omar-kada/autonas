import type { FileDiff } from '@/api/api';
import { GitCompare } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { FileDiffView } from '.';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../ui/card';
import { Skeleton } from '../ui/skeleton';

export function DeploymentDiff({
  fileDiffs,
  autoOpen,
}: {
  fileDiffs: FileDiff[];
  autoOpen?: boolean;
}) {
  const { t } = useTranslation();
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex">
          <GitCompare className="size-5 mx-1" />
          {t('UPDATED_FILES')}
        </CardTitle>
        <CardDescription>{t('X_UPDATED_FILES', { count: fileDiffs.length })}</CardDescription>
      </CardHeader>
      <CardContent>
        {fileDiffs.map((fileDiff) => (
          <FileDiffView
            fileDiff={fileDiff}
            key={fileDiff.oldFile}
            autoOpen={autoOpen}
            className="mb-2"
          />
        ))}
      </CardContent>
    </Card>
  );
}
export function DeploymentDiffSkeleton() {
  return (
    <div className="flex flex-col gap-4">
      <Skeleton className="h-6 w-2/3 " />
      <Skeleton className="h-4 w-50" />
      <Skeleton className="h-30 w-full " />
    </div>
  );
}
