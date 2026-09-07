import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { tokenService } from '@/modules/token/services/TokenService';
import type { RegisterableRole } from '@/modules/user/services/UserService';
import { userService } from '@/modules/user/services/UserService';
import type { UserWithToken } from '@/modules/user/types';
import { ApiError } from '@/shared/http/ApiError';

export function useUsers() {
  const { t } = useI18n();
  const users = ref<UserWithToken[]>([]);
  const loading = ref(false);
  const mutatingId = ref<string | null>(null);
  const creating = ref(false);
  const error = ref<string | null>(null);

  const formOpen = ref(false);
  const formEmail = ref('');
  const formRole = ref<RegisterableRole>('user');

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

  function openCreate(): void {
    formEmail.value = '';
    formRole.value = 'user';
    formOpen.value = true;
  }

  function closeForm(): void {
    formOpen.value = false;
  }

  async function submitCreate(): Promise<void> {
    error.value = null;
    creating.value = true;
    try {
      await userService.register(formEmail.value.trim(), formRole.value);
      closeForm();
      await fetchUsers();
    } catch (err) {
      error.value = toErrorMessage(err, t('users.errors.createFailed'));
    } finally {
      creating.value = false;
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
    creating,
    error,
    formOpen,
    formEmail,
    formRole,
    fetchUsers,
    openCreate,
    closeForm,
    submitCreate,
    revokeToken,
  };
}
