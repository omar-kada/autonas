import {
  getNotificationsAPIListQueryKey,
  notificationsAPIList,
  type Event,
  type NotificationsAPIList200,
  type NotificationsAPIListParams,
} from '@/api/api';
import { type InfiniteData, type QueryClient } from '@tanstack/react-query';
import type { AxiosResponse } from 'axios';

const initialParams = { limit: 10, offset: '' } as NotificationsAPIListParams;

export function getNotificationsQueryOptions() {
  return {
    queryKey: getNotificationsAPIListQueryKey(initialParams),
    queryFn: ({ pageParam = initialParams }: { pageParam: NotificationsAPIListParams }) =>
      notificationsAPIList(pageParam),
    initialPageParam: initialParams,
    select: (
      data: InfiniteData<AxiosResponse<NotificationsAPIList200>, NotificationsAPIListParams>,
    ): Event[] => {
      return data.pages.flatMap((page) => page.data.items ?? []);
    },
    getNextPageParam: (lastPage: AxiosResponse<NotificationsAPIList200>) => {
      if (lastPage.data.pageInfo.endCursor === '') {
        return undefined;
      }
      return { limit: initialParams.limit, offset: lastPage.data.pageInfo.endCursor };
    },
    gcTime: 10 * 60 * 1000,
  };
}

export function refetchNotifications(queryClient: QueryClient) {
  queryClient.refetchQueries(getNotificationsQueryOptions());
}
