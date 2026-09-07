<script setup lang="ts">
import { computed, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import type { UserWithToken } from '@/modules/user/types';
import { fromName } from '@/shared/slug/fromName';

const props = defineProps<{
  open: boolean;
  mode: 'create' | 'edit';
  name: string;
  accountableId: string;
  users: UserWithToken[];
  saving: boolean;
}>();

const emit = defineEmits<{
  close: [];
  submit: [];
  'update:name': [value: string];
  'update:accountableId': [value: string];
}>();

const { t } = useI18n();

const slugPreview = computed(() => fromName(props.name));

watch(
  () => props.open,
  (isOpen) => {
    const dialog = document.getElementById('team-form-modal') as HTMLDialogElement | null;
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
  <dialog id="team-form-modal" class="modal" @close="onClose">
    <div class="modal-box max-w-md">
      <h3 class="text-lg font-medium">
        {{ mode === 'create' ? t('teams.form.createTitle') : t('teams.form.editTitle') }}
      </h3>

      <form class="mt-4 space-y-4" @submit.prevent="emit('submit')">
        <label class="form-control w-full">
          <span class="label-text">{{ t('teams.form.name') }}</span>
          <input
            class="input input-bordered w-full"
            type="text"
            name="team-name"
            autocomplete="off"
            required
            maxlength="80"
            :value="name"
            :aria-label="t('teams.form.name')"
            @input="emit('update:name', ($event.target as HTMLInputElement).value)"
          />
        </label>

        <p class="text-sm text-base-content/60">
          {{ t('teams.form.slugPreview') }}:
          <span class="font-mono">{{ slugPreview || t('common.emDash') }}</span>
        </p>

        <label class="form-control w-full">
          <span class="label-text">{{ t('teams.form.accountable') }}</span>
          <select
            class="select select-bordered w-full"
            name="team-accountable"
            autocomplete="off"
            required
            :value="accountableId"
            :aria-label="t('teams.form.accountable')"
            @change="emit('update:accountableId', ($event.target as HTMLSelectElement).value)"
          >
            <option disabled value="">{{ t('teams.form.selectUser') }}</option>
            <option v-for="user in users" :key="user.id" :value="user.id">
              {{ user.email }}
            </option>
          </select>
        </label>

        <div class="modal-action">
          <button class="btn" type="button" @click="onClose">{{ t('teams.form.cancel') }}</button>
          <button class="btn btn-primary" type="submit" :disabled="saving">
            <span v-if="saving" class="loading loading-spinner loading-xs" />
            {{ mode === 'create' ? t('teams.form.create') : t('teams.form.save') }}
          </button>
        </div>
      </form>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button type="submit">{{ t('teams.form.cancel') }}</button>
    </form>
  </dialog>
</template>
