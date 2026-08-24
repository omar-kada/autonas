import {
  ServerMessageLogKind,
  ServerMessagePreviousLogsKind,
  type Error,
  type LogMessages,
  type ServerMessageLog,
  type ServerMessagePreviousLogs,
} from '@/api';
import { QueryClient, useQuery, type QueryKey, type UseQueryOptions } from '@tanstack/react-query';
import { useEffect } from 'react';
import { useWs, useWsStatusQuery } from '..';

export function getLogsQueryKey(): QueryKey {
  return ['logs'];
}

export function useLogs(
  previousLines = 0,
  queryOptions?: Partial<UseQueryOptions<LogMessages, Error>>,
) {
  const { startLogs, endLogs } = useWs();
  const { data: status } = useWsStatusQuery();
  /*
  useEffect(() => {
    startLogs(previousLines);
    return () => endLogs();
  }, [startLogs, endLogs]);*/

  useEffect(() => {
    if (status === 'connected') {
      startLogs(previousLines);
      return () => endLogs();
    }
  }, [status]);

  return useQuery<LogMessages, Error>({
    queryKey: getLogsQueryKey(),
    queryFn: async (): Promise<LogMessages> => [],
    enabled: true,
    refetchOnMount: false,
    refetchOnReconnect: false,
    ...queryOptions,
  });
}

export function onLogEvent(
  serverEvent: ServerMessageLog | ServerMessagePreviousLogs,
  queryClient: QueryClient,
) {
  switch (serverEvent.kind) {
    case ServerMessageLogKind.log:
      queryClient.setQueryData(getLogsQueryKey(), (prev: LogMessages = []) => [
        ...prev,
        serverEvent.value,
      ]);
      break;

    case ServerMessagePreviousLogsKind.previousLogs:
      queryClient.setQueryData(getLogsQueryKey(), serverEvent.value);
      break;
  }
}
