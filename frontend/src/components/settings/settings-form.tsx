import { GitBranch, Timer } from 'lucide-react';
import { Controller, type UseFormReturn } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLegend,
  FieldSeparator,
  FieldSet,
  FieldTitle,
} from '../ui/field';
import { Input } from '../ui/input';
import { type FormValues } from './settings-form-schema';

export function SettingsForm({ form }: { form: UseFormReturn<FormValues> }) {
  const { t } = useTranslation();

  return (
    <form>
      <FieldGroup>
        <FieldSet>
          <FieldLegend className="flex items-center gap-2">
            <GitBranch className="size-4" />
            {t('SETTINGS.FORM.GIT')}
          </FieldLegend>

          <SettingsField form={form} name="repo" />
          <SettingsField form={form} name="branch" />

          <SettingsField form={form} name="username" />
          <SettingsField form={form} name="token" />
        </FieldSet>
        <FieldSeparator />

        <FieldSet>
          <FieldLegend className="flex items-center gap-2">
            <Timer className="size-4" />
            {t('SETTINGS.FORM.AUTO_SYNC')}
          </FieldLegend>

          <SettingsField form={form} name="cron" withDescription />
        </FieldSet>
      </FieldGroup>
    </form>
  );
}
function SettingsField({
  form,
  name,
  withDescription = false,
  withPlaceholder = true,
}: {
  form: UseFormReturn<FormValues>;
  name: keyof FormValues;
  withDescription?: boolean;
  withPlaceholder?: boolean;
}) {
  const { t } = useTranslation();
  return (
    <Controller
      name={name}
      control={form.control}
      render={({ field, fieldState }) => (
        <Field data-invalid={fieldState.invalid}>
          <FieldTitle>{t(`SETTINGS.FORM.${name}`)}</FieldTitle>
          {withDescription && (
            <FieldDescription>{t(`SETTINGS.FORM.${name}_DESCRIPTION`)}</FieldDescription>
          )}
          <Input
            {...field}
            aria-invalid={fieldState.invalid}
            autoComplete="off"
            placeholder={withPlaceholder ? t(`SETTINGS.FORM.${name}_PLACEHOLDER`) : ''}
          />
          {fieldState.invalid && (
            <FieldError
              errors={[{ ...fieldState.error, message: t(fieldState.error?.message ?? '') }]}
            />
          )}
        </Field>
      )}
    />
  );
}
