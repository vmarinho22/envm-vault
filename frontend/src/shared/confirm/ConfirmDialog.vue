<script setup lang="ts">
import { nextTick, watch } from 'vue';
import { useConfirmDialog } from '@/shared/confirm/useConfirmDialog';

const { open, title, message, confirmLabel, cancelLabel, variant, accept, cancel } =
  useConfirmDialog();

watch(
  open,
  async (isOpen) => {
    await nextTick();
    const dialog = document.getElementById('app-confirm-dialog') as HTMLDialogElement | null;
    if (!dialog) return;
    if (isOpen) {
      if (!dialog.open) dialog.showModal();
    } else if (dialog.open) {
      dialog.close();
    }
  },
);

function onDialogClose(): void {
  if (open.value) {
    cancel();
  }
}
</script>

<template>
  <dialog id="app-confirm-dialog" class="modal" @close="onDialogClose">
    <div class="modal-box max-w-md">
      <h3 class="text-lg font-medium">{{ title }}</h3>
      <p class="mt-3 text-sm text-base-content/80">{{ message }}</p>
      <div class="modal-action">
        <button class="btn" type="button" data-testid="confirm-cancel" @click="cancel">
          {{ cancelLabel }}
        </button>
        <button
          class="btn"
          :class="variant === 'danger' ? 'btn-error' : 'btn-primary'"
          type="button"
          data-testid="confirm-accept"
          @click="accept"
        >
          {{ confirmLabel }}
        </button>
      </div>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button type="submit">{{ cancelLabel }}</button>
    </form>
  </dialog>
</template>
