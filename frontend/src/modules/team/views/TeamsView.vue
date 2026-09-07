<script setup lang="ts">
import { PencilSquareIcon, PlusIcon, TrashIcon } from '@heroicons/vue/24/outline';
import { computed, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import TeamFormModal from '@/modules/team/components/TeamFormModal.vue';
import { useTeams } from '@/modules/team/composables/useTeams';
import { formatDate } from '@/shared/datetime/format';

const { t } = useI18n();
const {
  teams,
  users,
  loading,
  mutating,
  error,
  nameQuery,
  hasActiveSearch,
  formOpen,
  editing,
  formName,
  formAccountableId,
  canMutate,
  init,
  applySearch,
  clearSearch,
  openCreate,
  openEdit,
  closeForm,
  submitForm,
  removeTeam,
} = useTeams();

const formMode = computed(() => (editing.value ? 'edit' : 'create'));

onMounted(() => {
  void init();
});
</script>

<template>
  <section class="px-8 py-10 md:px-12">
    <div class="flex flex-wrap items-start justify-between gap-4">
      <div>
        <h1 class="text-2xl font-medium tracking-tight md:text-3xl">{{ t('teams.title') }}</h1>
        <p class="mt-3 max-w-lg text-base text-base-content/70">
          {{ t('teams.description') }}
        </p>
      </div>
      <button
        v-if="canMutate"
        class="btn btn-primary gap-2"
        type="button"
        @click="openCreate"
      >
        <PlusIcon class="size-4" aria-hidden="true" />
        {{ t('teams.create') }}
      </button>
    </div>

    <form class="mt-8 flex flex-wrap items-end gap-3" @submit.prevent="applySearch">
      <div class="w-full max-w-xs">
        <input
          v-model="nameQuery"
          class="input input-bordered w-full"
          type="search"
          name="team-name-search"
          autocomplete="off"
          :placeholder="t('teams.search.placeholder')"
          :aria-label="t('teams.search.label')"
        />
      </div>
      <button class="btn btn-outline" type="submit">{{ t('teams.search.apply') }}</button>
      <button
        v-if="hasActiveSearch"
        class="btn btn-ghost"
        type="button"
        @click="clearSearch"
      >
        {{ t('teams.search.clear') }}
      </button>
    </form>

    <p v-if="error" class="mt-6 text-sm text-error" role="alert">{{ error }}</p>
    <p v-else-if="loading" class="mt-8 text-base-content/70">{{ t('common.loading') }}</p>

    <div v-else class="mt-8 overflow-x-auto">
      <table class="table">
        <thead>
          <tr>
            <th>{{ t('teams.columns.name') }}</th>
            <th>{{ t('teams.columns.slug') }}</th>
            <th>{{ t('teams.columns.accountable') }}</th>
            <th>{{ t('teams.columns.createdAt') }}</th>
            <th>{{ t('teams.columns.updatedAt') }}</th>
            <th v-if="canMutate">
              <span class="sr-only">{{ t('teams.columns.actions') }}</span>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="teams.length === 0">
            <td :colspan="canMutate ? 6 : 5" class="text-base-content/70">
              {{ hasActiveSearch ? t('teams.emptySearch') : t('teams.empty') }}
            </td>
          </tr>
          <tr v-for="team in teams" :key="team.id">
            <td>{{ team.name }}</td>
            <td class="font-mono text-sm">{{ team.slug }}</td>
            <td>{{ team.accountable.email }}</td>
            <td>{{ formatDate(team.created_at) }}</td>
            <td>{{ formatDate(team.updated_at) }}</td>
            <td v-if="canMutate" class="text-right">
              <div class="flex justify-end gap-1">
                <button
                  class="btn btn-ghost btn-square btn-sm"
                  type="button"
                  :disabled="mutating"
                  :aria-label="t('teams.edit')"
                  @click="openEdit(team)"
                >
                  <PencilSquareIcon class="size-4" aria-hidden="true" />
                </button>
                <button
                  class="btn btn-error btn-outline btn-square btn-sm"
                  type="button"
                  :disabled="mutating"
                  :aria-label="t('teams.delete')"
                  @click="removeTeam(team)"
                >
                  <TrashIcon class="size-4" aria-hidden="true" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <TeamFormModal
      v-if="canMutate"
      :open="formOpen"
      :mode="formMode"
      :name="formName"
      :accountable-id="formAccountableId"
      :users="users"
      :saving="mutating"
      @close="closeForm"
      @submit="submitForm"
      @update:name="formName = $event"
      @update:accountable-id="formAccountableId = $event"
    />
  </section>
</template>
