<script setup lang="ts">
import { ArrowPathIcon, TrashIcon } from '@heroicons/vue/24/outline';
import { onMounted } from 'vue';
import TokenOnceReveal from '@/modules/token/components/TokenOnceReveal.vue';
import { useApiKey } from '@/modules/token/composables/useApiKey';

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
  if (!value) return '—';
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
    <h1 class="text-2xl font-medium tracking-tight md:text-3xl">Chave de API</h1>
    <p class="mt-3 max-w-lg text-base text-base-content/70">
      O valor completo só aparece na geração ou ao rotacionar. Rotacionar invalida a chave antiga.
    </p>

    <p v-if="error" class="mt-6 text-sm text-error" role="alert">{{ error }}</p>
    <p v-else-if="loading" class="mt-8 text-base-content/70">Carregando...</p>

    <div v-else class="mt-8 flex max-w-xl flex-col gap-6">
      <dl class="grid gap-3 text-sm">
        <div class="flex flex-col gap-1">
          <dt class="text-base-content/60">Prefixo</dt>
          <dd class="font-mono">{{ prefix ?? '—' }}</dd>
        </div>
        <div class="flex flex-col gap-1">
          <dt class="text-base-content/60">Criada em</dt>
          <dd>{{ formatDate(createdAt) }}</dd>
        </div>
        <div class="flex flex-col gap-1">
          <dt class="text-base-content/60">Status</dt>
          <dd>{{ isActive ? 'Ativa' : 'Revogada' }}</dd>
        </div>
      </dl>

      <template v-if="plaintext">
        <TokenOnceReveal :token="plaintext" />
        <button class="btn btn-primary w-fit" type="button" @click="dismissPlaintext">
          Já guardei a chave
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
          Rotacionar
        </button>
        <button class="btn btn-error btn-outline gap-2" type="button" :disabled="mutating" @click="revoke">
          <TrashIcon class="size-5" aria-hidden="true" />
          Revogar
        </button>
      </div>
    </div>
  </section>
</template>
