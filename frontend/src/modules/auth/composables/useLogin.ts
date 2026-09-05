import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { authService } from '@/modules/auth/services/AuthService';
import { useAuthStore } from '@/modules/auth/stores/auth.store';
import { ApiError } from '@/shared/http/ApiError';

export type LoginStep = 'email' | 'code';

export function useLogin() {
  const router = useRouter();
  const auth = useAuthStore();
  const { t } = useI18n();

  const step = ref<LoginStep>('email');
  const email = ref('');
  const code = ref('');
  const loading = ref(false);
  const error = ref<string | null>(null);

  function toErrorMessage(err: unknown): string {
    if (err instanceof ApiError) {
      if (err.status === 429) {
        if (err.message.includes('verification')) {
          return t('login.errors.tooManyVerification');
        }
        return t('login.errors.tooManyRequests');
      }
      return err.message;
    }
    return t('login.errors.fallback');
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
      const hasToken = await auth.ensureHasApiToken();
      await router.push({ name: hasToken ? 'home' : 'setup-token' });
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
