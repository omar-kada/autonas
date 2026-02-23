import { EventType, type Settings } from '@/api/api';
import z from 'zod/v3';

export const formSchema = z.object({
  repo: z.string().min(1, { message: 'SETTINGS.FORM.REPO_REQUIRED' }),
  branch: z.string().optional(),
  cron: z.string().optional(),
  username: z.string().optional(),
  token: z.string().optional(),
  notificationURL: z.string().optional(),
  notificationTypes: z.array(z.nativeEnum(EventType)),
});
export type FormValues = z.infer<typeof formSchema>;

export function fromSettings(settings?: Settings): FormValues {
  if (!settings) {
    return {
      repo: '',
      notificationTypes: [],
    };
  }
  return settings;
}

export function toSettings(formValues: FormValues): Settings {
  return formValues;
}
