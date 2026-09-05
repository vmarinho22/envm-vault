<script setup lang="ts">
import { EnvelopeIcon, KeyIcon } from '@heroicons/vue/24/outline';
import { useI18n } from 'vue-i18n';
import OtpCodeInput from '@/modules/auth/components/OtpCodeInput.vue';
import { useLogin } from '@/modules/auth/composables/useLogin';

const { t } = useI18n();
const { step, email, code, loading, error, submitEmail, submitCode, backToEmail } = useLogin();
</script>

<template>
  <main class="vault-bg flex min-h-svh items-end px-6 pb-16 md:items-center md:px-16 md:pb-0">
    <section class="mx-auto w-full max-w-xl">
      <p
        class="font-brand vault-enter text-5xl leading-none tracking-tight text-base-content md:text-7xl"
      >
        envm-vault
      </p>
      <h1
        class="vault-enter-delay mt-8 max-w-md text-2xl font-medium leading-snug text-base-content md:text-3xl"
      >
        {{ t('login.title') }}
      </h1>
      <p class="vault-enter-delay mt-3 max-w-md text-base text-base-content/70">
        {{ t('login.subtitle') }}
      </p>

      <form
        v-if="step === 'email'"
        class="vault-step mt-10 flex max-w-md flex-col gap-4"
        @submit.prevent="submitEmail"
      >
        <label class="flex flex-col gap-2">
          <span class="text-sm font-medium">{{ t('login.email') }}</span>
          <input
            v-model="email"
            class="input input-bordered w-full"
            type="email"
            name="email"
            autocomplete="email"
            required
            :disabled="loading"
          />
        </label>
        <p v-if="error" class="text-sm text-error" role="alert">{{ error }}</p>
        <button class="btn btn-primary w-fit gap-2" type="submit" :disabled="loading">
          <span v-if="loading" class="loading loading-spinner loading-sm" />
          <EnvelopeIcon v-else class="size-5" />
          {{ t('login.sendCode') }}
        </button>
      </form>

      <form
        v-else
        class="vault-step mt-10 flex max-w-md flex-col gap-4"
        @submit.prevent="submitCode"
      >
        <div class="flex flex-col gap-2">
          <span class="text-sm font-medium">{{ t('code.label') }}</span>
          <OtpCodeInput v-model="code" :disabled="loading" />
        </div>
        <p v-if="error" class="text-sm text-error" role="alert">{{ error }}</p>
        <div class="flex flex-wrap items-center gap-4">
          <button
            class="btn btn-primary w-fit gap-2"
            type="submit"
            :disabled="loading || code.length !== 8"
          >
            <span v-if="loading" class="loading loading-spinner loading-sm" />
            <KeyIcon v-else class="size-5" />
            {{ t('code.submit') }}
          </button>
          <button
            class="btn btn-ghost btn-sm"
            type="button"
            :disabled="loading"
            @click="backToEmail"
          >
            {{ t('code.backToEmail') }}
          </button>
        </div>
      </form>
    </section>
  </main>
</template>
