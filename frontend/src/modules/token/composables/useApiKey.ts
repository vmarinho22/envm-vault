import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/modules/auth/stores/auth.store';
import { tokenService } from '@/modules/token/services/TokenService';
import type { TokenStatus } from '@/modules/token/types';
import { ApiError } from '@/shared/http/ApiError';

export function useApiKey() {
  const router = useRouter();
  const auth = useAuthStore();

  const loading = ref(false);
  const mutating = ref(false);
  const error = ref<string | null>(null);
  const status = ref<TokenStatus | null>(null);
  const plaintext = ref<string | null>(null);

  const prefix = computed(() => status.value?.prefix ?? null);
  const createdAt = computed(() => status.value?.created_at ?? null);
  const isActive = computed(() => status.value?.active === true);

  function toErrorMessage(err: unknown, fallback: string): string {
    if (err instanceof ApiError) {
      return err.message;
    }
    return fallback;
  }

  async function load(): Promise<void> {
    error.value = null;
    loading.value = true;
    try {
      status.value = await tokenService.getMine();
      auth.setHasApiToken(status.value.active);
    } catch (err) {
      error.value = toErrorMessage(err, 'não foi possível carregar a chave');
    } finally {
      loading.value = false;
    }
  }

  async function rotate(): Promise<void> {
    if (
      !window.confirm(
        'Rotacionar gera uma chave nova e invalida a atual. Quem ainda usa a chave antiga vai falhar. Continuar?',
      )
    ) {
      return;
    }

    error.value = null;
    mutating.value = true;
    try {
      const created = await tokenService.create();
      plaintext.value = created.token;
      status.value = {
        active: true,
        prefix: created.prefix,
        created_at: created.created_at,
      };
      auth.setHasApiToken(true);
    } catch (err) {
      error.value = toErrorMessage(err, 'não foi possível rotacionar a chave');
    } finally {
      mutating.value = false;
    }
  }

  function dismissPlaintext(): void {
    plaintext.value = null;
  }

  async function revoke(): Promise<void> {
    error.value = null;
    mutating.value = true;
    try {
      await tokenService.revokeOwn();
      auth.setHasApiToken(false);
      await router.push({ name: 'setup-token' });
    } catch (err) {
      error.value = toErrorMessage(err, 'não foi possível revogar a chave');
    } finally {
      mutating.value = false;
    }
  }

  return {
    loading,
    mutating,
    error,
    status,
    prefix,
    createdAt,
    isActive,
    plaintext,
    load,
    rotate,
    revoke,
    dismissPlaintext,
  };
}
