import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/modules/auth/stores/auth.store';
import { tokenService } from '@/modules/token/services/TokenService';
import { ApiError } from '@/shared/http/ApiError';

export function useTokenSetup() {
  const router = useRouter();
  const auth = useAuthStore();

  const loading = ref(false);
  const error = ref<string | null>(null);
  const plaintext = ref<string | null>(null);

  function toErrorMessage(err: unknown): string {
    if (err instanceof ApiError) {
      return err.message;
    }
    return 'não foi possível gerar a chave';
  }

  async function generate(): Promise<void> {
    error.value = null;
    loading.value = true;
    try {
      const created = await tokenService.create();
      plaintext.value = created.token;
      auth.setHasApiToken(true);
    } catch (err) {
      error.value = toErrorMessage(err);
    } finally {
      loading.value = false;
    }
  }

  async function continueToApp(): Promise<void> {
    await router.push({ name: 'home' });
  }

  return {
    loading,
    error,
    plaintext,
    generate,
    continueToApp,
  };
}
