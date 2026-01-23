import { ScrollArea } from '@/components/ui/scroll-area';
import { getConfigQueryOptions, getFeaturesQueryOptions, useUpdateConfig } from '@/hooks';
import { zodResolver } from '@hookform/resolvers/zod';
import { useQuery } from '@tanstack/react-query';
import { RotateCcw, Save } from 'lucide-react';
import { useCallback, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import { ConfigForm, formSchema, fromConfig, toConfig, type FormValues } from './config';
import { Button } from './ui/button';
import { Skeleton } from './ui/skeleton';
import { Spinner } from './ui/spinner';
import { ErrorAlert, HeaderLayout, InfoEmpty } from './view';

export function ConfigPage() {
  const { t } = useTranslation();
  const {
    data: features,
    isPending: featuresPending,
    error: featuresError,
  } = useQuery(getFeaturesQueryOptions());
  const {
    data: config,
    isPending,
    error,
  } = useQuery(
    getConfigQueryOptions({
      enabled: !!features?.displayConfig,
    }),
  );
  const { updateConfig, isPending: isSavePending } = useUpdateConfig();
  const onSubmit = (data: FormValues) => updateConfig(toConfig(data));

  const disabled = !features?.editConfig;
  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    disabled,
  });
  const resetForm = useCallback(() => form.reset(fromConfig(config)), [form, config]);

  useEffect(() => {
    // init and reset when config changes
    if (config) {
      form.reset(fromConfig(config));
    }
  }, [config, form]);

  if (featuresPending || (features?.displayConfig && isPending)) {
    return (
      <div className="flex flex-col gap-10 m-4">
        <EnvVarSkeleton repeat={2} />
        <EnvVarSkeleton repeat={1} />
      </div>
    );
  }
  if (features && !features.displayConfig) {
    return (
      <InfoEmpty
        title="CONFIGURATION.DISABLED_TITLE"
        details="CONFIGURATION.DISABLED_DESCRIPTION"
      ></InfoEmpty>
    );
  }

  const mergedError = featuresError || error;

  return (
    <HeaderLayout
      header={
        !disabled && (
          <div className="flex w-full justify-end-safe gap-2">
            {form.formState.isDirty && (
              <Button variant="outline" onClick={resetForm}>
                <RotateCcw />
                {t('ACTION.RESET')}
              </Button>
            )}
            {!disabled && (
              <Button onClick={form.handleSubmit(onSubmit)} disabled={!form.formState.isDirty}>
                {isSavePending ? <Spinner /> : <Save />}
                {t('ACTION.SAVE')}
              </Button>
            )}
          </div>
        )
      }
    >
      <ErrorAlert
        title={mergedError && 'ALERT.LOAD_CONFIGURATION_ERROR'}
        details={mergedError?.message}
        className="m-4"
      />
      <ScrollArea className="h-full flex-1">
        {config && (
          <ConfigForm className="flex-1 w-full p-4" form={form} disabled={disabled}></ConfigForm>
        )}
      </ScrollArea>
    </HeaderLayout>
  );
}

function EnvVarSkeleton({ repeat }: { repeat: number }) {
  return (
    <div className="flex flex-col gap-4">
      <Skeleton className="h-6 w-25"></Skeleton>
      {Array(repeat)
        .fill({})
        .map((_, index) => (
          <span className="flex gap-4" key={index}>
            <Skeleton className="h-6 w-50" />
            <Skeleton className="h-6 w-50" />
          </span>
        ))}
    </div>
  );
}
