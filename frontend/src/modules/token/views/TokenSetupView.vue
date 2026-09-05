<script setup lang="ts">
import { ArrowRightIcon, ArrowRightStartOnRectangleIcon, KeyIcon } from '@heroicons/vue/24/outline';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import LocaleSwitcher from '@/app/components/LocaleSwitcher.vue';
import { useAuthStore } from '@/modules/auth/stores/auth.store';
import TokenOnceReveal from '@/modules/token/components/TokenOnceReveal.vue';
import { useTokenSetup } from '@/modules/token/composables/useTokenSetup';

const { t } = useI18n();
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
      <div class="flex items-center gap-3">
        <LocaleSwitcher />
        <button class="btn btn-ghost btn-sm gap-2" type="button" data-testid="logout" @click="logout">
          <ArrowRightStartOnRectangleIcon class="size-4" aria-hidden="true" />
          {{ t('common.logout') }}
        </button>
      </div>
    </header>

    <main class="flex flex-1 items-end px-6 pb-16 md:items-center md:px-16 md:pb-0">
      <section class="mx-auto w-full max-w-xl">
        <h1 class="text-2xl font-medium leading-snug md:text-3xl">{{ t('tokenSetup.title') }}</h1>
        <p class="mt-3 max-w-md text-base text-base-content/70">
          {{ t('tokenSetup.description') }}
        </p>

        <p v-if="error" class="mt-6 text-sm text-error" role="alert">{{ error }}</p>

        <div v-if="!plaintext" class="mt-10">
          <button class="btn btn-primary gap-2" type="button" :disabled="loading" @click="generate">
            <span v-if="loading" class="loading loading-spinner loading-sm" />
            <KeyIcon v-else class="size-5" aria-hidden="true" />
            {{ t('tokenSetup.generate') }}
          </button>
        </div>

        <div v-else class="mt-10 flex flex-col gap-6">
          <TokenOnceReveal :token="plaintext" />
          <button class="btn btn-primary w-fit gap-2" type="button" @click="continueToApp">
            {{ t('tokenSetup.continue') }}
            <ArrowRightIcon class="size-5" aria-hidden="true" />
          </button>
        </div>
      </section>
    </main>
  </div>
</template>
