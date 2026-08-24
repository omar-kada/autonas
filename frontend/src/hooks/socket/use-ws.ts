import { useQuery, useQueryClient, type QueryKey } from '@tanstack/react-query';
import { useCallback, useContext, useEffect, useRef, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { type SocketEmitterActions } from './socket-emitter';
import { BASE_DELAY, MAX_DELAY, MAX_RETRIES, WsContext, type WsStatus } from './socket-provider';
import { wsToast, wsToastOnStatus } from './socket-toast';

export function useWs(): SocketEmitterActions {
  const ctx = useContext(WsContext);
  if (!ctx) throw new Error('useWs must be used inside <WebSocketProvider>');
  return ctx;
}

export function getWsStatusKey(): QueryKey {
  return ['ws-status'];
}

export function useWsStatusQuery() {
  return useQuery<WsStatus, Error>({
    queryKey: getWsStatusKey(),
    queryFn: async (): Promise<WsStatus> => {
      return 'off';
    },
    enabled: true,
    refetchOnMount: false,
    refetchOnReconnect: false,
    initialData: 'off',
  });
}

export function useWsStatus(enabled: boolean) {
  const { t } = useTranslation();
  const queryClient = useQueryClient();

  const { data: status } = useWsStatusQuery();
  const setStatus = useCallback(
    (status: WsStatus) => {
      queryClient.setQueryData(getWsStatusKey(), status);
    },
    [queryClient],
  );
  const enabledRef = useRef(enabled);

  useEffect(() => {
    enabledRef.current = enabled;
  }, [enabled]);

  const updateStatus = useCallback(
    (newStatus: WsStatus, attempt = 0) => {
      if (enabledRef.current) {
        wsToastOnStatus(t, status, newStatus, attempt, MAX_RETRIES);
      } else {
        wsToast.dismiss();
      }
      setStatus(newStatus);
    },
    [t, status, setStatus, enabledRef],
  );

  return { status, updateStatus };
}

export function useWsRetry() {
  const retriesRef = useRef(0);
  const retryTimeoutRef = useRef<ReturnType<typeof setTimeout>>(0);
  const [attempt, setAttempt] = useState(0);

  const scheduleRetry = useCallback(() => {
    if (retriesRef.current >= MAX_RETRIES) return false;
    const delay = Math.min(BASE_DELAY * 2 ** retriesRef.current, MAX_DELAY);
    const jitter = Math.random() * 0.3 * delay;
    retriesRef.current++;
    retryTimeoutRef.current = setTimeout(() => setAttempt((n) => n + 1), delay + jitter);
    return true;
  }, [retriesRef, retryTimeoutRef]);

  const reset = useCallback(() => {
    retriesRef.current = 0;
  }, [retriesRef]);

  const cancel = useCallback(() => {
    clearTimeout(retryTimeoutRef.current);
  }, [retryTimeoutRef]);

  return { attempt, scheduleRetry, reset, cancel, retriesRef };
}
