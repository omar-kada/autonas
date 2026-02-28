import {
  getNotificationsAPIListQueryKey,
  notificationsAPIList,
  type Error,
  type Event,
  type NotificationsAPIList200,
  type NotificationsAPIListParams,
} from '@/api/api';
import {
  type InfiniteData,
  type QueryClient,
  type QueryKey,
  type UseInfiniteQueryOptions,
} from '@tanstack/react-query';
import type { AxiosError, AxiosResponse } from 'axios';

const initialParams = { limit: 10, offset: '' } as NotificationsAPIListParams;

export function getNotificationsQueryOptions(): UseInfiniteQueryOptions<
  AxiosResponse<NotificationsAPIList200>,
  AxiosError<Error>,
  Event[],
  QueryKey,
  NotificationsAPIListParams
> {
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
    refetchOnMount: true,
  };
}

export function refetchNotifications(queryClient: QueryClient) {
  queryClient.refetchQueries(getNotificationsQueryOptions());
}
