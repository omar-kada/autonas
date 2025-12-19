import type { FileDiff } from '@/api/api';
import { useTheme } from '@/hooks/theme-provider';
import { DiffModeEnum, DiffView } from '@git-diff-view/react';
import { ChevronDown, ChevronUp, FileDiff as DiffIcon } from 'lucide-react';
import { useState } from 'react';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '../ui/collapsible';

export function FileDiffView({ fileDiff, className }: { fileDiff: FileDiff; className?: string }) {
  const { theme } = useTheme();
  const [isOpen, setIsOpen] = useState(false);

  return (
    <Collapsible open={isOpen} onOpenChange={setIsOpen} className={className}>
      <CollapsibleTrigger
        className={`w-full  font-light justify-between items-center flex bg-accent p-2 cursor-pointer ${isOpen ? 'rounded-t-lg' : 'rounded-lg'}`}
      >
        <span className="flex items-center">
          <DiffIcon className="size-4 mr-2"></DiffIcon>
          <span className="font-light">
            {fileDiff.oldFile +
              (fileDiff.oldFile !== fileDiff.newFile ? ` > ${fileDiff.newFile}` : '')}
          </span>
        </span>
        {isOpen ? <ChevronUp /> : <ChevronDown />}
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
