import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { auditService } from '@/modules/audit/services/AuditService';
import type {
  AuditListFilters,
  AuditListQuery,
  AuditLog,
  AuditSort,
} from '@/modules/audit/types';
import type { CatalogOption } from '@/modules/project/components/CatalogMultiSelect.vue';
import { projectService } from '@/modules/project/services/ProjectService';
import { PROJECT_PAGE_SIZE, type Project } from '@/modules/project/types';
import { userService } from '@/modules/user/services/UserService';
import type { UserWithToken } from '@/modules/user/types';
import { ApiError } from '@/shared/http/ApiError';

const PAGE_SIZE = 25;

const defaultFilters = (): AuditListFilters => ({
  actorUserIds: [],
  targetUserIds: [],
  projectIds: [],
  occurredFrom: '',
  occurredTo: '',
});

const defaultQuery = (): AuditListQuery => ({
  page: 1,
  sort: 'desc',
});

export function useAuditLogs() {
  const { t } = useI18n();

  const logs = ref<AuditLog[]>([]);
  const users = ref<UserWithToken[]>([]);
  const projects = ref<Project[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);

  const filters = ref<AuditListFilters>(defaultFilters());
  const query = ref<AuditListQuery>(defaultQuery());
  const total = ref(0);

  const selectedLog = ref<AuditLog | null>(null);
  const detailOpen = ref(false);

  const totalPages = computed(() => Math.max(1, Math.ceil(total.value / PAGE_SIZE)));
  const canPrev = computed(() => query.value.page > 1);
  const canNext = computed(() => query.value.page < totalPages.value);

  const userOptions = computed(() =>
    users.value.map((u) => ({ id: u.id, email: u.email })),
  );

  const projectOptions = computed<CatalogOption[]>(() =>
    projects.value.map((p) => ({ id: p.id, label: p.title })),
  );

  function toErrorMessage(err: unknown, fallback: string): string {
    if (err instanceof ApiError) {
      return err.message;
    }
    return fallback;
  }

  async function loadUsers(): Promise<void> {
    users.value = await userService.list();
  }

  async function loadProjects(): Promise<void> {
    const items: Project[] = [];
    let page = 1;
    for (;;) {
      const resp = await projectService.list({ page, page_size: PROJECT_PAGE_SIZE });
      items.push(...resp.items);
      if (items.length >= resp.total || resp.items.length === 0) break;
      page += 1;
    }
    projects.value = items;
  }

  async function fetchLogs(): Promise<void> {
    error.value = null;
    loading.value = true;
    try {
      const { actorUserIds, targetUserIds, projectIds, occurredFrom, occurredTo } = filters.value;
      const resp = await auditService.list({
        page: query.value.page,
        page_size: PAGE_SIZE,
        sort: query.value.sort,
        actor_user_ids: actorUserIds.length > 0 ? actorUserIds : undefined,
        target_user_ids: targetUserIds.length > 0 ? targetUserIds : undefined,
        project_ids: projectIds.length > 0 ? projectIds : undefined,
        occurred_from: occurredFrom || undefined,
        occurred_to: occurredTo || undefined,
      });
      logs.value = resp.items;
      total.value = resp.total;
      query.value.page = resp.page;
    } catch (err) {
      error.value = toErrorMessage(err, t('audit.errors.loadFailed'));
    } finally {
      loading.value = false;
    }
  }

  async function init(): Promise<void> {
    error.value = null;
    loading.value = true;
    try {
      await Promise.all([loadUsers(), loadProjects()]);
      await fetchLogs();
    } catch (err) {
      error.value = toErrorMessage(err, t('audit.errors.loadFailed'));
    } finally {
      loading.value = false;
    }
  }

  async function applyFilters(): Promise<void> {
    query.value.page = 1;
    await fetchLogs();
  }

  async function refresh(): Promise<void> {
    await fetchLogs();
  }

  async function toggleSort(): Promise<void> {
    const next: AuditSort = query.value.sort === 'desc' ? 'asc' : 'desc';
    query.value.sort = next;
    await applyFilters();
  }

  async function goPrev(): Promise<void> {
    if (!canPrev.value) return;
    query.value.page -= 1;
    await fetchLogs();
  }

  async function goNext(): Promise<void> {
    if (!canNext.value) return;
    query.value.page += 1;
    await fetchLogs();
  }

  function openDetail(log: AuditLog): void {
    selectedLog.value = log;
    detailOpen.value = true;
  }

  function closeDetail(): void {
    detailOpen.value = false;
    selectedLog.value = null;
  }

  return {
    logs,
    users,
    userOptions,
    projectOptions,
    loading,
    error,
    filters,
    query,
    total,
    totalPages,
    selectedLog,
    detailOpen,
    canPrev,
    canNext,
    init,
    applyFilters,
    refresh,
    toggleSort,
    goPrev,
    goNext,
    openDetail,
    closeDetail,
  };
}
