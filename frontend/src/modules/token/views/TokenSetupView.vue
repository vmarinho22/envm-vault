<script setup lang="ts">
import { ArrowRightIcon, ArrowRightStartOnRectangleIcon, KeyIcon } from '@heroicons/vue/24/outline';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/modules/auth/stores/auth.store';
import TokenOnceReveal from '@/modules/token/components/TokenOnceReveal.vue';
import { useTokenSetup } from '@/modules/token/composables/useTokenSetup';

const auth = useAuthStore();
const router = useRouter();
const { loading, error, plaintext, generate, continueToApp } = useTokenSetup();

function logout(): void {
  auth.logout();
  void router.push({ name: 'login' });
}
</script>

<template>
  <div class="vault-bg flex min-h-svh flex-col">
    <header class="flex h-14 items-center justify-between px-6 md:px-10">
      <p class="font-brand text-xl tracking-tight">envm-vault</p>
      <button class="btn btn-ghost btn-sm gap-2" type="button" @click="logout">
        <ArrowRightStartOnRectangleIcon class="size-4" aria-hidden="true" />
        Sair
      </button>
    </header>

    <main class="flex flex-1 items-end px-6 pb-16 md:items-center md:px-16 md:pb-0">
      <section class="mx-auto w-full max-w-xl">
        <h1 class="text-2xl font-medium leading-snug md:text-3xl">Gere sua chave de API</h1>
        <p class="mt-3 max-w-md text-base text-base-content/70">
          Esta chave identifica você nas requisições de variáveis de ambiente. Você só verá o valor
          completo uma vez.
        </p>

        <p v-if="error" class="mt-6 text-sm text-error" role="alert">{{ error }}</p>

        <div v-if="!plaintext" class="mt-10">
          <button class="btn btn-primary gap-2" type="button" :disabled="loading" @click="generate">
            <span v-if="loading" class="loading loading-spinner loading-sm" />
            <KeyIcon v-else class="size-5" aria-hidden="true" />
            Gerar chave
          </button>
        </div>

        <div v-else class="mt-10 flex flex-col gap-6">
          <TokenOnceReveal :token="plaintext" />
          <button class="btn btn-primary w-fit gap-2" type="button" @click="continueToApp">
            Continuar
            <ArrowRightIcon class="size-5" aria-hidden="true" />
          </button>
        </div>
      </section>
    </main>
  </div>
</template>
