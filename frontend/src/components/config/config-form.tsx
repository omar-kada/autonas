import { Button } from '@/components/ui/button';
import { Card, CardAction, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { cn, ServiceLogo } from '@/lib';
import { Plus, Trash2 } from 'lucide-react';
import { Controller, useFieldArray, type UseFormReturn } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import { Field, FieldError, FieldGroup } from '../ui/field';
import { Input } from '../ui/input';
import { type FormValues } from './config-form-schema';
import { EnvVarArrayForm } from './env-vars-array-form';

export function ConfigForm({
  form,
  className,
  disabled,
}: {
  form: UseFormReturn<FormValues>;
  className?: string;
  disabled?: boolean;
}) {
  const { t } = useTranslation();

  const {
    fields: servicesFields,
    append: appendService,
    remove: removeService,
  } = useFieldArray({
    control: form.control,
    name: 'services',
  });

  return (
    <form className={cn('space-y-4', className)}>
      <h2 className="text-xl">{t('CONFIGURATION.FORM.GLOBAL_ENV_VARS')}</h2>
      <Card>
        <CardContent>
          <EnvVarArrayForm
            control={form.control}
            name="globalEnvVars"
            disabled={disabled}
          ></EnvVarArrayForm>
        </CardContent>
      </Card>
      <h2 className="text-xl">{t('CONFIGURATION.FORM.SERVICES')}</h2>

      {servicesFields.map((service, serviceIndex) => (
        <Card key={service.id}>
          <CardHeader>
            <CardTitle>
              <div className="flex items-start gap-2">
                <ServiceLogo service={service.name ?? ''} className="size-6 m-2" />
                <FieldGroup className="max-w-sm">
                  <Controller
                    name={`services.${serviceIndex}.name`}
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
                            errors={[
                              { ...fieldState.error, message: t(fieldState.error?.message ?? '') },
                            ]}
                          />
                        )}
                      </Field>
                    )}
                  />
                </FieldGroup>
              </div>
            </CardTitle>
            <CardAction>
              {!disabled && (
                <Button
                  aria-label={t('CONFIGURATION.FORM.REMOVE_SERVICE')}
                  variant="ghost"
                  disabled={disabled}
                  onClick={() => removeService(serviceIndex)}
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
                name={`services.${serviceIndex}.envVars`}
                disabled={disabled}
              ></EnvVarArrayForm>
            </div>
          </CardContent>
        </Card>
      ))}
      {!disabled && (
        <Button
          variant="default"
          type="button"
          onClick={() =>
            appendService({ name: '', envVars: [{ key: '', value: '' }] }, { shouldFocus: true })
          }
        >
          <Plus /> {t('CONFIGURATION.FORM.ADD_SERVICE')}
        </Button>
      )}
    </form>
  );
}
