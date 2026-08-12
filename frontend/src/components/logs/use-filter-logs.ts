import type { LogMessages } from '@/api';
import { useMemo, useState } from 'react';
import { LEVELS } from './log-levels';

export function useLogFilter(logs: LogMessages) {
  const [text, setText] = useState('');
  const [activeLevels, setActiveLevels] = useState(new Set(LEVELS));

  const filtered = useMemo(
    () =>
      logs.filter(
        (line) =>
          activeLevels.has(line.level) &&
          (!text || line.msg.toLowerCase().includes(text.toLowerCase())),
      ),
    [logs, activeLevels, text],
  );

  return { text, setText, activeLevels, setActiveLevels, filtered };
}
