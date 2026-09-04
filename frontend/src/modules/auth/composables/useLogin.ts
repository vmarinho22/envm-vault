import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { authService } from '@/modules/auth/services/AuthService';
import { useAuthStore } from '@/modules/auth/stores/auth.store';
import { ApiError } from '@/shared/http/ApiError';

export type LoginStep = 'email' | 'code';

export function useLogin() {
  const router = useRouter();
  const auth = useAuthStore();

  const step = ref<LoginStep>('email');
  const email = ref('');
  const code = ref('');
  const loading = ref(false);
  const error = ref<string | null>(null);

  function toErrorMessage(err: unknown): string {
    if (err instanceof ApiError) {
      if (err.status === 429) {
        if (err.message.includes('verification')) {
          return 'Muitas tentativas. Solicite um novo código.';
        }
        return 'Muitas requisições. Tente novamente em instantes.';
      }
      return err.message;
    }
    return 'não foi possível concluir o login';
  }

  async function submitEmail(): Promise<void> {
    error.value = null;
    loading.value = true;
    try {
      await authService.requestLogin(email.value.trim());
      code.value = '';
      step.value = 'code';
    } catch (err) {
      error.value = toErrorMessage(err);
    } finally {
      loading.value = false;
    }
  }

  async function submitCode(): Promise<void> {
    error.value = null;
    loading.value = true;
    try {
      const result = await authService.verifyLogin(email.value.trim(), code.value.trim());
      auth.setSession({ token: result.token, user: result.user });
      await router.push({ name: 'home' });
    } catch (err) {
      error.value = toErrorMessage(err);
      if (
        err instanceof ApiError &&
        err.status === 429 &&
        err.message.includes('verification')
      ) {
        code.value = '';
      }
    } finally {
      loading.value = false;
    }
  }

  function backToEmail(): void {
    step.value = 'email';
    code.value = '';
    error.value = null;
  }

  return {
    step,
    email,
    code,
    loading,
    error,
    submitEmail,
    submitCode,
    backToEmail,
  };
}
