import { useLogin } from '@/hooks';
import { useTranslation } from 'react-i18next';
import { LoginForm } from './login';
import { Card, CardContent } from './ui/card';

export function LoginPage() {
  const { t } = useTranslation();

  const { login, isPending } = useLogin();

  return (
    <div className="p-4 space-y-4 h-full flex items-center flex-col justify-center">
      <h2 className="text-xl">{t('LOGIN.FORM.TITLE')}</h2>
      <Card>
        <CardContent className="space-y-4">
          <LoginForm
            onSubmit={login}
            className="m-auto max-w-lg flex-1"
            loading={isPending}
          ></LoginForm>
        </CardContent>
      </Card>
    </div>
  );
}
