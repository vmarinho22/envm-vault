import { computed, onUnmounted, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute } from 'vue-router';
import type { CatalogOption } from '@/modules/project/components/CatalogMultiSelect.vue';
import { projectService } from '@/modules/project/services/ProjectService';
import type { Project } from '@/modules/project/types';
import { secretService } from '@/modules/secret/services/SecretService';
import type { CatalogSecret, Secret } from '@/modules/secret/types';
import { useConfirmDialog } from '@/shared/confirm/useConfirmDialog';
import { ApiError } from '@/shared/http/ApiError';

const SEARCH_DEBOUNCE_MS = 300;

export function useProjectSecrets() {
  const { t } = useI18n();
  const route = useRoute();
  const { confirm } = useConfirmDialog();

  const project = ref<Project | null>(null);
  const secrets = ref<Secret[]>([]);
  const environmentQuery = ref('');
  const catalogHits = ref<CatalogSecret[]>([]);
  const catalogSearching = ref(false);
  const loading = ref(false);
  const mutating = ref(false);
  const error = ref<string | null>(null);

  const formOpen = ref(false);
  const formName = ref('');
  const editingSecret = ref<Secret | null>(null);
  const selectedEnvironmentIds = ref<string[]>([]);
  const valuesByEnv = ref<Record<string, string>>({});
  const prefixesByEnv = ref<Record<string, string>>({});

  let debounceTimer: ReturnType<typeof setTimeout> | null = null;

  const slug = computed(() => String(route.params.slug ?? ''));

  const environmentOptions = computed<CatalogOption[]>(() =>
    (project.value?.environments ?? []).map((e) => ({ id: e.id, label: e.title })),
  );

  function toErrorMessage(err: unknown, fallback: string): string {
    if (err instanceof ApiError) {
      return err.message;
    }
    return fallback;
  }

  function envTitlesForSecret(secret: Secret): string[] {
    if (!project.value) return [];
    const titles: string[] = [];
    for (const env of project.value.environments) {
      if (secret.values[env.id]?.has_value) {
        titles.push(env.title);
      }
    }
    return titles;
  }

  const filteredSecrets = computed(() => {
    const term = environmentQuery.value.trim().toLowerCase();
    if (!term) return secrets.value;
    return secrets.value.filter((secret) => {
      if (secret.name.toLowerCase().includes(term)) return true;
      return envTitlesForSecret(secret).some((title) => title.toLowerCase().includes(term));
    });
  });

  const hasEnvironmentFilter = computed(() => environmentQuery.value.trim() !== '');

  async function refreshSecrets(): Promise<void> {
    if (!project.value) return;
    const list = await secretService.list(project.value.id);
    secrets.value = list.items;
  }

  async function load(): Promise<void> {
    error.value = null;
    loading.value = true;
    try {
      const p = await projectService.getBySlug(slug.value);
      project.value = p;
      await refreshSecrets();
    } catch (err) {
      error.value = toErrorMessage(err, t('secrets.errors.loadFailed'));
      project.value = null;
      secrets.value = [];
    } finally {
      loading.value = false;
    }
  }

  async function fetchCatalogPreview(term: string): Promise<void> {
    if (!project.value || !term) {
      catalogHits.value = [];
      return;
    }
    catalogSearching.value = true;
    try {
      const resp = await secretService.listCatalog({
        title: term,
        excludeProjectId: project.value.id,
        page: 1,
        pageSize: 20,
      });
      catalogHits.value = resp.items;
    } catch {
      catalogHits.value = [];
    } finally {
      catalogSearching.value = false;
    }
  }

  watch(environmentQuery, (value) => {
    if (debounceTimer) clearTimeout(debounceTimer);
    const term = value.trim();
    if (!term) {
      catalogHits.value = [];
      return;
    }
    debounceTimer = setTimeout(() => {
      void fetchCatalogPreview(term);
    }, SEARCH_DEBOUNCE_MS);
  });

  onUnmounted(() => {
    if (debounceTimer) clearTimeout(debounceTimer);
  });

  async function runCatalogSearch(): Promise<void> {
    if (!project.value) return;
    const term = environmentQuery.value.trim();
    if (!term) {
      catalogHits.value = [];
      return;
    }
    error.value = null;
    catalogSearching.value = true;
    try {
      const resp = await secretService.listCatalog({
        title: term,
        excludeProjectId: project.value.id,
        page: 1,
        pageSize: 20,
      });
      catalogHits.value = resp.items;
      if (resp.items.length === 0) {
        openCreate(term);
      }
    } catch (err) {
      error.value = toErrorMessage(err, t('secrets.errors.searchFailed'));
      catalogHits.value = [];
    } finally {
      catalogSearching.value = false;
    }
  }

  async function attachCatalogSecret(hit: CatalogSecret): Promise<void> {
    if (!project.value) return;
    error.value = null;
    mutating.value = true;
    try {
      await secretService.attach(project.value.id, hit.id);
      catalogHits.value = catalogHits.value.filter((h) => h.id !== hit.id);
      await refreshSecrets();
    } catch (err) {
      error.value = toErrorMessage(err, t('secrets.errors.attachFailed'));
    } finally {
      mutating.value = false;
    }
  }

  function openCreate(prefillName = ''): void {
    editingSecret.value = null;
    formName.value = prefillName;
    selectedEnvironmentIds.value = [];
    valuesByEnv.value = {};
    prefixesByEnv.value = {};
    formOpen.value = true;
  }

  async function openEdit(secret: Secret): Promise<void> {
    if (!project.value) return;
    error.value = null;
    mutating.value = true;
    editingSecret.value = secret;
    formName.value = secret.name;
    // Only environments linked to this project — shared secrets may have
    // values for other projects' environments.
    const projectEnvIds = new Set(project.value.environments.map((e) => e.id));
    const envIds = Object.keys(secret.values).filter(
      (id) =>
        projectEnvIds.has(id) &&
        (secret.values[id]?.has_value || Boolean(secret.values[id]?.prefix)),
    );
    // Prefer envs that have values; still include prefix-only for editing.
    const valueEnvIds = envIds.filter((id) => secret.values[id]?.has_value);
    selectedEnvironmentIds.value = valueEnvIds.length > 0 ? valueEnvIds : envIds;
    valuesByEnv.value = Object.fromEntries(selectedEnvironmentIds.value.map((id) => [id, '']));
    prefixesByEnv.value = Object.fromEntries(
      selectedEnvironmentIds.value.map((id) => [id, secret.values[id]?.prefix ?? '']),
    );
    formOpen.value = true;

    try {
      const revealed = await Promise.all(
        selectedEnvironmentIds.value.map(async (environmentId) => {
          if (!secret.values[environmentId]?.has_value) {
            return [environmentId, ''] as const;
          }
          const resp = await secretService.revealValue(project.value!.id, secret.id, environmentId);
          return [environmentId, resp.value] as const;
        }),
      );
      // Ignore if user closed the form while loading.
      if (editingSecret.value?.id !== secret.id) return;
      valuesByEnv.value = Object.fromEntries(revealed);
    } catch (err) {
      closeForm();
      error.value = toErrorMessage(err, t('secrets.errors.revealFailed'));
    } finally {
      mutating.value = false;
    }
  }

  function closeForm(): void {
    formOpen.value = false;
    editingSecret.value = null;
    formName.value = '';
    selectedEnvironmentIds.value = [];
    valuesByEnv.value = {};
    prefixesByEnv.value = {};
  }

  function setSelectedEnvironmentIds(ids: string[]): void {
    selectedEnvironmentIds.value = ids;
    const nextValues: Record<string, string> = {};
    const nextPrefixes: Record<string, string> = {};
    for (const id of ids) {
      nextValues[id] = valuesByEnv.value[id] ?? '';
      nextPrefixes[id] =
        prefixesByEnv.value[id] ?? editingSecret.value?.values[id]?.prefix ?? '';
    }
    valuesByEnv.value = nextValues;
    prefixesByEnv.value = nextPrefixes;
  }

  function setEnvValue(environmentId: string, value: string): void {
    valuesByEnv.value = { ...valuesByEnv.value, [environmentId]: value };
  }

  function setEnvPrefix(environmentId: string, prefix: string): void {
    prefixesByEnv.value = { ...prefixesByEnv.value, [environmentId]: prefix };
  }

  async function persistPrefixes(secretId: string, environmentIds: string[]): Promise<void> {
    if (!project.value) return;
    for (const id of environmentIds) {
      const prefix = (prefixesByEnv.value[id] ?? '').trim();
      await secretService.setPrefix(project.value.id, secretId, id, prefix);
    }
  }

  async function submitForm(): Promise<void> {
    if (!project.value) return;
    error.value = null;
    mutating.value = true;
    try {
      const name = formName.value.trim();
      if (editingSecret.value) {
        const secret = editingSecret.value;
        if (name !== secret.name) {
          await secretService.updateName(project.value.id, secret.id, name);
        }

        const previousIds = new Set(
          Object.keys(secret.values).filter((id) => secret.values[id]?.has_value),
        );
        const previousPrefixIds = new Set(
          Object.keys(secret.values).filter((id) => Boolean(secret.values[id]?.prefix)),
        );
        const nextIds = new Set(selectedEnvironmentIds.value);

        for (const id of previousIds) {
          if (!nextIds.has(id)) {
            await secretService.clearValue(project.value.id, secret.id, id);
          }
        }
        for (const id of previousPrefixIds) {
          if (!nextIds.has(id)) {
            await secretService.clearPrefix(project.value.id, secret.id, id);
          }
        }

        for (const id of nextIds) {
          const plaintext = valuesByEnv.value[id] ?? '';
          if (plaintext.length === 0) {
            throw new ApiError(400, t('secrets.errors.valueRequired'));
          }
          await secretService.setValue(project.value.id, secret.id, id, plaintext);
        }
        await persistPrefixes(secret.id, selectedEnvironmentIds.value);
      } else {
        const values: Record<string, string> = {};
        for (const id of selectedEnvironmentIds.value) {
          values[id] = valuesByEnv.value[id] ?? '';
        }
        const created = await secretService.create(project.value.id, { name, values });
        await persistPrefixes(created.id, selectedEnvironmentIds.value);
      }
      closeForm();
      await refreshSecrets();
    } catch (err) {
      error.value = toErrorMessage(
        err,
        editingSecret.value ? t('secrets.errors.updateFailed') : t('secrets.errors.createFailed'),
      );
    } finally {
      mutating.value = false;
    }
  }

  async function removeSecret(secret: Secret): Promise<void> {
    if (!project.value) return;
    const ok = await confirm({
      title: t('secrets.delete'),
      message: t('secrets.confirmDetach', { name: secret.name }),
      confirmLabel: t('common.confirm'),
      cancelLabel: t('common.cancel'),
      variant: 'danger',
    });
    if (!ok) return;

    error.value = null;
    mutating.value = true;
    try {
      await secretService.remove(project.value.id, secret.id);
      await refreshSecrets();
    } catch (err) {
      error.value = toErrorMessage(err, t('secrets.errors.deleteFailed'));
    } finally {
      mutating.value = false;
    }
  }

  return {
    project,
    secrets,
    filteredSecrets,
    environmentQuery,
    catalogHits,
    catalogSearching,
    hasEnvironmentFilter,
    loading,
    mutating,
    error,
    formOpen,
    formName,
    editingSecret,
    selectedEnvironmentIds,
    valuesByEnv,
    prefixesByEnv,
    environmentOptions,
    slug,
    load,
    runCatalogSearch,
    attachCatalogSecret,
    openCreate,
    openEdit,
    closeForm,
    setSelectedEnvironmentIds,
    setEnvValue,
    setEnvPrefix,
    submitForm,
    removeSecret,
    envTitlesForSecret,
  };
}
