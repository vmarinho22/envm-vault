<script setup lang="ts">
import { EnvelopeIcon, KeyIcon } from '@heroicons/vue/24/outline';
import { useLogin } from '@/modules/auth/composables/useLogin';

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
        Entre para gerenciar os ambientes.
      </h1>
      <p class="vault-enter-delay mt-3 max-w-md text-base text-base-content/70">
        Enviamos um código de 8 dígitos para o seu e-mail.
      </p>

      <form
        v-if="step === 'email'"
        class="vault-step mt-10 flex max-w-md flex-col gap-4"
        @submit.prevent="submitEmail"
      >
        <label class="flex flex-col gap-2">
          <span class="text-sm font-medium">E-mail</span>
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
          Enviar código
        </button>
      </form>

      <form
        v-else
        class="vault-step mt-10 flex max-w-md flex-col gap-4"
        @submit.prevent="submitCode"
      >
        <label class="flex flex-col gap-2">
          <span class="text-sm font-medium">Código</span>
          <input
            v-model="code"
            class="input input-bordered w-full tracking-[0.35em]"
            type="text"
            name="code"
            inputmode="numeric"
            autocomplete="one-time-code"
            maxlength="8"
            minlength="8"
            pattern="[0-9]{8}"
            required
            :disabled="loading"
          />
        </label>
        <p v-if="error" class="text-sm text-error" role="alert">{{ error }}</p>
        <div class="flex flex-wrap items-center gap-4">
          <button class="btn btn-primary w-fit gap-2" type="submit" :disabled="loading">
            <span v-if="loading" class="loading loading-spinner loading-sm" />
            <KeyIcon v-else class="size-5" />
            Entrar
          </button>
          <button
            class="btn btn-ghost btn-sm"
            type="button"
            :disabled="loading"
            @click="backToEmail"
          >
            Usar outro e-mail
          </button>
        </div>
      </form>
    </section>
  </main>
</template>
