<script setup lang="ts">
import {
  ArrowPathIcon,
  ChevronDownIcon,
  ChevronUpIcon,
  EyeIcon,
} from '@heroicons/vue/24/outline';
import { VueDatePicker } from '@vuepic/vue-datepicker';
import { ptBR } from 'date-fns/locale';
import { computed, onMounted, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { actionLabel } from '@/modules/audit/actionLabel';
import AuditLogDetailModal from '@/modules/audit/components/AuditLogDetailModal.vue';
import UserMultiSelect from '@/modules/audit/components/UserMultiSelect.vue';
import { channelLabel } from '@/modules/audit/channelLabel';
import { useAuditLogs } from '@/modules/audit/composables/useAuditLogs';
import CatalogMultiSelect from '@/modules/project/components/CatalogMultiSelect.vue';
import { formatDateTime, toDateParam } from '@/shared/datetime/format';

const { t } = useI18n();
const {
  logs,
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
} = useAuditLogs();

const dateRange = ref<[Date, Date] | null>(null);
const weekDayNames = ['Dom', 'Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sáb'];
const datePickerTimeConfig = {
  enableTimePicker: false,
  startTime: [
    { hours: 0, minutes: 0, seconds: 0 },
    { hours: 0, minutes: 0, seconds: 0 },
  ],
};

const sortAria = computed(() =>
  query.value.sort === 'asc' ? t('audit.sort.asc') : t('audit.sort.desc'),
);

function projectLabel(log: { project_title: string; project_slug: string }): string {
  if (log.project_title) return log.project_title;
  if (log.project_slug) return log.project_slug;
  return '—';
}

watch(dateRange, async (range) => {
  // Only filter when cleared or when both ends of the range are set.
  if (range == null) {
    filters.value.occurredFrom = '';
    filters.value.occurredTo = '';
    await applyFilters();
    return;
  }
  if (range.length < 2 || range[0] == null || range[1] == null) {
    return;
  }
  // Date-only filter: always treat bounds as local midnight (00:00).
  filters.value.occurredFrom = toDateParam(range[0]);
  filters.value.occurredTo = toDateParam(range[1]);
  await applyFilters();
});

async function onActorsChange(ids: string[]): Promise<void> {
  filters.value.actorUserIds = ids;
  await applyFilters();
}

async function onTargetsChange(ids: string[]): Promise<void> {
  filters.value.targetUserIds = ids;
  await applyFilters();
}

async function onProjectsChange(ids: string[]): Promise<void> {
  filters.value.projectIds = ids;
  await applyFilters();
}

onMounted(() => {
  void init();
});
</script>

<template>
  <section class="px-8 py-10 md:px-12">
    <div class="flex flex-wrap items-start justify-between gap-4">
      <div>
        <h1 class="text-2xl font-medium tracking-tight md:text-3xl">{{ t('audit.title') }}</h1>
        <p class="mt-3 max-w-lg text-base text-base-content/70">
          {{ t('audit.description') }}
        </p>
      </div>
      <button
        class="btn btn-outline btn-sm gap-2"
        type="button"
        :disabled="loading"
        data-testid="audit-refresh"
        @click="refresh"
      >
        <ArrowPathIcon class="size-4" aria-hidden="true" />
        {{ t('audit.refresh') }}
      </button>
    </div>

    <div class="mt-8 flex flex-wrap items-end gap-3">
      <UserMultiSelect
        :model-value="filters.actorUserIds"
        :options="userOptions"
        :label="t('audit.filters.actor')"
        @update:model-value="onActorsChange"
      />
      <UserMultiSelect
        :model-value="filters.targetUserIds"
        :options="userOptions"
        :label="t('audit.filters.target')"
        @update:model-value="onTargetsChange"
      />
      <div class="w-full max-w-xs">
        <CatalogMultiSelect
          :model-value="filters.projectIds"
          :options="projectOptions"
          :label="t('audit.filters.project')"
          :placeholder="t('audit.filters.allProjects')"
          @update:model-value="onProjectsChange"
        />
      </div>
      <div class="form-control flex w-full max-w-xs flex-col gap-[4px]">
        <span class="label-text">{{ t('audit.filters.dateRange') }}</span>
        <VueDatePicker
          v-model="dateRange"
          :range="{ partialRange: false }"
          dark
          :locale="ptBR"
          :day-names="weekDayNames"
          :week-start="0"
          :time-config="datePickerTimeConfig"
          auto-apply
          format="dd/MM/yyyy"
          :placeholder="t('audit.filters.dateRange')"
          class="audit-datepicker"
        />
      </div>
    </div>

    <p v-if="error" class="mt-6 text-sm text-error" role="alert">{{ error }}</p>
    <p v-else-if="loading" class="mt-8 text-base-content/70">{{ t('common.loading') }}</p>

    <template v-else>
      <div class="mt-8 overflow-x-auto">
        <table class="table">
          <thead>
            <tr>
              <th>{{ t('audit.columns.id') }}</th>
              <th>{{ t('audit.columns.action') }}</th>
              <th>{{ t('audit.columns.project') }}</th>
              <th>{{ t('audit.columns.secret') }}</th>
              <th>{{ t('audit.columns.actor') }}</th>
              <th>{{ t('audit.columns.target') }}</th>
              <th>{{ t('audit.columns.channel') }}</th>
              <th :aria-sort="query.sort === 'asc' ? 'ascending' : 'descending'">
                <button
                  class="btn btn-ghost btn-xs gap-1 px-0 font-bold"
                  type="button"
                  :aria-label="t('audit.sort.toggleAria')"
                  :title="sortAria"
                  @click="toggleSort"
                >
                  {{ t('audit.columns.occurredAt') }}
                  <ChevronUpIcon
                    v-if="query.sort === 'asc'"
                    class="size-4"
                    aria-hidden="true"
                  />
                  <ChevronDownIcon v-else class="size-4" aria-hidden="true" />
                </button>
              </th>
              <th>
                <span class="sr-only">{{ t('audit.columns.actions') }}</span>
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="logs.length === 0">
              <td colspan="9" class="text-base-content/70">{{ t('audit.empty') }}</td>
            </tr>
            <tr v-for="log in logs" :key="log.id">
              <td class="max-w-[8rem] truncate font-mono text-xs" :title="log.id">
                {{ log.id }}
              </td>
              <td class="text-sm">{{ actionLabel(log.action, t) }}</td>
              <td class="text-sm">{{ projectLabel(log) }}</td>
              <td class="font-mono text-sm">{{ log.secret_name || '—' }}</td>
              <td>{{ log.actor_email || '—' }}</td>
              <td>{{ log.target_email || '—' }}</td>
              <td>{{ channelLabel(log.auth_channel, t) }}</td>
              <td class="whitespace-nowrap text-sm">{{ formatDateTime(log.occurred_at) }}</td>
              <td class="text-right">
                <button
                  class="btn btn-ghost btn-sm btn-square"
                  type="button"
                  :aria-label="t('audit.detail.openAria')"
                  @click="openDetail(log)"
                >
                  <EyeIcon class="size-4" aria-hidden="true" />
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="mt-6 flex items-center justify-between gap-4">
        <p class="text-sm text-base-content/70">
          {{ t('audit.pagination.pageOf', { page: query.page, totalPages, total }) }}
        </p>
        <div class="flex gap-2">
          <button class="btn btn-outline btn-sm" type="button" :disabled="!canPrev" @click="goPrev">
            {{ t('audit.pagination.prev') }}
          </button>
          <button class="btn btn-outline btn-sm" type="button" :disabled="!canNext" @click="goNext">
            {{ t('audit.pagination.next') }}
          </button>
        </div>
      </div>
    </template>

    <AuditLogDetailModal :open="detailOpen" :log="selectedLog" @close="closeDetail" />
  </section>
</template>
