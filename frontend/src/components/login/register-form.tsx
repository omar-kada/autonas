import type { Credentials } from '@/api/api';
import { Button } from '@/components/ui/button';
import { InputGroup, InputGroupAddon, InputGroupInput } from '@/components/ui/input-group'; // Add this import at the top of the file
import { cn } from '@/lib';
import { zodResolver } from '@hookform/resolvers/zod';
import { EyeIcon, EyeOffIcon, Lock, User as UserIcon } from 'lucide-react';
import { useState } from 'react';
import { Controller, useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import * as z from 'zod';
import { Field, FieldError, FieldLabel, FieldSet } from '../ui/field';
import { Spinner } from '../ui/spinner';

const registerSchema = z.object({
  username: z
    .string({
      error: 'REGISTER.FORM.username_MIN_3',
    })
    .min(3, { error: 'REGISTER.FORM.username_MIN_3' }),
  password: z
    .string({ error: 'REGISTER.FORM.password_MIN_12' })
    .min(12, { error: 'REGISTER.FORM.password_MIN_12' }),
});

type RegisterFormValues = z.infer<typeof registerSchema>;

export function RegisterForm({
  className,
  onSubmit,
  loading,
}: {
  className?: string;
  onSubmit: (data: Credentials) => void;
  loading: boolean;
}) {
  const { t } = useTranslation();
  const { handleSubmit, control } = useForm<RegisterFormValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      username: '',
      password: '',
    },
  });
  const [showPassword, setShowPassword] = useState(false);

  return (
    <form className={cn('space-y-4', className)} onSubmit={handleSubmit(onSubmit)}>
      <FieldSet>
        <Controller
          name={`username`}
          control={control}
          render={({ field, fieldState }) => (
            <Field data-invalid={fieldState.invalid}>
              <FieldLabel>{t('REGISTER.FORM.username')}</FieldLabel>
              <InputGroup>
                <InputGroupInput {...field} autoComplete="off" aria-invalid={fieldState.invalid} />
                <InputGroupAddon align="inline-start">
                  <UserIcon />
                </InputGroupAddon>
              </InputGroup>
              {fieldState.invalid && (
                <FieldError
                  errors={[{ ...fieldState.error, message: t(fieldState.error?.message ?? '') }]}
                />
              )}
            </Field>
          )}
        />

        <Controller
          name={`password`}
          control={control}
          render={({ field, fieldState }) => (
            <Field data-invalid={fieldState.invalid}>
              <FieldLabel>{t('REGISTER.FORM.password')}</FieldLabel>
              <InputGroup>
                <InputGroupInput
                  {...field}
                  type={showPassword ? 'text' : 'password'}
                  autoComplete="off"
                  aria-invalid={fieldState.invalid}
                />
                <InputGroupAddon align="inline-start">
                  <Lock />
                </InputGroupAddon>
                <InputGroupAddon align="inline-end">
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    onClick={() => setShowPassword(!showPassword)}
                    aria-label={showPassword ? 'Hide password' : 'Show password'}
                  >
                    {showPassword ? <EyeOffIcon /> : <EyeIcon />}
                  </Button>
                </InputGroupAddon>
              </InputGroup>
              {fieldState.invalid && (
                <FieldError
                  errors={[{ ...fieldState.error, message: t(fieldState.error?.message ?? '') }]}
                />
              )}
            </Field>
          )}
        />
      </FieldSet>

      <div className="flex justify-between items-center">
        <Button type="submit" disabled={loading}>
          {t('REGISTER.FORM.SUBMIT')}
        </Button>
        {loading && <Spinner />}
      </div>
    </form>
  );
}
