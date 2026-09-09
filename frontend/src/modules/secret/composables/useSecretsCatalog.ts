import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { secretService } from '@/modules/secret/services/SecretService';
import { SECRET_PAGE_SIZE, type CatalogSecret } from '@/modules/secret/types';
import { ApiError } from '@/shared/http/ApiError';

export function useSecretsCatalog() {
  const { t } = useI18n();

  const items = ref<CatalogSecret[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);

  const titleQuery = ref('');
  const appliedTitle = ref('');
  const page = ref(1);
  const total = ref(0);

  const totalPages = computed(() => Math.max(1, Math.ceil(total.value / SECRET_PAGE_SIZE)));
  const canPrev = computed(() => page.value > 1);
  const canNext = computed(() => page.value < totalPages.value);
  const hasActiveSearch = computed(() => appliedTitle.value.trim() !== '');

  function toErrorMessage(err: unknown, fallback: string): string {
    if (err instanceof ApiError) return err.message;
    return fallback;
  }

  async function fetchPage(): Promise<void> {
    error.value = null;
    loading.value = true;
    try {
      const resp = await secretService.listCatalog({
        title: appliedTitle.value.trim() || undefined,
        page: page.value,
        pageSize: SECRET_PAGE_SIZE,
      });
      items.value = resp.items;
      total.value = resp.total;
      page.value = resp.page;
    } catch (err) {
      error.value = toErrorMessage(err, t('secrets.catalog.errors.loadFailed'));
      items.value = [];
      total.value = 0;
    } finally {
      loading.value = false;
    }
  }

  async function init(): Promise<void> {
    page.value = 1;
    await fetchPage();
  }

  async function applySearch(): Promise<void> {
    appliedTitle.value = titleQuery.value.trim();
    page.value = 1;
    await fetchPage();
  }

  async function clearSearch(): Promise<void> {
    titleQuery.value = '';
    appliedTitle.value = '';
    page.value = 1;
    await fetchPage();
  }

  async function goPrev(): Promise<void> {
    if (!canPrev.value) return;
    page.value -= 1;
    await fetchPage();
  }

  async function goNext(): Promise<void> {
    if (!canNext.value) return;
    page.value += 1;
    await fetchPage();
  }

  return {
    items,
    loading,
    error,
    titleQuery,
    hasActiveSearch,
    page,
    total,
    totalPages,
    canPrev,
    canNext,
    init,
    applySearch,
    clearSearch,
    goPrev,
    goNext,
  };
}
