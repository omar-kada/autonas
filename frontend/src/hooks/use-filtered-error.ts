import type { Error } from '@/api/api';
import { isInvalidToken } from '@/lib';
import {
  useQuery,
  type QueryKey,
  type QueryObserverLoadingResult,
  type UseQueryOptions,
  type UseQueryResult,
} from '@tanstack/react-query';
import type { AxiosError } from 'axios';
import { useEffect, useState } from 'react';

function useFilteredError<D>(
  query: UseQueryResult<D, AxiosError<Error>>,
): UseQueryResult<D, AxiosError<Error> | null> {
  const [res, setRes] = useState<UseQueryResult<D, AxiosError<Error>>>(query);

  useEffect(() => {
    if (isInvalidToken(query.error)) {
      setRes({
        ...query,
        error: null,
        isPending: true,
      } as QueryObserverLoadingResult<D, AxiosError<Error>>);
    } else {
      setRes(query);
    }
  }, [query.error, query.data, query.error, setRes]);

  return res;
}

export function useFilteredQuery<FnD, TData = FnD, TQueryKey extends QueryKey = QueryKey>(
  options: UseQueryOptions<FnD, AxiosError<Error>, TData, TQueryKey>,
) {
  return useFilteredError(useQuery(options));
}
