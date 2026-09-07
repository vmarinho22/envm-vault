import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useAuthStore } from '@/modules/auth/stores/auth.store';
import { teamService } from '@/modules/team/services/TeamService';
import type { CreateTeamInput, Team, UpdateTeamInput } from '@/modules/team/types';
import { userService } from '@/modules/user/services/UserService';
import type { UserWithToken } from '@/modules/user/types';
import { useConfirmDialog } from '@/shared/confirm/useConfirmDialog';
import { ApiError } from '@/shared/http/ApiError';

export function useTeams() {
  const { t } = useI18n();
  const auth = useAuthStore();
  const { confirm } = useConfirmDialog();

  const teams = ref<Team[]>([]);
  const users = ref<UserWithToken[]>([]);
  const loading = ref(false);
  const mutating = ref(false);
  const error = ref<string | null>(null);
  const nameQuery = ref('');
  const appliedName = ref('');

  const formOpen = ref(false);
  const editing = ref<Team | null>(null);
  const formName = ref('');
  const formAccountableId = ref('');

  const canMutate = computed(() => auth.isAdmin);
  const hasActiveSearch = computed(() => appliedName.value.trim() !== '');

  function toErrorMessage(err: unknown, fallback: string): string {
    if (err instanceof ApiError) {
      return err.message;
    }
    return fallback;
  }

  async function fetchTeams(): Promise<void> {
    error.value = null;
    loading.value = true;
    try {
      const name = appliedName.value.trim();
      teams.value = await teamService.list(name ? { name } : undefined);
    } catch (err) {
      error.value = toErrorMessage(err, t('teams.errors.loadFailed'));
    } finally {
      loading.value = false;
    }
  }

  async function fetchUsersIfAdmin(): Promise<void> {
    if (!auth.isAdmin) return;
    try {
      users.value = await userService.list();
    } catch (err) {
      error.value = toErrorMessage(err, t('teams.errors.loadUsersFailed'));
    }
  }

  async function applySearch(): Promise<void> {
    appliedName.value = nameQuery.value.trim();
    await fetchTeams();
  }

  async function clearSearch(): Promise<void> {
    nameQuery.value = '';
    appliedName.value = '';
    await fetchTeams();
  }

  function openCreate(): void {
    editing.value = null;
    formName.value = '';
    formAccountableId.value = users.value[0]?.id ?? '';
    formOpen.value = true;
  }

  function openEdit(team: Team): void {
    editing.value = team;
    formName.value = team.name;
    formAccountableId.value = team.accountable.id;
    formOpen.value = true;
  }

  function closeForm(): void {
    formOpen.value = false;
    editing.value = null;
  }

  async function submitForm(): Promise<void> {
    if (!canMutate.value) return;
    error.value = null;
    mutating.value = true;
    try {
      if (editing.value) {
        const input: UpdateTeamInput = {
          name: formName.value,
          accountable_id: formAccountableId.value,
        };
        await teamService.update(editing.value.id, input);
      } else {
        const input: CreateTeamInput = {
          name: formName.value,
          accountable_id: formAccountableId.value,
        };
        await teamService.create(input);
      }
      closeForm();
      await fetchTeams();
    } catch (err) {
      error.value = toErrorMessage(
        err,
        editing.value ? t('teams.errors.updateFailed') : t('teams.errors.createFailed'),
      );
    } finally {
      mutating.value = false;
    }
  }

  async function removeTeam(team: Team): Promise<void> {
    if (!canMutate.value) return;
    const ok = await confirm({
      title: t('teams.delete'),
      message: t('teams.confirmDelete', { name: team.name }),
      confirmLabel: t('teams.delete'),
      cancelLabel: t('common.cancel'),
      variant: 'danger',
    });
    if (!ok) return;

    error.value = null;
    mutating.value = true;
    try {
      await teamService.remove(team.id);
      await fetchTeams();
    } catch (err) {
      error.value = toErrorMessage(err, t('teams.errors.deleteFailed'));
    } finally {
      mutating.value = false;
    }
  }

  async function init(): Promise<void> {
    await Promise.all([fetchTeams(), fetchUsersIfAdmin()]);
  }

  return {
    teams,
    users,
    loading,
    mutating,
    error,
    nameQuery,
    appliedName,
    hasActiveSearch,
    formOpen,
    editing,
    formName,
    formAccountableId,
    canMutate,
    init,
    fetchTeams,
    applySearch,
    clearSearch,
    openCreate,
    openEdit,
    closeForm,
    submitForm,
    removeTeam,
  };
}
