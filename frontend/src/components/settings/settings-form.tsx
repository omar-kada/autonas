import { Controller, type UseFormReturn } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import { Field, FieldError, FieldGroup, FieldLabel } from '../ui/field';
import { Input } from '../ui/input';
import { type FormValues } from './settings-form-schema';

export function SettingsForm({ form }: { form: UseFormReturn<FormValues> }) {
  const { t } = useTranslation();

  return (
    <form>
      <FieldGroup>
        <Controller
          name="repo"
          control={form.control}
          render={({ field, fieldState }) => (
            <Field data-invalid={fieldState.invalid}>
              <FieldLabel>{t('SETTINGS.FORM.REPO')}</FieldLabel>
              <Input
                {...field}
                aria-invalid={fieldState.invalid}
                autoComplete="off"
                placeholder={t('SETTINGS.FORM.REPO_PLACEHOLDER')}
              />
              {fieldState.invalid && (
                <FieldError
                  errors={[{ ...fieldState.error, message: t(fieldState.error?.message ?? '') }]}
                />
              )}
            </Field>
          )}
        />
        <Controller
          name="branch"
          control={form.control}
          render={({ field, fieldState }) => (
            <Field data-invalid={fieldState.invalid}>
              <FieldLabel>{t('SETTINGS.FORM.BRANCH')}</FieldLabel>
              <Input
                {...field}
                aria-invalid={fieldState.invalid}
                autoComplete="off"
                placeholder={t('SETTINGS.FORM.BRANCH_PLACEHOLDER')}
              />
              {fieldState.invalid && (
                <FieldError
                  errors={[{ ...fieldState.error, message: t(fieldState.error?.message ?? '') }]}
                />
              )}
            </Field>
          )}
        />
        <Controller
          name="cron"
          control={form.control}
          render={({ field, fieldState }) => (
            <Field data-invalid={fieldState.invalid}>
              <FieldLabel>{t('SETTINGS.FORM.CRON')}</FieldLabel>
              <Input
                {...field}
                aria-invalid={fieldState.invalid}
                autoComplete="off"
                placeholder={t('SETTINGS.FORM.CRON_PLACEHOLDER')}
              />
              {fieldState.invalid && (
                <FieldError
                  errors={[{ ...fieldState.error, message: t(fieldState.error?.message ?? '') }]}
                />
              )}
            </Field>
          )}
        />
      </FieldGroup>
    </form>
  );
}
