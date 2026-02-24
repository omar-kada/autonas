import { authAPIRefresh } from '@/api/api';
import { ROUTES } from '@/lib';
import { QueryCache, QueryClient } from '@tanstack/react-query';
import type { AxiosError } from 'axios';
import { debounce } from 'lodash';

const debouncedRefreshToken = debounce(refreshToken, 5000, { leading: true }); // 5 seconds debounce

export const queryClient = new QueryClient({
  queryCache: new QueryCache({
    onError: (error, _) => {
      if ((error as any as AxiosError).status === 401) {
        debouncedRefreshToken();
      }
    },
  }),
  defaultOptions: {
    queries: {
      retry: 0,
      refetchOnWindowFocus: false,
      staleTime: 1000 * 60, // 1 min
      gcTime: 10 * 60 * 1000, // 10 min
    },
  },
});

function refreshToken() {
  authAPIRefresh()
    .then((success) => success && queryClient.refetchQueries())
    .catch((err: AxiosError) => {
      if (err.status === 401) {
        window.location.href = ROUTES.LOGIN;
      }
    });
}
