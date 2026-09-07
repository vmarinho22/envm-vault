<script setup lang="ts">
import { computed, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import CatalogMultiSelect, {
  type CatalogOption,
} from '@/modules/project/components/CatalogMultiSelect.vue';

const props = defineProps<{
  open: boolean;
  mode: 'create' | 'edit';
  name: string;
  saving: boolean;
  environmentOptions: CatalogOption[];
  selectedEnvironmentIds: string[];
  valuesByEnv: Record<string, string>;
}>();

const emit = defineEmits<{
  close: [];
  submit: [];
  'update:name': [value: string];
  'update:selectedEnvironmentIds': [value: string[]];
  'update:value': [payload: { environmentId: string; value: string }];
}>();

const { t } = useI18n();

const selectedEnvs = computed(() =>
  props.environmentOptions.filter((o) => props.selectedEnvironmentIds.includes(o.id)),
);

const canSubmit = computed(() => {
  if (!props.name.trim()) return false;
  if (props.selectedEnvironmentIds.length === 0) return false;
  return props.selectedEnvironmentIds.every((id) => (props.valuesByEnv[id] ?? '').length > 0);
});

watch(
  () => props.open,
  (open) => {
    const dialog = document.getElementById('secret-form-modal') as HTMLDialogElement | null;
    if (!dialog) return;
    if (open && !dialog.open) dialog.showModal();
    if (!open && dialog.open) dialog.close();
  },
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

        <div v-if="selectedEnvs.length > 0" class="space-y-3">
          <p class="text-sm font-medium text-base-content/80">{{ t('secrets.form.valuesTitle') }}</p>
          <label
            v-for="env in selectedEnvs"
            :key="env.id"
            class="form-control w-full"
          >
            <span class="label-text font-mono">{{ env.label }}</span>
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
