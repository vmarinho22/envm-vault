<script setup lang="ts">
import { ClipboardDocumentIcon } from '@heroicons/vue/24/outline';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
  token: string;
}>();

const { t } = useI18n();
const copied = ref(false);

async function copyToken(): Promise<void> {
  try {
    await navigator.clipboard.writeText(props.token);
    copied.value = true;
  } catch {
    copied.value = false;
  }
}
</script>

<template>
  <div class="rounded-lg border border-warning/40 bg-warning/10 p-4">
    <p class="text-sm font-medium text-warning">
      {{ t('tokenReveal.warning') }}
    </p>
    <div class="mt-3 flex flex-col gap-2 sm:flex-row sm:items-center">
      <input
        class="input input-bordered min-w-0 flex-1 font-mono text-sm"
        data-testid="api-token-value"
        type="text"
        readonly
        :value="token"
      />
      <button
        class="btn btn-outline btn-sm gap-2"
        type="button"
        data-testid="copy-token"
        @click.stop.prevent="copyToken"
      >
        <ClipboardDocumentIcon class="size-4" aria-hidden="true" />
        {{ copied ? t('common.copied') : t('common.copy') }}
      </button>
    </div>
  </div>
</template>
