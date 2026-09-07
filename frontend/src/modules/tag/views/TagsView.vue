<script setup lang="ts">
import { PencilSquareIcon, PlusIcon, TrashIcon } from '@heroicons/vue/24/outline';
import { computed, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import TagFormModal from '@/modules/tag/components/TagFormModal.vue';
import { useTags } from '@/modules/tag/composables/useTags';
import { formatDate } from '@/shared/datetime/format';

const { t } = useI18n();
const {
  tags,
  loading,
  mutating,
  error,
  titleQuery,
  hasActiveSearch,
  formOpen,
  editing,
  formTitle,
  formColor,
  init,
  applySearch,
  clearSearch,
  openCreate,
  openEdit,
  closeForm,
  submitForm,
  removeTag,
} = useTags();

const formMode = computed(() => (editing.value ? 'edit' : 'create'));

onMounted(() => {
  void init();
});
</script>

<template>
  <section class="px-8 py-10 md:px-12">
    <div class="flex flex-wrap items-start justify-between gap-4">
      <div>
        <h1 class="text-2xl font-medium tracking-tight md:text-3xl">{{ t('tags.title') }}</h1>
        <p class="mt-3 max-w-lg text-base text-base-content/70">
          {{ t('tags.description') }}
        </p>
      </div>
      <button class="btn btn-primary gap-2" type="button" @click="openCreate">
        <PlusIcon class="size-4" aria-hidden="true" />
        {{ t('tags.create') }}
      </button>
    </div>

    <form class="mt-8 flex flex-wrap items-end gap-3" @submit.prevent="applySearch">
      <div class="w-full max-w-xs">
        <input
          v-model="titleQuery"
          class="input input-bordered w-full"
          type="search"
          name="tag-title-search"
          autocomplete="off"
          :placeholder="t('tags.search.placeholder')"
          :aria-label="t('tags.search.label')"
        />
      </div>
      <button class="btn btn-outline" type="submit">{{ t('tags.search.apply') }}</button>
      <button
        v-if="hasActiveSearch"
        class="btn btn-ghost"
        type="button"
        @click="clearSearch"
      >
        {{ t('tags.search.clear') }}
      </button>
    </form>

    <p v-if="error" class="mt-6 text-sm text-error" role="alert">{{ error }}</p>
    <p v-else-if="loading" class="mt-8 text-base-content/70">{{ t('common.loading') }}</p>

    <div v-else class="mt-8 overflow-x-auto">
      <table class="table">
        <thead>
          <tr>
            <th>{{ t('tags.columns.title') }}</th>
            <th>{{ t('tags.columns.color') }}</th>
            <th>{{ t('tags.columns.createdAt') }}</th>
            <th>{{ t('tags.columns.updatedAt') }}</th>
            <th>
              <span class="sr-only">{{ t('tags.columns.actions') }}</span>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="tags.length === 0">
            <td colspan="5" class="text-base-content/70">
              {{ hasActiveSearch ? t('tags.emptySearch') : t('tags.empty') }}
            </td>
          </tr>
          <tr v-for="tag in tags" :key="tag.id">
            <td>{{ tag.title }}</td>
            <td>
              <span class="inline-flex items-center gap-2">
                <span
                  class="inline-block size-4 shrink-0 rounded-full border border-base-300"
                  :style="{ backgroundColor: tag.color }"
                  aria-hidden="true"
                />
                <span class="font-mono text-sm">{{ tag.color }}</span>
              </span>
            </td>
            <td>{{ formatDate(tag.created_at) }}</td>
            <td>{{ formatDate(tag.updated_at) }}</td>
            <td class="text-right">
              <div class="flex justify-end gap-1">
                <button
                  class="btn btn-ghost btn-square btn-sm"
                  type="button"
                  :disabled="mutating"
                  :aria-label="t('tags.edit')"
                  @click="openEdit(tag)"
                >
                  <PencilSquareIcon class="size-4" aria-hidden="true" />
                </button>
                <button
                  class="btn btn-error btn-outline btn-square btn-sm"
                  type="button"
                  :disabled="mutating"
                  :aria-label="t('tags.delete')"
                  @click="removeTag(tag)"
                >
                  <TrashIcon class="size-4" aria-hidden="true" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <TagFormModal
      :open="formOpen"
      :mode="formMode"
      :title="formTitle"
      :color="formColor"
      :saving="mutating"
      @close="closeForm"
      @submit="submitForm"
      @update:title="formTitle = $event"
      @update:color="formColor = $event"
    />
  </section>
</template>
