<script setup lang="ts">
import { PencilSquareIcon, PlusIcon, TrashIcon } from '@heroicons/vue/24/outline';
import { computed, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import EnvironmentFormModal from '@/modules/environment/components/EnvironmentFormModal.vue';
import { useEnvironments } from '@/modules/environment/composables/useEnvironments';
import { formatDate } from '@/shared/datetime/format';

const { t } = useI18n();
const {
  environments,
  loading,
  mutating,
  error,
  titleQuery,
  hasActiveSearch,
  formOpen,
  editing,
  formTitle,
  canMutate,
  isProtectedEnvironment,
  init,
  applySearch,
  clearSearch,
  openCreate,
  openEdit,
  closeForm,
  submitForm,
  removeEnvironment,
} = useEnvironments();

const formMode = computed(() => (editing.value ? 'edit' : 'create'));

onMounted(() => {
  void init();
});
</script>

<template>
  <section class="px-8 py-10 md:px-12">
    <div class="flex flex-wrap items-start justify-between gap-4">
      <div>
        <h1 class="text-2xl font-medium tracking-tight md:text-3xl">
          {{ t('environments.title') }}
        </h1>
        <p class="mt-3 max-w-lg text-base text-base-content/70">
          {{ t('environments.description') }}
        </p>
      </div>
      <button
        v-if="canMutate"
        class="btn btn-primary gap-2"
        type="button"
        @click="openCreate"
      >
        <PlusIcon class="size-4" aria-hidden="true" />
        {{ t('environments.create') }}
      </button>
    </div>

    <form class="mt-8 flex flex-wrap items-end gap-3" @submit.prevent="applySearch">
      <div class="w-full max-w-xs">
        <input
          v-model="titleQuery"
          class="input input-bordered w-full"
          type="search"
          name="environment-title-search"
          autocomplete="off"
          :placeholder="t('environments.search.placeholder')"
          :aria-label="t('environments.search.label')"
        />
      </div>
      <button class="btn btn-outline" type="submit">{{ t('environments.search.apply') }}</button>
      <button
        v-if="hasActiveSearch"
        class="btn btn-ghost"
        type="button"
        @click="clearSearch"
      >
        {{ t('environments.search.clear') }}
      </button>
    </form>

    <p v-if="error" class="mt-6 text-sm text-error" role="alert">{{ error }}</p>
    <p v-else-if="loading" class="mt-8 text-base-content/70">{{ t('common.loading') }}</p>

    <div v-else class="mt-8 overflow-x-auto">
      <table class="table">
        <thead>
          <tr>
            <th>{{ t('environments.columns.title') }}</th>
            <th>{{ t('environments.columns.createdAt') }}</th>
            <th>{{ t('environments.columns.updatedAt') }}</th>
            <th v-if="canMutate">
              <span class="sr-only">{{ t('environments.columns.actions') }}</span>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="environments.length === 0">
            <td :colspan="canMutate ? 4 : 3" class="text-base-content/70">
              {{
                hasActiveSearch ? t('environments.emptySearch') : t('environments.empty')
              }}
            </td>
          </tr>
          <tr v-for="environment in environments" :key="environment.id">
            <td>{{ environment.title }}</td>
            <td>{{ formatDate(environment.created_at) }}</td>
            <td>{{ formatDate(environment.updated_at) }}</td>
            <td v-if="canMutate" class="text-right">
              <div class="flex justify-end gap-1">
                <button
                  class="btn btn-ghost btn-square btn-sm"
                  type="button"
                  :disabled="mutating || isProtectedEnvironment(environment)"
                  :aria-label="t('environments.edit')"
                  :title="
                    isProtectedEnvironment(environment)
                      ? t('environments.protected')
                      : t('environments.edit')
                  "
                  @click="openEdit(environment)"
                >
                  <PencilSquareIcon class="size-4" aria-hidden="true" />
                </button>
                <button
                  class="btn btn-error btn-outline btn-square btn-sm"
                  type="button"
                  :disabled="mutating || isProtectedEnvironment(environment)"
                  :aria-label="t('environments.delete')"
                  :title="
                    isProtectedEnvironment(environment)
                      ? t('environments.protected')
                      : t('environments.delete')
                  "
                  @click="removeEnvironment(environment)"
                >
                  <TrashIcon class="size-4" aria-hidden="true" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <EnvironmentFormModal
      v-if="canMutate"
      :open="formOpen"
      :mode="formMode"
      :title="formTitle"
      :saving="mutating"
      @close="closeForm"
      @submit="submitForm"
      @update:title="formTitle = $event"
    />
  </section>
</template>
