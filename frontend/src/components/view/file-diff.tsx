import type { FileDiff } from '@/api/api';
import { useTheme } from '@/hooks/theme-provider';
import { DiffModeEnum, DiffView } from '@git-diff-view/react';
import { ChevronsUpDown } from 'lucide-react';
import { useState } from 'react';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '../ui/collapsible';

export function FileDiffView({ fileDiff, className }: { fileDiff: FileDiff; className?: string }) {
  const { theme } = useTheme();
  const [isOpen, setIsOpen] = useState(false);

  return (
    <Collapsible open={isOpen} onOpenChange={setIsOpen} className={className}>
      <CollapsibleTrigger
        className={`w-full justify-between flex bg-accent p-2 cursor-pointer ${isOpen ? 'rounded-t-lg' : 'rounded-lg'}`}
      >
        {fileDiff.oldFile + (fileDiff.oldFile !== fileDiff.newFile ? ` > ${fileDiff.newFile}` : '')}
        <ChevronsUpDown />
      </CollapsibleTrigger>
      <CollapsibleContent>
        <DiffView<string>
          className="border-color overflow-hidden border rounded-b-lg"
          data={{
            oldFile: { fileName: fileDiff.oldFile },
            newFile: { fileName: fileDiff.newFile },
            hunks: [fileDiff.diff],
          }}
          diffViewTheme={theme}
          diffViewHighlight
          diffViewMode={DiffModeEnum.Split}
          diffViewWrap
        />
      </CollapsibleContent>
    </Collapsible>
  );
}
