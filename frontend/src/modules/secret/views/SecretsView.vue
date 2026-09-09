<script setup lang="ts">
import { onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useSecretsCatalog } from '@/modules/secret/composables/useSecretsCatalog';
import { formatDateTime } from '@/shared/datetime/format';

const { t } = useI18n();
const router = useRouter();
const {
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
} = useSecretsCatalog();

function openSecret(id: string): void {
  void router.push({ name: 'secret-detail', params: { id } });
}

onMounted(() => {
  void init();
});
</script>

<template>
  <section class="px-8 py-10 md:px-12">
    <div>
      <h1 class="text-2xl font-medium tracking-tight md:text-3xl">
        {{ t('secrets.catalog.title') }}
      </h1>
      <p class="mt-3 max-w-lg text-base text-base-content/70">
        {{ t('secrets.catalog.description') }}
      </p>
    </div>

    <form class="mt-8 flex flex-wrap items-end gap-3" @submit.prevent="applySearch">
      <div class="w-full max-w-xs">
        <input
          v-model="titleQuery"
          class="input input-bordered w-full"
          type="search"
          name="secret-catalog-search"
          autocomplete="off"
          :placeholder="t('secrets.catalog.search.placeholder')"
          :aria-label="t('secrets.catalog.search.label')"
        />
      </div>
      <button class="btn btn-outline" type="submit">{{ t('secrets.catalog.search.apply') }}</button>
      <button
        v-if="hasActiveSearch"
        class="btn btn-ghost"
        type="button"
        @click="clearSearch"
      >
        {{ t('secrets.catalog.search.clear') }}
      </button>
    </form>

    <p v-if="error" class="mt-6 text-sm text-error" role="alert">{{ error }}</p>
    <p v-else-if="loading" class="mt-8 text-base-content/70">{{ t('common.loading') }}</p>

    <div v-else class="mt-8 overflow-x-auto">
      <table class="table">
        <thead>
          <tr>
            <th>{{ t('secrets.catalog.columns.name') }}</th>
            <th>{{ t('secrets.catalog.columns.projects') }}</th>
            <th>{{ t('secrets.catalog.columns.createdAt') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="items.length === 0">
            <td colspan="3" class="text-base-content/70">
              {{ hasActiveSearch ? t('secrets.catalog.emptySearch') : t('secrets.catalog.empty') }}
            </td>
          </tr>
          <tr
            v-for="secret in items"
            :key="secret.id"
            class="cursor-pointer transition-colors hover:bg-base-200/70"
            tabindex="0"
            role="link"
            :aria-label="secret.name"
            @click="openSecret(secret.id)"
            @keydown.enter.prevent="openSecret(secret.id)"
          >
            <td class="font-mono text-sm font-medium">{{ secret.name }}</td>
            <td>{{ secret.project_count }}</td>
            <td>{{ formatDateTime(secret.created_at) }}</td>
          </tr>
        </tbody>
      </table>

      <div class="mt-4 flex flex-wrap items-center justify-between gap-3">
        <p class="text-sm text-base-content/70">
          {{ t('secrets.catalog.pagination.pageOf', { page, totalPages, total }) }}
        </p>
        <div class="flex gap-2">
          <button class="btn btn-outline btn-sm" type="button" :disabled="!canPrev" @click="goPrev">
            {{ t('secrets.catalog.pagination.prev') }}
          </button>
          <button class="btn btn-outline btn-sm" type="button" :disabled="!canNext" @click="goNext">
            {{ t('secrets.catalog.pagination.next') }}
          </button>
        </div>
      </div>
    </div>
  </section>
</template>
