<script setup lang="ts">
import { watch } from 'vue';
import { useI18n } from 'vue-i18n';
import type { RegisterableRole } from '@/modules/user/services/UserService';
import { roleLabel } from '@/modules/user/roleLabel';

const props = defineProps<{
  open: boolean;
  email: string;
  role: RegisterableRole;
  saving: boolean;
}>();

const emit = defineEmits<{
  close: [];
  submit: [];
  'update:email': [value: string];
  'update:role': [value: RegisterableRole];
}>();

const { t } = useI18n();

const roles: RegisterableRole[] = ['admin', 'user'];

watch(
  () => props.open,
  (isOpen) => {
    const dialog = document.getElementById('user-form-modal') as HTMLDialogElement | null;
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

function onRoleChange(event: Event): void {
  const value = (event.target as HTMLSelectElement).value;
  if (value === 'admin' || value === 'user') {
    emit('update:role', value);
  }
}
</script>

<template>
  <dialog id="user-form-modal" class="modal" @close="onClose">
    <div class="modal-box max-w-md">
      <h3 class="text-lg font-medium">{{ t('users.form.createTitle') }}</h3>

      <form class="mt-4 space-y-4" @submit.prevent="emit('submit')">
        <label class="form-control flex w-full flex-col gap-[4px]">
          <span class="label-text">{{ t('users.form.email') }}</span>
          <input
            class="input input-bordered w-full"
            type="email"
            name="user-email"
            autocomplete="off"
            required
            :value="email"
            :aria-label="t('users.form.email')"
            @input="emit('update:email', ($event.target as HTMLInputElement).value)"
          />
        </label>

        <label class="form-control flex w-full flex-col gap-[4px]">
          <span class="label-text">{{ t('users.form.role') }}</span>
          <select
            class="select select-bordered w-full"
            name="user-role"
            autocomplete="off"
            required
            :value="role"
            :aria-label="t('users.form.role')"
            @change="onRoleChange"
          >
            <option v-for="option in roles" :key="option" :value="option">
              {{ roleLabel(option, t) }}
            </option>
          </select>
        </label>

        <div class="modal-action">
          <button class="btn" type="button" @click="onClose">{{ t('users.form.cancel') }}</button>
          <button class="btn btn-primary" type="submit" :disabled="saving">
            <span v-if="saving" class="loading loading-spinner loading-xs" />
            {{ t('users.form.create') }}
          </button>
        </div>
      </form>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button type="submit">{{ t('users.form.cancel') }}</button>
    </form>
  </dialog>
</template>
