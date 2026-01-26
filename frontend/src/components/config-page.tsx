import { ScrollArea } from '@/components/ui/scroll-area';
import {
  getConfigQueryOptions,
  getFeaturesQueryOptions,
  useIsMobile,
  useUpdateConfig,
} from '@/hooks';
import { cn } from '@/lib';
import { zodResolver } from '@hookform/resolvers/zod';
import { useQuery } from '@tanstack/react-query';
import { AlertCircleIcon, Code, RotateCcw, Save } from 'lucide-react';
import { useCallback, useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import {
  ConfigForm,
  ConfigViewer,
  formSchema,
  fromConfig,
  toConfig,
  toYaml,
  type FormValues,
} from './config';
import { Alert, AlertDescription, AlertTitle } from './ui/alert';
import { Button } from './ui/button';
import { ResizableHandle, ResizablePanel, ResizablePanelGroup } from './ui/resizable';
import { Skeleton } from './ui/skeleton';
import { Spinner } from './ui/spinner';
import { Toggle } from './ui/toggle';
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
  const isMobile = useIsMobile();
  const [showYaml, setShowYaml] = useState(!isMobile);
  const hideFile = useCallback(() => setShowYaml(false), [setShowYaml]);

  useEffect(() => {
    if (showYaml && isMobile) {
      setShowYaml(false);
    }
  }, [isMobile]);

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
        <div className="flex w-full justify-between gap-2">
          <div>
            <Toggle
              pressed={showYaml}
              onPressedChange={setShowYaml}
              aria-label={showYaml ? t('ACTION.HIDE_FILE') : t('ACTION.VIEW_FILE')}
              variant="outline"
            >
              <Code />
              {showYaml ? t('ACTION.HIDE_FILE') : t('ACTION.VIEW_FILE')}
            </Toggle>
          </div>
          {!disabled && (
            <div className="flex">
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
          )}
        </div>
      }
    >
      <ErrorAlert
        title={mergedError && 'ALERT.LOAD_CONFIGURATION_ERROR'}
        details={mergedError?.message}
        className="m-4"
      />
      <ScrollArea className="h-full flex-1">
        {disabled && (
          <Alert variant="default" className="w-auto m-4">
            <AlertCircleIcon />
            <AlertTitle>{t('CONFIGURATION.DISABLED_EDIT_TITLE')}</AlertTitle>
            <AlertDescription>{t('CONFIGURATION.DISABLED_EDIT_DESCRIPTION')}</AlertDescription>
          </Alert>
        )}
        <ResizablePanelGroup direction="horizontal">
          <ResizablePanel hidden={isMobile && showYaml} className="min-w-1/3">
            {config && (
              <ConfigForm className={cn('flex-1 p-4')} form={form} disabled={disabled}></ConfigForm>
            )}
          </ResizablePanel>
          <ResizableHandle withHandle={!isMobile} />
          <ResizablePanel hidden={!showYaml} className={showYaml ? 'min-w-1/3' : 'w-0 min-w-0'}>
            <ConfigViewer
              className={cn(
                'transform transition-all',
                showYaml ? 'translate-x-0' : 'translate-x-full w-0',
              )}
              text={toYaml(form.getValues())}
              onClose={hideFile}
            />
          </ResizablePanel>
        </ResizablePanelGroup>
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
