<script setup lang="ts">
import { ArrowPathIcon } from '@heroicons/vue/24/outline';
import { watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { DEFAULT_TAG_COLOR, randomTagColor } from '@/modules/tag/types';

const props = defineProps<{
  open: boolean;
  mode: 'create' | 'edit';
  title: string;
  color: string;
  saving: boolean;
}>();

const emit = defineEmits<{
  close: [];
  submit: [];
  'update:title': [value: string];
  'update:color': [value: string];
}>();

const { t } = useI18n();

watch(
  () => props.open,
  (isOpen) => {
    const dialog = document.getElementById('tag-form-modal') as HTMLDialogElement | null;
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

function normalizeHex(raw: string): string {
  const value = raw.trim();
  if (/^#[0-9A-Fa-f]{6}$/.test(value)) {
    return value.toUpperCase();
  }
  return value;
}

function onColorPickerInput(event: Event): void {
  const value = (event.target as HTMLInputElement).value.toUpperCase();
  emit('update:color', value);
}

function onHexInput(event: Event): void {
  emit('update:color', normalizeHex((event.target as HTMLInputElement).value));
}

function randomizeColor(): void {
  emit('update:color', randomTagColor());
}
</script>

<template>
  <dialog id="tag-form-modal" class="modal" @close="onClose">
    <div class="modal-box max-w-md">
      <h3 class="text-lg font-medium">
        {{ mode === 'create' ? t('tags.form.createTitle') : t('tags.form.editTitle') }}
      </h3>

      <form class="mt-4 space-y-4" @submit.prevent="emit('submit')">
        <label class="form-control w-full">
          <span class="label-text">{{ t('tags.form.title') }}</span>
          <input
            class="input input-bordered w-full"
            type="text"
            name="tag-title"
            autocomplete="off"
            required
            maxlength="80"
            :value="title"
            :aria-label="t('tags.form.title')"
            @input="emit('update:title', ($event.target as HTMLInputElement).value)"
          />
        </label>

        <div class="form-control w-full">
          <span class="label-text">{{ t('tags.form.color') }}</span>
          <div class="mt-1 flex items-center gap-2">
            <input
              class="h-10 w-14 cursor-pointer rounded border border-base-300 bg-transparent p-1"
              type="color"
              name="tag-color-picker"
              :value="/^#[0-9A-Fa-f]{6}$/.test(color) ? color : DEFAULT_TAG_COLOR"
              :aria-label="t('tags.form.colorPicker')"
              @input="onColorPickerInput"
            />
            <input
              class="input input-bordered min-w-0 flex-1 font-mono uppercase"
              type="text"
              name="tag-color-hex"
              autocomplete="off"
              required
              maxlength="7"
              pattern="#[0-9A-Fa-f]{6}"
              :value="color"
              :aria-label="t('tags.form.colorHex')"
              @input="onHexInput"
            />
            <button
              class="btn btn-ghost btn-square"
              type="button"
              :aria-label="t('tags.form.randomizeColor')"
              :title="t('tags.form.randomizeColor')"
              @click="randomizeColor"
            >
              <ArrowPathIcon class="size-5" aria-hidden="true" />
            </button>
          </div>
        </div>

        <div class="modal-action">
          <button class="btn" type="button" @click="onClose">{{ t('tags.form.cancel') }}</button>
          <button class="btn btn-primary" type="submit" :disabled="saving">
            <span v-if="saving" class="loading loading-spinner loading-xs" />
            {{ mode === 'create' ? t('tags.form.create') : t('tags.form.save') }}
          </button>
        </div>
      </form>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button type="submit">{{ t('tags.form.cancel') }}</button>
    </form>
  </dialog>
</template>
