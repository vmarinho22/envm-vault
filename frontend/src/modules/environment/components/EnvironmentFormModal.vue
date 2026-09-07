<script setup lang="ts">
import { watch } from 'vue';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
  open: boolean;
  mode: 'create' | 'edit';
  title: string;
  saving: boolean;
}>();

const emit = defineEmits<{
  close: [];
  submit: [];
  'update:title': [value: string];
}>();

const { t } = useI18n();

watch(
  () => props.open,
  (isOpen) => {
    const dialog = document.getElementById('environment-form-modal') as HTMLDialogElement | null;
    if (!dialog) return;
    if (isOpen) {
      if (!dialog.open) dialog.showModal();
    } else if (dialog.open) {
      dialog.close();
    }
  },
);

function onClose(): void {
  emit('close');
}
</script>

<template>
  <dialog id="environment-form-modal" class="modal" @close="onClose">
    <div class="modal-box max-w-md">
      <h3 class="text-lg font-medium">
        {{
          mode === 'create'
            ? t('environments.form.createTitle')
            : t('environments.form.editTitle')
        }}
      </h3>

      <form class="mt-4 space-y-4" @submit.prevent="emit('submit')">
        <label class="form-control w-full">
          <span class="label-text">{{ t('environments.form.title') }}</span>
          <input
            class="input input-bordered w-full"
            type="text"
            name="environment-title"
            autocomplete="off"
            required
            maxlength="80"
            :value="title"
            :aria-label="t('environments.form.title')"
            @input="emit('update:title', ($event.target as HTMLInputElement).value)"
          />
        </label>

        <div class="modal-action">
          <button class="btn" type="button" @click="onClose">
            {{ t('environments.form.cancel') }}
          </button>
          <button class="btn btn-primary" type="submit" :disabled="saving">
            <span v-if="saving" class="loading loading-spinner loading-xs" />
            {{ mode === 'create' ? t('environments.form.create') : t('environments.form.save') }}
          </button>
        </div>
      </form>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button type="submit">{{ t('environments.form.cancel') }}</button>
    </form>
  </dialog>
</template>
