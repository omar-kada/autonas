import { lazy, Suspense } from 'react';

export function Devtools() {
  if (!import.meta.env.DEV) return null;

  const LazyDevtools = lazy(() =>
    import('@tanstack/react-query-devtools').then((mod) => ({
      default: mod.ReactQueryDevtools,
    })),
  );

  return (
    <Suspense fallback={null}>
      <LazyDevtools initialIsOpen={false} />
    </Suspense>
  );
}
