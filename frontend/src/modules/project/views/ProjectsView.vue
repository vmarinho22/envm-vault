<script setup lang="ts">
import { PencilSquareIcon, PlusIcon, TrashIcon } from '@heroicons/vue/24/outline';
import { computed, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { RouterLink } from 'vue-router';
import CatalogMultiSelect from '@/modules/project/components/CatalogMultiSelect.vue';
import ProjectFormModal from '@/modules/project/components/ProjectFormModal.vue';
import { useProjects } from '@/modules/project/composables/useProjects';
import { formatDate } from '@/shared/datetime/format';

const { t } = useI18n();
const {
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
} = useProjects();

const formMode = computed(() => (editing.value ? 'edit' : 'create'));

onMounted(() => {
  void init();
});
</script>

<template>
  <section class="px-8 py-10 md:px-12">
    <div class="flex flex-wrap items-start justify-between gap-4">
      <div>
        <h1 class="text-2xl font-medium tracking-tight md:text-3xl">{{ t('projects.title') }}</h1>
        <p class="mt-3 max-w-lg text-base text-base-content/70">
          {{ t('projects.description') }}
        </p>
      </div>
      <button class="btn btn-primary gap-2" type="button" @click="openCreate">
        <PlusIcon class="size-4" aria-hidden="true" />
        {{ t('projects.create') }}
      </button>
    </div>

    <form class="mt-8 flex flex-wrap items-end gap-3" @submit.prevent="applySearch">
      <div class="w-full max-w-xs">
        <input
          v-model="titleQuery"
          class="input input-bordered w-full"
          type="search"
          name="project-title-search"
          autocomplete="off"
          :placeholder="t('projects.search.placeholder')"
          :aria-label="t('projects.search.label')"
        />
      </div>
      <div class="w-full max-w-xs">
        <CatalogMultiSelect
          v-model="filterTeamIds"
          :options="teamOptions"
          :label="t('projects.search.teams')"
          :placeholder="t('projects.search.allTeams')"
        />
      </div>
      <div class="w-full max-w-xs">
        <CatalogMultiSelect
          v-model="filterTagIds"
          :options="tagOptions"
          :label="t('projects.search.tags')"
          :placeholder="t('projects.search.allTags')"
        />
      </div>
      <button class="btn btn-outline" type="submit">{{ t('projects.search.apply') }}</button>
      <button
        v-if="hasActiveSearch"
        class="btn btn-ghost"
        type="button"
        @click="clearSearch"
      >
        {{ t('projects.search.clear') }}
      </button>
    </form>

    <p v-if="error" class="mt-6 text-sm text-error" role="alert">{{ error }}</p>
    <p v-else-if="loading" class="mt-8 text-base-content/70">{{ t('common.loading') }}</p>

    <div v-else class="mt-8 overflow-x-auto">
      <table class="table">
        <thead>
          <tr>
            <th>{{ t('projects.columns.title') }}</th>
            <th>{{ t('projects.columns.teams') }}</th>
            <th>{{ t('projects.columns.tags') }}</th>
            <th>{{ t('projects.columns.createdAt') }}</th>
            <th>
              <span class="sr-only">{{ t('projects.columns.actions') }}</span>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="projects.length === 0">
            <td colspan="5" class="text-base-content/70">
              {{ hasActiveSearch ? t('projects.emptySearch') : t('projects.empty') }}
            </td>
          </tr>
          <tr v-for="project in projects" :key="project.id">
            <td>
              <RouterLink
                class="link link-hover font-medium"
                :to="{ name: 'project-detail', params: { slug: project.slug } }"
              >
                {{ project.title }}
              </RouterLink>
            </td>
            <td>
              <span class="flex flex-wrap gap-1">
                <span v-for="team in project.teams" :key="team.id" class="badge badge-ghost badge-sm">
                  {{ team.name }}
                </span>
                <span v-if="project.teams.length === 0">{{ t('common.emDash') }}</span>
              </span>
            </td>
            <td>
              <span class="flex flex-wrap gap-1">
                <span
                  v-for="tag in project.tags"
                  :key="tag.id"
                  class="badge badge-sm gap-1 border border-base-300"
                >
                  <span
                    class="inline-block size-2 rounded-full"
                    :style="{ backgroundColor: tag.color }"
                    aria-hidden="true"
                  />
                  {{ tag.title }}
                </span>
                <span v-if="project.tags.length === 0">{{ t('common.emDash') }}</span>
              </span>
            </td>
            <td>{{ formatDate(project.created_at) }}</td>
            <td class="text-right">
              <div class="flex justify-end gap-1">
                <button
                  class="btn btn-ghost btn-square btn-sm"
                  type="button"
                  :disabled="mutating"
                  :aria-label="t('projects.edit')"
                  @click="openEdit(project)"
                >
                  <PencilSquareIcon class="size-4" aria-hidden="true" />
                </button>
                <button
                  v-if="canDelete"
                  class="btn btn-error btn-outline btn-square btn-sm"
                  type="button"
                  :disabled="mutating"
                  :aria-label="t('projects.delete')"
                  @click="removeProject(project)"
                >
                  <TrashIcon class="size-4" aria-hidden="true" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <div class="mt-4 flex flex-wrap items-center justify-between gap-3">
        <p class="text-sm text-base-content/70">
          {{ t('projects.pagination.pageOf', { page, totalPages, total }) }}
        </p>
        <div class="flex gap-2">
          <button class="btn btn-outline btn-sm" type="button" :disabled="!canPrev" @click="goPrev">
            {{ t('projects.pagination.prev') }}
          </button>
          <button class="btn btn-outline btn-sm" type="button" :disabled="!canNext" @click="goNext">
            {{ t('projects.pagination.next') }}
          </button>
        </div>
      </div>
    </div>

    <ProjectFormModal
      :open="formOpen"
      :mode="formMode"
      :title="formTitle"
      :team-ids="formTeamIds"
      :environment-ids="formEnvironmentIds"
      :tag-ids="formTagIds"
      :team-options="teamOptions"
      :environment-options="environmentOptions"
      :tag-options="tagOptions"
      :saving="mutating"
      :can-add-team="canAddTeam"
      :can-add-environment="canAddEnvironment"
      :can-add-tag="canAddTag"
      :catalog-quick="catalogQuick"
      :quick-users="quickUsers"
      :quick-saving="quickSaving"
      :quick-error="quickError"
      :team-quick-name="teamQuickName"
      :team-quick-accountable-id="teamQuickAccountableId"
      :environment-quick-title="environmentQuickTitle"
      :tag-quick-title="tagQuickTitle"
      :tag-quick-color="tagQuickColor"
      @close="closeForm"
      @submit="submitForm"
      @update:title="formTitle = $event"
      @update:team-ids="formTeamIds = $event"
      @update:environment-ids="formEnvironmentIds = $event"
      @update:tag-ids="formTagIds = $event"
      @add-team="openQuickTeam"
      @add-environment="openQuickEnvironment"
      @add-tag="openQuickTag"
      @close-catalog-quick="closeCatalogQuick"
      @submit-quick-team="submitQuickTeam"
      @submit-quick-environment="submitQuickEnvironment"
      @submit-quick-tag="submitQuickTag"
      @update:team-quick-name="teamQuickName = $event"
      @update:team-quick-accountable-id="teamQuickAccountableId = $event"
      @update:environment-quick-title="environmentQuickTitle = $event"
      @update:tag-quick-title="tagQuickTitle = $event"
      @update:tag-quick-color="tagQuickColor = $event"
    />
  </section>
</template>
