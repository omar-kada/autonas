import {
  deployementAPIList,
  getDeployementAPIListQueryKey,
  type DeployementAPIList200,
  type DeployementAPIListParams,
  type Deployment,
} from '@/api/api';
import { type InfiniteData, type QueryClient } from '@tanstack/react-query';
import type { AxiosResponse } from 'axios';

const initialParams = { limit: 10, offset: '' } as DeployementAPIListParams;

export function getDeploymentsQueryOptions() {
  return {
    queryKey: getDeployementAPIListQueryKey(initialParams),
    queryFn: ({ pageParam = initialParams }: { pageParam: DeployementAPIListParams }) =>
      deployementAPIList(pageParam),
    initialPageParam: initialParams,
    select: (
      data: InfiniteData<AxiosResponse<DeployementAPIList200>, DeployementAPIListParams>,
    ): Deployment[] => {
      return data.pages.flatMap((page) => page.data.items ?? []);
    },
    getNextPageParam: (lastPage: AxiosResponse<DeployementAPIList200>) => {
      if (lastPage.data.pageInfo.endCursor === '') {
        return undefined;
      }
      return { limit: initialParams.limit, offset: lastPage.data.pageInfo.endCursor };
    },
  };
}

export function refetchDeployments(queryClient: QueryClient) {
  queryClient.refetchQueries(getDeploymentsQueryOptions());
}
