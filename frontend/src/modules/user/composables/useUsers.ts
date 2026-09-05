import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { tokenService } from '@/modules/token/services/TokenService';
import { userService } from '@/modules/user/services/UserService';
import type { UserWithToken } from '@/modules/user/types';
import { ApiError } from '@/shared/http/ApiError';

export function useUsers() {
  const { t } = useI18n();
  const users = ref<UserWithToken[]>([]);
  const loading = ref(false);
  const mutatingId = ref<string | null>(null);
  const error = ref<string | null>(null);

  function toErrorMessage(err: unknown, fallback: string): string {
    if (err instanceof ApiError) {
      return err.message;
    }
    return fallback;
  }

  async function fetchUsers(): Promise<void> {
    error.value = null;
    loading.value = true;
    try {
      users.value = await userService.list();
    } catch (err) {
      error.value = toErrorMessage(err, t('users.errors.loadFailed'));
    } finally {
      loading.value = false;
    }
  }

  async function revokeToken(userId: string): Promise<void> {
    error.value = null;
    mutatingId.value = userId;
    try {
      await tokenService.revokeForUser(userId);
      await fetchUsers();
    } catch (err) {
      error.value = toErrorMessage(err, t('users.errors.revokeFailed'));
    } finally {
      mutatingId.value = null;
    }
  }

  return {
    users,
    loading,
    mutatingId,
    error,
    fetchUsers,
    revokeToken,
  };
}
