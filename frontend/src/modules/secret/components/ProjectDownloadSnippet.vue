<script setup lang="ts">
import { ClipboardDocumentIcon } from '@heroicons/vue/24/outline';
import { computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { useAuthStore } from '@/modules/auth/stores/auth.store';
import {
  buildDownloadSnippet,
  type DownloadOs,
} from '@/modules/secret/downloadSnippet';

const PRODUCTION_ENV = 'production';

const props = defineProps<{
  projectSlug: string;
  environments: { id: string; title: string }[];
}>();

const { t } = useI18n();
const auth = useAuthStore();

const osTabs: { id: DownloadOs; labelKey: string }[] = [
  { id: 'linux', labelKey: 'secrets.download.linux' },
  { id: 'macos', labelKey: 'secrets.download.macos' },
  { id: 'windows', labelKey: 'secrets.download.windows' },
];

const downloadableEnvironments = computed(() => {
  if (auth.isSuperAdmin) return props.environments;
  return props.environments.filter((e) => e.title !== PRODUCTION_ENV);
});

const selectedOs = ref<DownloadOs>('linux');
const selectedEnvTitle = ref(downloadableEnvironments.value[0]?.title ?? '');
const copied = ref(false);

watch(
  downloadableEnvironments,
  (envs) => {
    if (!envs.some((e) => e.title === selectedEnvTitle.value)) {
      selectedEnvTitle.value = envs[0]?.title ?? '';
    }
  },
  { immediate: true },
);

const snippet = computed(() => {
  if (!props.projectSlug || !selectedEnvTitle.value) return '';
  return buildDownloadSnippet(selectedOs.value, props.projectSlug, selectedEnvTitle.value);
});

const showProductionNote = computed(
  () =>
    !auth.isSuperAdmin && props.environments.some((e) => e.title === PRODUCTION_ENV),
);

async function copySnippet(): Promise<void> {
  if (!snippet.value) return;
  try {
    await navigator.clipboard.writeText(snippet.value);
    copied.value = true;
    window.setTimeout(() => {
      copied.value = false;
    }, 2000);
  } catch {
    copied.value = false;
  }
}
</script>

<template>
  <div class="rounded-lg border border-base-300 bg-base-200/40 p-4 md:p-5">
    <div class="flex flex-wrap items-end justify-between gap-3">
      <div class="min-w-0 flex-1">
        <h2 class="text-base font-medium">{{ t('secrets.download.title') }}</h2>
        <p class="mt-1 max-w-2xl text-sm text-base-content/70">
          {{ t('secrets.download.description') }}
        </p>
      </div>
      <div class="flex flex-wrap items-end gap-2">
        <label class="form-control w-40 sm:w-48">
          <span class="label-text">{{ t('secrets.download.environment') }}</span>
          <select
            v-model="selectedEnvTitle"
            class="select select-bordered select-sm"
            :disabled="downloadableEnvironments.length === 0"
            :aria-label="t('secrets.download.environment')"
          >
            <option v-if="downloadableEnvironments.length === 0" value="">
              {{ t('secrets.download.noEnvironments') }}
            </option>
            <option
              v-for="env in downloadableEnvironments"
              :key="env.id"
              :value="env.title"
            >
              {{ env.title }}
            </option>
          </select>
        </label>
        <button
          class="btn btn-outline btn-sm gap-2"
          type="button"
          data-testid="copy-download-snippet"
          :disabled="!snippet"
          @click="copySnippet"
        >
          <ClipboardDocumentIcon class="size-4" aria-hidden="true" />
          {{ copied ? t('common.copied') : t('common.copy') }}
        </button>
      </div>
    </div>

    <div role="tablist" class="tabs tabs-boxed mt-4 w-fit bg-base-300/50">
      <button
        v-for="tab in osTabs"
        :key="tab.id"
        type="button"
        role="tab"
        class="tab"
        :class="{ 'tab-active': selectedOs === tab.id }"
        :aria-selected="selectedOs === tab.id"
        @click="selectedOs = tab.id"
      >
        {{ t(tab.labelKey) }}
      </button>
    </div>

    <pre
      class="mt-3 max-h-56 overflow-auto rounded-md border border-base-300 bg-base-100 p-3 font-mono text-xs leading-relaxed"
      data-testid="download-snippet"
    >{{ snippet }}</pre>

    <p class="mt-3 text-xs text-base-content/60">
      {{ t('secrets.download.tokenHint') }}
    </p>
    <p
      v-if="showProductionNote"
      class="mt-1 text-xs text-base-content/60"
      data-testid="production-download-note"
    >
      {{ t('secrets.download.productionRestricted') }}
    </p>
  </div>
</template>
