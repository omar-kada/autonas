import type { FileDiff } from '@/api/api';
import { GitCompare } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { FileDiffView } from '.';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../ui/card';

export function DeploymentDiff({ fileDiffs }: { fileDiffs: FileDiff[] }) {
  const { t } = useTranslation();
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex">
          <GitCompare className="size-5 mx-1" />
          {t('UPDATED_FILES')}
        </CardTitle>
        <CardDescription>{fileDiffs.length} updated files</CardDescription>
      </CardHeader>
      <CardContent>
        {fileDiffs.map((fileDiff) => (
          <FileDiffView fileDiff={fileDiff} key={fileDiff.oldFile} className="mb-2" />
        ))}
      </CardContent>
    </Card>
  );
}
