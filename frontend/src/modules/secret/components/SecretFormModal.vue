<script setup lang="ts">
import { PlusIcon } from '@heroicons/vue/24/outline';
import { computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import CatalogMultiSelect, {
  type CatalogOption,
} from '@/modules/project/components/CatalogMultiSelect.vue';

const props = defineProps<{
  open: boolean;
  mode: 'create' | 'edit';
  name: string;
  saving: boolean;
  projectCount?: number;
  environmentOptions: CatalogOption[];
  selectedEnvironmentIds: string[];
  valuesByEnv: Record<string, string>;
  prefixesByEnv: Record<string, string>;
}>();

const emit = defineEmits<{
  close: [];
  submit: [];
  'update:name': [value: string];
  'update:selectedEnvironmentIds': [value: string[]];
  'update:value': [payload: { environmentId: string; value: string }];
  'update:prefix': [payload: { environmentId: string; prefix: string }];
}>();

const { t } = useI18n();

/** Env ids whose prefix field is visible. */
const prefixOpenByEnv = ref<Record<string, boolean>>({});

const selectedEnvs = computed(() =>
  props.environmentOptions.filter((o) => props.selectedEnvironmentIds.includes(o.id)),
);

const canSubmit = computed(() => {
  if (!props.name.trim()) return false;
  if (props.selectedEnvironmentIds.length === 0) return false;
  return props.selectedEnvironmentIds.every((id) => (props.valuesByEnv[id] ?? '').length > 0);
});

function previewName(environmentId: string): string {
  const prefix = (props.prefixesByEnv[environmentId] ?? '').trim();
  const base = props.name.trim() || t('secrets.form.namePreviewPlaceholder');
  return `${prefix}${base}`;
}

function isPrefixVisible(environmentId: string): boolean {
  return Boolean(prefixOpenByEnv.value[environmentId]);
}

function showPrefix(environmentId: string): void {
  prefixOpenByEnv.value = { ...prefixOpenByEnv.value, [environmentId]: true };
}

function syncPrefixVisibility(): void {
  const next: Record<string, boolean> = {};
  for (const id of props.selectedEnvironmentIds) {
    const hasPrefix = Boolean((props.prefixesByEnv[id] ?? '').trim());
    next[id] = hasPrefix || Boolean(prefixOpenByEnv.value[id]);
  }
  prefixOpenByEnv.value = next;
}

watch(
  () => props.open,
  (open) => {
    const dialog = document.getElementById('secret-form-modal') as HTMLDialogElement | null;
    if (!dialog) return;
    if (open && !dialog.open) {
      prefixOpenByEnv.value = {};
      syncPrefixVisibility();
      dialog.showModal();
    }
    if (!open && dialog.open) {
      prefixOpenByEnv.value = {};
      dialog.close();
    }
  },
);

watch(
  () => [props.selectedEnvironmentIds, props.prefixesByEnv] as const,
  () => {
    if (!props.open) return;
    syncPrefixVisibility();
  },
  { deep: true },
);

function onClose(): void {
  emit('close');
}
</script>

<template>
  <dialog id="secret-form-modal" class="modal" @close="onClose">
    <div class="modal-box max-h-[90vh] max-w-lg overflow-y-auto">
      <h3 class="text-lg font-medium">
        {{ mode === 'create' ? t('secrets.form.createTitle') : t('secrets.form.editTitle') }}
      </h3>
      <p
        v-if="mode === 'edit' && (projectCount ?? 1) > 1"
        class="mt-2 text-sm text-warning"
        role="status"
      >
        {{ t('secrets.form.sharedWarning', { count: projectCount }) }}
      </p>
      <form class="mt-4 space-y-4" @submit.prevent="emit('submit')">
        <label class="form-control w-full">
          <span class="label-text">{{ t('secrets.form.name') }}</span>
          <input
            class="input input-bordered w-full font-mono"
            type="text"
            name="secret-name"
            autocomplete="off"
            required
            pattern="[A-Za-z_][A-Za-z0-9_]*"
            maxlength="128"
            :value="name"
            :disabled="saving"
            :aria-label="t('secrets.form.name')"
            @input="emit('update:name', ($event.target as HTMLInputElement).value)"
          />
          <span class="label-text-alt text-base-content/60">{{ t('secrets.form.nameHint') }}</span>
        </label>

        <CatalogMultiSelect
          :model-value="selectedEnvironmentIds"
          :options="environmentOptions"
          :label="t('secrets.form.environments')"
          :placeholder="t('secrets.form.selectEnvironments')"
          @update:model-value="emit('update:selectedEnvironmentIds', $event)"
        />

        <div v-if="selectedEnvs.length > 0" class="space-y-4">
          <p class="text-sm font-medium text-base-content/80">{{ t('secrets.form.valuesTitle') }}</p>
          <div
            v-for="env in selectedEnvs"
            :key="env.id"
            class="space-y-2 rounded-box border border-base-300 p-3"
          >
            <p class="font-mono text-sm font-medium">{{ env.label }}</p>
            <label class="form-control w-full">
              <span class="label-text">{{ t('secrets.form.value') }}</span>
              <input
                class="input input-bordered w-full font-mono text-sm"
                type="text"
                :name="`secret-value-${env.id}`"
                autocomplete="off"
                required
                :value="valuesByEnv[env.id] ?? ''"
                :disabled="saving"
                :aria-label="t('secrets.form.valueForEnv', { env: env.label })"
                @input="
                  emit('update:value', {
                    environmentId: env.id,
                    value: ($event.target as HTMLInputElement).value,
                  })
                "
              />
            </label>

            <label v-if="isPrefixVisible(env.id)" class="form-control w-full">
              <span class="label-text">{{ t('secrets.form.prefix') }}</span>
              <input
                class="input input-bordered w-full font-mono text-sm"
                type="text"
                :name="`secret-prefix-${env.id}`"
                autocomplete="off"
                pattern="[A-Za-z_][A-Za-z0-9_]*"
                maxlength="64"
                :value="prefixesByEnv[env.id] ?? ''"
                :disabled="saving"
                :placeholder="t('secrets.form.prefixPlaceholder')"
                :aria-label="t('secrets.form.prefixForEnv', { env: env.label })"
                @input="
                  emit('update:prefix', {
                    environmentId: env.id,
                    prefix: ($event.target as HTMLInputElement).value,
                  })
                "
              />
              <span class="label-text-alt text-base-content/60">
                {{ t('secrets.form.prefixHint') }}
                <span class="font-mono">{{ previewName(env.id) }}</span>
              </span>
            </label>
            <div v-else class="flex justify-end">
              <button
                class="btn btn-ghost btn-xs gap-1 text-base-content/70"
                type="button"
                :disabled="saving"
                @click="showPrefix(env.id)"
              >
                <PlusIcon class="size-3.5" aria-hidden="true" />
                {{ t('secrets.form.addPrefix') }}
              </button>
            </div>
          </div>
        </div>

        <div class="modal-action">
          <button class="btn" type="button" :disabled="saving" @click="onClose">
            {{ t('secrets.form.cancel') }}
          </button>
          <button class="btn btn-primary" type="submit" :disabled="saving || !canSubmit">
            {{ mode === 'create' ? t('secrets.form.create') : t('secrets.form.save') }}
          </button>
        </div>
      </form>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button type="submit">{{ t('secrets.form.cancel') }}</button>
    </form>
  </dialog>
</template>
