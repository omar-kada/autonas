import { Button } from '@/components/ui/button';
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetDescription,
  SheetFooter,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet';
import { getFeaturesQueryOptions, getSettingsQueryOptions } from '@/hooks';
import { useUpdateSettings } from '@/hooks/use-update-settings';
import { zodResolver } from '@hookform/resolvers/zod';
import { useQuery } from '@tanstack/react-query';
import { useCallback, useEffect, useState, type ReactNode } from 'react';
import { useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import { ScrollArea } from '../ui/scroll-area';
import { Skeleton } from '../ui/skeleton';
import { ErrorAlert } from '../view';
import { SettingsForm } from './settings-form';
import { formSchema, fromSettings, toSettings, type FormValues } from './settings-form-schema';

export function SettingsSheet({ children }: { children: ReactNode }) {
  const { t } = useTranslation();
  const { data: features, error: featuresError } = useQuery(getFeaturesQueryOptions());
  const disabled = !features?.editSettings;
  const { data: settings, isPending, error } = useQuery(getSettingsQueryOptions());
  const mergedError = error ?? featuresError;

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    disabled,
  });
  useEffect(() => {
    // init once loaded
    form.reset(fromSettings(settings));
  }, [settings, form]);

  const { updateSettings } = useUpdateSettings();
  const [open, setOpen] = useState(false);

  const onSubmit = useCallback(() => {
    updateSettings(toSettings(form.getValues()));
    setOpen(false);
  }, [form, setOpen]);

  return (
    <Sheet open={open} onOpenChange={setOpen}>
      <SheetTrigger asChild>{children}</SheetTrigger>
      <SheetContent
        className="w-full md:w-none flex flex-col h-full"
        onOpenAutoFocus={(e) => e.preventDefault()}
      >
        <SheetHeader>
          <SheetTitle>{t('SETTINGS.SETTINGS')}</SheetTitle>
          {disabled && (
            <SheetDescription>
              {t('SETTINGS.DISABLED_TITLE')}, {t('SETTINGS.DISABLED_DESCRIPTION')}
            </SheetDescription>
          )}
        </SheetHeader>
        <ScrollArea className="h-1 flex-1">
          <div className=" grid auto-rows-min gap-6 px-4 mb-25">
            {mergedError && (
              <ErrorAlert
                title={t('ALERT.LOAD_SETTINGS_ERROR')}
                details={mergedError.message}
              ></ErrorAlert>
            )}
            {isPending ? <SettingsSkeleton /> : settings && <SettingsForm form={form} />}
          </div>
        </ScrollArea>
        <SheetFooter>
          <Button type="submit" onClick={form.handleSubmit(onSubmit)}>
            {t('ACTION.SAVE')}
          </Button>
          <SheetClose asChild>
            <Button variant="outline">{t('ACTION.CLOSE')}</Button>
          </SheetClose>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}

function SettingsSkeleton() {
  return (
    <div className="flex flex-col w-full">
      <Skeleton className="h-12"></Skeleton>
      <Skeleton className="h-12"></Skeleton>
      <Skeleton className="h-12"></Skeleton>
    </div>
  );
}
