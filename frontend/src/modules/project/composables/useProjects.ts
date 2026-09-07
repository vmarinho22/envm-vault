import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useAuthStore } from '@/modules/auth/stores/auth.store';
import { environmentService } from '@/modules/environment/services/EnvironmentService';
import type { CatalogOption } from '@/modules/project/components/CatalogMultiSelect.vue';
import { projectService } from '@/modules/project/services/ProjectService';
import {
  PROJECT_PAGE_SIZE,
  type CreateProjectInput,
  type Project,
  type UpdateProjectInput,
} from '@/modules/project/types';
import { tagService } from '@/modules/tag/services/TagService';
import { randomTagColor } from '@/modules/tag/types';
import { teamService } from '@/modules/team/services/TeamService';
import { userService } from '@/modules/user/services/UserService';
import type { UserWithToken } from '@/modules/user/types';
import { useConfirmDialog } from '@/shared/confirm/useConfirmDialog';
import { ApiError } from '@/shared/http/ApiError';

export type CatalogQuickKind = 'team' | 'environment' | 'tag';

export function useProjects() {
  const { t } = useI18n();
  const auth = useAuthStore();
  const { confirm } = useConfirmDialog();

  const projects = ref<Project[]>([]);
  const teamOptions = ref<CatalogOption[]>([]);
  const environmentOptions = ref<CatalogOption[]>([]);
  const tagOptions = ref<CatalogOption[]>([]);

  const loading = ref(false);
  const mutating = ref(false);
  const error = ref<string | null>(null);

  const titleQuery = ref('');
  const filterTeamIds = ref<string[]>([]);
  const filterTagIds = ref<string[]>([]);
  const appliedTitle = ref('');
  const appliedTeamIds = ref<string[]>([]);
  const appliedTagIds = ref<string[]>([]);
  const page = ref(1);
  const total = ref(0);

  const formOpen = ref(false);
  const editing = ref<Project | null>(null);
  const formTitle = ref('');
  const formTeamIds = ref<string[]>([]);
  const formEnvironmentIds = ref<string[]>([]);
  const formTagIds = ref<string[]>([]);

  const catalogQuick = ref<CatalogQuickKind | null>(null);
  const quickUsers = ref<UserWithToken[]>([]);
  const quickSaving = ref(false);
  const quickError = ref<string | null>(null);

  const teamQuickName = ref('');
  const teamQuickAccountableId = ref('');
  const environmentQuickTitle = ref('');
  const tagQuickTitle = ref('');
  const tagQuickColor = ref(randomTagColor());

  const canDelete = computed(() => auth.isAdmin);
  const canAddTeam = computed(() => auth.isAdmin);
  const canAddEnvironment = computed(() => auth.isAdmin);
  const canAddTag = computed(() => true);
  const totalPages = computed(() => Math.max(1, Math.ceil(total.value / PROJECT_PAGE_SIZE)));
  const canPrev = computed(() => page.value > 1);
  const canNext = computed(() => page.value < totalPages.value);
  const hasActiveSearch = computed(
    () =>
      appliedTitle.value.trim() !== '' ||
      appliedTeamIds.value.length > 0 ||
      appliedTagIds.value.length > 0,
  );

  function toErrorMessage(err: unknown, fallback: string): string {
    if (err instanceof ApiError) return err.message;
    return fallback;
  }

  async function fetchCatalogs(): Promise<void> {
    const [teams, environments, tags] = await Promise.all([
      teamService.list(),
      environmentService.list(),
      tagService.list(),
    ]);
    teamOptions.value = teams.map((x) => ({ id: x.id, label: x.name }));
    environmentOptions.value = environments.map((x) => ({ id: x.id, label: x.title }));
    tagOptions.value = tags.map((x) => ({ id: x.id, label: x.title }));
  }

  async function refreshTeams(): Promise<void> {
    const teams = await teamService.list();
    teamOptions.value = teams.map((x) => ({ id: x.id, label: x.name }));
  }

  async function refreshEnvironments(): Promise<void> {
    const environments = await environmentService.list();
    environmentOptions.value = environments.map((x) => ({ id: x.id, label: x.title }));
  }

  async function refreshTags(): Promise<void> {
    const tags = await tagService.list();
    tagOptions.value = tags.map((x) => ({ id: x.id, label: x.title }));
  }

  async function ensureQuickUsers(): Promise<void> {
    if (!auth.isAdmin) return;
    if (quickUsers.value.length > 0) return;
    quickUsers.value = await userService.list();
  }

  async function fetchProjects(): Promise<void> {
    error.value = null;
    loading.value = true;
    try {
      const resp = await projectService.list({
        title: appliedTitle.value.trim() || undefined,
        team_ids: appliedTeamIds.value.length ? appliedTeamIds.value : undefined,
        tag_ids: appliedTagIds.value.length ? appliedTagIds.value : undefined,
        page: page.value,
        page_size: PROJECT_PAGE_SIZE,
      });
      projects.value = resp.items;
      total.value = resp.total;
      page.value = resp.page;
    } catch (err) {
      error.value = toErrorMessage(err, t('projects.errors.loadFailed'));
    } finally {
      loading.value = false;
    }
  }

  async function applySearch(): Promise<void> {
    appliedTitle.value = titleQuery.value.trim();
    appliedTeamIds.value = [...filterTeamIds.value];
    appliedTagIds.value = [...filterTagIds.value];
    page.value = 1;
    await fetchProjects();
  }

  async function clearSearch(): Promise<void> {
    titleQuery.value = '';
    filterTeamIds.value = [];
    filterTagIds.value = [];
    appliedTitle.value = '';
    appliedTeamIds.value = [];
    appliedTagIds.value = [];
    page.value = 1;
    await fetchProjects();
  }

  async function goPrev(): Promise<void> {
    if (!canPrev.value) return;
    page.value -= 1;
    await fetchProjects();
  }

  async function goNext(): Promise<void> {
    if (!canNext.value) return;
    page.value += 1;
    await fetchProjects();
  }

  function closeQuickForms(): void {
    catalogQuick.value = null;
    quickError.value = null;
  }

  function openCreate(): void {
    editing.value = null;
    formTitle.value = '';
    formTeamIds.value = [];
    formEnvironmentIds.value = [];
    formTagIds.value = [];
    closeQuickForms();
    formOpen.value = true;
  }

  function openEdit(project: Project): void {
    editing.value = project;
    formTitle.value = project.title;
    formTeamIds.value = project.teams.map((x) => x.id);
    formEnvironmentIds.value = project.environments.map((x) => x.id);
    formTagIds.value = project.tags.map((x) => x.id);
    closeQuickForms();
    formOpen.value = true;
  }

  function closeForm(): void {
    formOpen.value = false;
    editing.value = null;
    closeQuickForms();
  }

  async function submitForm(): Promise<void> {
    error.value = null;
    mutating.value = true;
    try {
      if (editing.value) {
        const input: UpdateProjectInput = {
          title: formTitle.value,
          team_ids: formTeamIds.value,
          environment_ids: formEnvironmentIds.value,
          tag_ids: formTagIds.value,
        };
        await projectService.update(editing.value.id, input);
      } else {
        const input: CreateProjectInput = {
          title: formTitle.value,
          team_ids: formTeamIds.value,
          environment_ids: formEnvironmentIds.value,
          tag_ids: formTagIds.value,
        };
        await projectService.create(input);
      }
      closeForm();
      await fetchProjects();
    } catch (err) {
      error.value = toErrorMessage(
        err,
        editing.value ? t('projects.errors.updateFailed') : t('projects.errors.createFailed'),
      );
    } finally {
      mutating.value = false;
    }
  }

  async function removeProject(project: Project): Promise<void> {
    if (!canDelete.value) return;
    const ok = await confirm({
      title: t('projects.delete'),
      message: t('projects.confirmDelete', { title: project.title }),
      confirmLabel: t('projects.delete'),
      cancelLabel: t('common.cancel'),
      variant: 'danger',
    });
    if (!ok) return;

    error.value = null;
    mutating.value = true;
    try {
      await projectService.remove(project.id);
      await fetchProjects();
    } catch (err) {
      error.value = toErrorMessage(err, t('projects.errors.deleteFailed'));
    } finally {
      mutating.value = false;
    }
  }

  async function openQuickTeam(): Promise<void> {
    if (!canAddTeam.value) return;
    quickError.value = null;
    try {
      await ensureQuickUsers();
      teamQuickName.value = '';
      teamQuickAccountableId.value = quickUsers.value[0]?.id ?? '';
      catalogQuick.value = 'team';
    } catch (err) {
      quickError.value = toErrorMessage(err, t('projects.errors.loadUsersFailed'));
      catalogQuick.value = 'team';
    }
  }

  function openQuickEnvironment(): void {
    if (!canAddEnvironment.value) return;
    quickError.value = null;
    environmentQuickTitle.value = '';
    catalogQuick.value = 'environment';
  }

  function openQuickTag(): void {
    quickError.value = null;
    tagQuickTitle.value = '';
    tagQuickColor.value = randomTagColor();
    catalogQuick.value = 'tag';
  }

  function closeCatalogQuick(): void {
    closeQuickForms();
  }

  async function submitQuickTeam(): Promise<void> {
    if (!canAddTeam.value) return;
    quickError.value = null;
    quickSaving.value = true;
    try {
      const created = await teamService.create({
        name: teamQuickName.value,
        accountable_id: teamQuickAccountableId.value,
      });
      await refreshTeams();
      if (!formTeamIds.value.includes(created.id)) {
        formTeamIds.value = [...formTeamIds.value, created.id];
      }
      closeCatalogQuick();
    } catch (err) {
      quickError.value = toErrorMessage(err, t('projects.errors.createTeamFailed'));
    } finally {
      quickSaving.value = false;
    }
  }

  async function submitQuickEnvironment(): Promise<void> {
    if (!canAddEnvironment.value) return;
    quickError.value = null;
    quickSaving.value = true;
    try {
      const created = await environmentService.create({ title: environmentQuickTitle.value });
      await refreshEnvironments();
      if (!formEnvironmentIds.value.includes(created.id)) {
        formEnvironmentIds.value = [...formEnvironmentIds.value, created.id];
      }
      closeCatalogQuick();
    } catch (err) {
      quickError.value = toErrorMessage(err, t('projects.errors.createEnvironmentFailed'));
    } finally {
      quickSaving.value = false;
    }
  }

  async function submitQuickTag(): Promise<void> {
    quickError.value = null;
    quickSaving.value = true;
    try {
      const created = await tagService.create({
        title: tagQuickTitle.value,
        color: tagQuickColor.value,
      });
      await refreshTags();
      if (!formTagIds.value.includes(created.id)) {
        formTagIds.value = [...formTagIds.value, created.id];
      }
      closeCatalogQuick();
    } catch (err) {
      quickError.value = toErrorMessage(err, t('projects.errors.createTagFailed'));
    } finally {
      quickSaving.value = false;
    }
  }

  async function init(): Promise<void> {
    error.value = null;
    loading.value = true;
    try {
      await fetchCatalogs();
      await fetchProjects();
    } catch (err) {
      error.value = toErrorMessage(err, t('projects.errors.loadFailed'));
    } finally {
      loading.value = false;
    }
  }

  return {
    projects,
    teamOptions,
    environmentOptions,
    tagOptions,
    loading,
    mutating,
    error,
    titleQuery,
    filterTeamIds,
    filterTagIds,
    hasActiveSearch,
    page,
    total,
    totalPages,
    canPrev,
    canNext,
    formOpen,
    editing,
    formTitle,
    formTeamIds,
    formEnvironmentIds,
    formTagIds,
    canDelete,
    canAddTeam,
    canAddEnvironment,
    canAddTag,
    catalogQuick,
    quickUsers,
    quickSaving,
    quickError,
    teamQuickName,
    teamQuickAccountableId,
    environmentQuickTitle,
    tagQuickTitle,
    tagQuickColor,
    init,
    applySearch,
    clearSearch,
    goPrev,
    goNext,
    openCreate,
    openEdit,
    closeForm,
    submitForm,
    removeProject,
    openQuickTeam,
    openQuickEnvironment,
    openQuickTag,
    closeCatalogQuick,
    submitQuickTeam,
    submitQuickEnvironment,
    submitQuickTag,
  };
}
