<script setup lang="ts">
import { ArrowPathIcon, KeyIcon } from '@heroicons/vue/24/outline';
import { onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import TokenOnceReveal from '@/modules/token/components/TokenOnceReveal.vue';
import { useApiKey } from '@/modules/token/composables/useApiKey';

const { t } = useI18n();

const {
  loading,
  mutating,
  error,
  prefix,
  createdAt,
  isActive,
  plaintext,
  load,
  rotate,
  revoke,
  dismissPlaintext,
} = useApiKey();

function formatDate(value: string | null): string {
  if (!value) return t('common.emDash');
  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) return value;
  return `${parsed.toISOString().slice(0, 16).replace('T', ' ')} UTC`;
}

onMounted(() => {
  void load();
});
</script>

<template>
  <section class="px-8 py-10 md:px-12">
    <h1 class="text-2xl font-medium tracking-tight md:text-3xl">{{ t('apiKey.title') }}</h1>
    <p class="mt-3 max-w-lg text-base text-base-content/70">
      {{ t('apiKey.description') }}
    </p>

    <p v-if="error" class="mt-6 text-sm text-error" role="alert">{{ error }}</p>
    <p v-else-if="loading" class="mt-8 text-base-content/70">{{ t('common.loading') }}</p>

    <div v-else class="mt-8 flex max-w-xl flex-col gap-6">
      <dl class="grid gap-3 text-sm">
        <div class="flex flex-col gap-1">
          <dt class="text-base-content/60">{{ t('apiKey.prefix') }}</dt>
          <dd class="font-mono">{{ prefix ?? t('common.emDash') }}</dd>
        </div>
        <div class="flex flex-col gap-1">
          <dt class="text-base-content/60">{{ t('apiKey.createdAt') }}</dt>
          <dd>{{ formatDate(createdAt) }}</dd>
        </div>
        <div class="flex flex-col gap-1">
          <dt class="text-base-content/60">{{ t('apiKey.status') }}</dt>
          <dd>{{ isActive ? t('apiKey.active') : t('apiKey.revoked') }}</dd>
        </div>
      </dl>

      <template v-if="plaintext">
        <TokenOnceReveal :token="plaintext" />
        <button class="btn btn-primary w-fit" type="button" @click="dismissPlaintext">
          {{ t('apiKey.dismissPlaintext') }}
        </button>
      </template>

      <div v-else class="flex flex-wrap gap-3">
        <button
          class="btn btn-outline gap-2"
          type="button"
          data-testid="rotate-token"
          :disabled="mutating"
          @click="rotate"
        >
          <span v-if="mutating" class="loading loading-spinner loading-sm" />
          <ArrowPathIcon v-else class="size-5" aria-hidden="true" />
          {{ t('apiKey.rotate') }}
        </button>
        <button class="btn btn-error btn-outline gap-2" type="button" :disabled="mutating" @click="revoke">
          <KeyIcon class="size-5" aria-hidden="true" />
          {{ t('apiKey.revoke') }}
        </button>
      </div>
    </div>
  </section>
</template>
