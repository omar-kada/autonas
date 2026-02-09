import { useRegister } from '@/hooks';
import { useTranslation } from 'react-i18next';
import { RegisterForm } from './login';
import { Card, CardContent } from './ui/card';

export function RegisterPage() {
  const { t } = useTranslation();

  const { register, isPending } = useRegister();

  return (
    <div className="p-4 space-y-4 h-full flex items-center flex-col justify-center">
      <h2 className="text-xl">{t('REGISTER.FORM.TITLE')}</h2>
      <Card>
        <CardContent className="space-y-4">
          <RegisterForm
            onSubmit={register}
            className="m-auto max-w-lg flex-1"
            loading={isPending}
          ></RegisterForm>
        </CardContent>
      </Card>
    </div>
  );
}
