import { Button } from '@/components/ui/button';
import { Card, CardAction, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ServiceLogo } from '@/lib';
import { Trash2 } from 'lucide-react';
import { Controller, type FieldArrayWithId, type UseFormReturn } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import { Field, FieldError, FieldGroup } from '../ui/field';
import { Input } from '../ui/input';
import { type FormValues } from './config-form-schema';
import { EnvVarArrayForm } from './env-vars-array-form';

export function ServiceCard({
  form,
  service,
  name,
  disabled,
  onRemove,
}: {
  form: UseFormReturn<FormValues>;
  service: FieldArrayWithId<FormValues, 'services', 'id'>;
  name: `services.${number}`;
  disabled?: boolean;
  onRemove: () => void;
}) {
  const { t } = useTranslation();
  return (
    <Card key={service.id}>
      <CardHeader>
        <CardTitle className="flex items-start gap-2">
          <ServiceNameForm service={service} name={name} form={form} />
        </CardTitle>
        <CardAction>
          {!disabled && (
            <Button
              aria-label={t('CONFIGURATION.FORM.REMOVE_SERVICE')}
              variant="ghost"
              disabled={disabled}
              onClick={onRemove}
            >
              <Trash2 />
            </Button>
          )}
        </CardAction>
      </CardHeader>

      <CardContent>
        <div className="mt-4">
          <EnvVarArrayForm
            control={form.control}
            name={`${name}.envVars`}
            disabled={disabled}
          ></EnvVarArrayForm>
        </div>
      </CardContent>
    </Card>
  );
}
function ServiceNameForm({
  service,
  name,
  form,
}: {
  service: FieldArrayWithId<FormValues, 'services', 'id'>;
  name: `services.${number}`;
  form: UseFormReturn<FormValues>;
}) {
  const { t } = useTranslation();
  return (
    <>
      <ServiceLogo service={service.name ?? ''} className="size-6 m-2" />
      <FieldGroup className="max-w-sm">
        <Controller
          name={`${name}.name`}
          control={form.control}
          render={({ field, fieldState }) => (
            <Field data-invalid={fieldState.invalid}>
              <Input
                {...field}
                aria-invalid={fieldState.invalid}
                autoComplete="off"
                placeholder={t('CONFIGURATION.FORM.SERVICE_NAME')}
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
    </>
  );
}
