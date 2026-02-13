import { getUserAPIChangePasswordMutationOptions } from '@/api/api';
import { useMutation } from '@tanstack/react-query';

export const useChangePass = () => {
  const changePassMutation = useMutation(
    getUserAPIChangePasswordMutationOptions({
      mutation: {
        onSuccess(data) {
          return data.data?.success;
        },
      },
    }),
  );

  return {
    ...changePassMutation,
    changePass: changePassMutation.mutateAsync,
  };
};
