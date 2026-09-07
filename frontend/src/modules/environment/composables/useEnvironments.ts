import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useAuthStore } from '@/modules/auth/stores/auth.store';
import { environmentService } from '@/modules/environment/services/EnvironmentService';
import type {
  CreateEnvironmentInput,
  Environment,
  UpdateEnvironmentInput,
} from '@/modules/environment/types';
import { useConfirmDialog } from '@/shared/confirm/useConfirmDialog';
import { ApiError } from '@/shared/http/ApiError';

export function useEnvironments() {
  const { t } = useI18n();
  const auth = useAuthStore();
  const { confirm } = useConfirmDialog();

  const environments = ref<Environment[]>([]);
  const loading = ref(false);
  const mutating = ref(false);
  const error = ref<string | null>(null);
  const titleQuery = ref('');
  const appliedTitle = ref('');

  const formOpen = ref(false);
  const editing = ref<Environment | null>(null);
  const formTitle = ref('');

  const canMutate = computed(() => auth.isAdmin);
  const hasActiveSearch = computed(() => appliedTitle.value.trim() !== '');

  function toErrorMessage(err: unknown, fallback: string): string {
    if (err instanceof ApiError) {
      return err.message;
    }
    return fallback;
  }

  async function fetchEnvironments(): Promise<void> {
    error.value = null;
    loading.value = true;
    try {
      const title = appliedTitle.value.trim();
      environments.value = await environmentService.list(title ? { title } : undefined);
    } catch (err) {
      error.value = toErrorMessage(err, t('environments.errors.loadFailed'));
    } finally {
      loading.value = false;
    }
  }

  async function applySearch(): Promise<void> {
    appliedTitle.value = titleQuery.value.trim();
    await fetchEnvironments();
  }

  async function clearSearch(): Promise<void> {
    titleQuery.value = '';
    appliedTitle.value = '';
    await fetchEnvironments();
  }

  function openCreate(): void {
    editing.value = null;
    formTitle.value = '';
    formOpen.value = true;
  }

  function openEdit(environment: Environment): void {
    editing.value = environment;
    formTitle.value = environment.title;
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
        const input: UpdateEnvironmentInput = {
          title: formTitle.value,
        };
        await environmentService.update(editing.value.id, input);
      } else {
        const input: CreateEnvironmentInput = {
          title: formTitle.value,
        };
        await environmentService.create(input);
      }
      closeForm();
      await fetchEnvironments();
    } catch (err) {
      error.value = toErrorMessage(
        err,
        editing.value
          ? t('environments.errors.updateFailed')
          : t('environments.errors.createFailed'),
      );
    } finally {
      mutating.value = false;
    }
  }

  async function removeEnvironment(environment: Environment): Promise<void> {
    if (!canMutate.value) return;
    const ok = await confirm({
      title: t('environments.delete'),
      message: t('environments.confirmDelete', { title: environment.title }),
      confirmLabel: t('environments.delete'),
      cancelLabel: t('common.cancel'),
      variant: 'danger',
    });
    if (!ok) return;

    error.value = null;
    mutating.value = true;
    try {
      await environmentService.remove(environment.id);
      await fetchEnvironments();
    } catch (err) {
      error.value = toErrorMessage(err, t('environments.errors.deleteFailed'));
    } finally {
      mutating.value = false;
    }
  }

  async function init(): Promise<void> {
    await fetchEnvironments();
  }

  return {
    environments,
    loading,
    mutating,
    error,
    titleQuery,
    appliedTitle,
    hasActiveSearch,
    formOpen,
    editing,
    formTitle,
    canMutate,
    init,
    fetchEnvironments,
    applySearch,
    clearSearch,
    openCreate,
    openEdit,
    closeForm,
    submitForm,
    removeEnvironment,
  };
}
