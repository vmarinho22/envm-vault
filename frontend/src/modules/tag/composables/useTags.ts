import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { tagService } from '@/modules/tag/services/TagService';
import { randomTagColor, type CreateTagInput, type Tag, type UpdateTagInput } from '@/modules/tag/types';
import { useConfirmDialog } from '@/shared/confirm/useConfirmDialog';
import { ApiError } from '@/shared/http/ApiError';

export function useTags() {
  const { t } = useI18n();
  const { confirm } = useConfirmDialog();

  const tags = ref<Tag[]>([]);
  const loading = ref(false);
  const mutating = ref(false);
  const error = ref<string | null>(null);
  const titleQuery = ref('');
  const appliedTitle = ref('');

  const formOpen = ref(false);
  const editing = ref<Tag | null>(null);
  const formTitle = ref('');
  const formColor = ref(randomTagColor());

  const hasActiveSearch = computed(() => appliedTitle.value.trim() !== '');

  function toErrorMessage(err: unknown, fallback: string): string {
    if (err instanceof ApiError) {
      return err.message;
    }
    return fallback;
  }

  async function fetchTags(): Promise<void> {
    error.value = null;
    loading.value = true;
    try {
      const title = appliedTitle.value.trim();
      tags.value = await tagService.list(title ? { title } : undefined);
    } catch (err) {
      error.value = toErrorMessage(err, t('tags.errors.loadFailed'));
    } finally {
      loading.value = false;
    }
  }

  async function applySearch(): Promise<void> {
    appliedTitle.value = titleQuery.value.trim();
    await fetchTags();
  }

  async function clearSearch(): Promise<void> {
    titleQuery.value = '';
    appliedTitle.value = '';
    await fetchTags();
  }

  function openCreate(): void {
    editing.value = null;
    formTitle.value = '';
    formColor.value = randomTagColor();
    formOpen.value = true;
  }

  function openEdit(tag: Tag): void {
    editing.value = tag;
    formTitle.value = tag.title;
    formColor.value = tag.color;
    formOpen.value = true;
  }

  function closeForm(): void {
    formOpen.value = false;
    editing.value = null;
  }

  async function submitForm(): Promise<void> {
    error.value = null;
    mutating.value = true;
    try {
      if (editing.value) {
        const input: UpdateTagInput = {
          title: formTitle.value,
          color: formColor.value,
        };
        await tagService.update(editing.value.id, input);
      } else {
        const input: CreateTagInput = {
          title: formTitle.value,
          color: formColor.value,
        };
        await tagService.create(input);
      }
      closeForm();
      await fetchTags();
    } catch (err) {
      error.value = toErrorMessage(
        err,
        editing.value ? t('tags.errors.updateFailed') : t('tags.errors.createFailed'),
      );
    } finally {
      mutating.value = false;
    }
  }

  async function removeTag(tag: Tag): Promise<void> {
    const ok = await confirm({
      title: t('tags.delete'),
      message: t('tags.confirmDelete', { title: tag.title }),
      confirmLabel: t('tags.delete'),
      cancelLabel: t('common.cancel'),
      variant: 'danger',
    });
    if (!ok) return;

    error.value = null;
    mutating.value = true;
    try {
      await tagService.remove(tag.id);
      await fetchTags();
    } catch (err) {
      error.value = toErrorMessage(err, t('tags.errors.deleteFailed'));
    } finally {
      mutating.value = false;
    }
  }

  async function init(): Promise<void> {
    await fetchTags();
  }

  return {
    tags,
    loading,
    mutating,
    error,
    titleQuery,
    appliedTitle,
    hasActiveSearch,
    formOpen,
    editing,
    formTitle,
    formColor,
    init,
    fetchTags,
    applySearch,
    clearSearch,
    openCreate,
    openEdit,
    closeForm,
    submitForm,
    removeTag,
  };
}
