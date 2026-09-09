<script setup lang="ts">
import { ArrowLeftIcon, MagnifyingGlassIcon, PencilSquareIcon, PlusIcon, TrashIcon } from '@heroicons/vue/24/outline';
import { computed, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { RouterLink } from 'vue-router';
import SecretFormModal from '@/modules/secret/components/SecretFormModal.vue';
import ProjectDownloadSnippet from '@/modules/secret/components/ProjectDownloadSnippet.vue';
import { useProjectSecrets } from '@/modules/secret/composables/useProjectSecrets';
import { formatDateTime } from '@/shared/datetime/format';

const { t } = useI18n();
const {
  project,
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
} = useProjectSecrets();

const formMode = computed(() => (editingSecret.value ? 'edit' : 'create'));
const formProjectCount = computed(() => editingSecret.value?.project_count ?? 1);

onMounted(() => {
  void load();
});
</script>

<template>
  <section class="px-8 py-10 md:px-12">
    <RouterLink
      :to="{ name: 'home' }"
      class="inline-flex items-center gap-2 text-sm text-base-content/70 hover:text-base-content"
    >
      <ArrowLeftIcon class="size-4" aria-hidden="true" />
      {{ t('secrets.back') }}
    </RouterLink>

    <div class="mt-4 flex flex-wrap items-start justify-between gap-4">
      <div>
        <h1 class="text-2xl font-medium tracking-tight md:text-3xl">
          {{ project?.title ?? t('common.loading') }}
        </h1>
        <p v-if="project" class="mt-1 font-mono text-sm text-base-content/60">
          {{ project.slug }}
        </p>
      </div>
      <button
        class="btn btn-primary gap-2"
        type="button"
        :disabled="!project || mutating"
        @click="openCreate()"
      >
        <PlusIcon class="size-4" aria-hidden="true" />
        {{ t('secrets.create') }}
      </button>
    </div>

    <p v-if="error" class="mt-6 text-sm text-error" role="alert">{{ error }}</p>
    <p v-else-if="loading" class="mt-8 text-base-content/70">{{ t('common.loading') }}</p>
    <p v-else-if="!project" class="mt-8 text-base-content/70">{{ t('secrets.emptyProject') }}</p>

    <template v-else>
      <div class="mt-8">
        <ProjectDownloadSnippet
          :project-slug="project.slug"
          :environments="project.environments"
        />
      </div>

      <form
        class="mt-8 flex flex-wrap items-end gap-3"
        @submit.prevent="runCatalogSearch"
      >
        <div class="w-full max-w-xs">
          <input
            v-model="environmentQuery"
            class="input input-bordered w-full"
            type="search"
            name="secret-environment-search"
            autocomplete="off"
            :placeholder="t('secrets.search.environmentPlaceholder')"
            :aria-label="t('secrets.search.environmentLabel')"
          />
        </div>
        <button
          class="btn btn-outline gap-2"
          type="submit"
          :disabled="mutating || catalogSearching || !environmentQuery.trim()"
        >
          <MagnifyingGlassIcon class="size-4" aria-hidden="true" />
          {{ t('secrets.search.submit') }}
        </button>
      </form>

      <div
        v-if="catalogHits.length > 0"
        class="mt-4 rounded-box border border-base-300 bg-base-200/40 p-3"
      >
        <p class="text-sm font-medium text-base-content/80">
          {{ t('secrets.search.catalogTitle') }}
        </p>
        <ul class="mt-2 space-y-1">
          <li
            v-for="hit in catalogHits"
            :key="hit.id"
            class="flex flex-wrap items-center justify-between gap-2 rounded-lg px-2 py-1.5 hover:bg-base-200"
          >
            <span class="font-mono text-sm">{{ hit.name }}</span>
            <button
              class="btn btn-ghost btn-xs"
              type="button"
              :disabled="mutating"
              @click="attachCatalogSecret(hit)"
            >
              {{ t('secrets.search.attach') }}
            </button>
          </li>
        </ul>
      </div>

      <div class="mt-4 overflow-x-auto">
        <table class="table">
          <thead>
            <tr>
              <th>{{ t('secrets.columns.name') }}</th>
              <th>{{ t('secrets.columns.environments') }}</th>
              <th>{{ t('secrets.columns.createdAt') }}</th>
              <th>
                <span class="sr-only">{{ t('secrets.columns.actions') }}</span>
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="filteredSecrets.length === 0">
              <td colspan="4" class="text-base-content/70">
                {{
                  hasEnvironmentFilter ? t('secrets.emptySearch') : t('secrets.empty')
                }}
              </td>
            </tr>
            <tr v-for="secret in filteredSecrets" :key="secret.id">
              <td class="font-mono text-sm">{{ secret.name }}</td>
              <td>
                <span class="flex flex-wrap gap-1">
                  <span
                    v-for="title in envTitlesForSecret(secret)"
                    :key="title"
                    class="badge badge-ghost badge-sm"
                  >
                    {{ title }}
                  </span>
                  <span v-if="envTitlesForSecret(secret).length === 0">{{ t('common.emDash') }}</span>
                </span>
              </td>
              <td>{{ formatDateTime(secret.created_at) }}</td>
              <td class="text-right">
                <div class="flex justify-end gap-1">
                  <button
                    class="btn btn-ghost btn-square btn-sm"
                    type="button"
                    :disabled="mutating"
                    :aria-label="t('secrets.edit')"
                    @click="openEdit(secret)"
                  >
                    <PencilSquareIcon class="size-4" aria-hidden="true" />
                  </button>
                  <button
                    class="btn btn-error btn-outline btn-square btn-sm"
                    type="button"
                    :disabled="mutating"
                    :aria-label="t('secrets.delete')"
                    @click="removeSecret(secret)"
                  >
                    <TrashIcon class="size-4" aria-hidden="true" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>

    <SecretFormModal
      :open="formOpen"
      :mode="formMode"
      :name="formName"
      :saving="mutating"
      :project-count="formProjectCount"
      :environment-options="environmentOptions"
      :selected-environment-ids="selectedEnvironmentIds"
      :values-by-env="valuesByEnv"
      :prefixes-by-env="prefixesByEnv"
      @close="closeForm"
      @submit="submitForm"
      @update:name="formName = $event"
      @update:selected-environment-ids="setSelectedEnvironmentIds"
      @update:value="setEnvValue($event.environmentId, $event.value)"
      @update:prefix="setEnvPrefix($event.environmentId, $event.prefix)"
    />
  </section>
</template>
