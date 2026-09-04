<script setup lang="ts">
import { ArrowRightStartOnRectangleIcon, HomeIcon } from '@heroicons/vue/24/outline';
import { RouterLink, RouterView, useRouter } from 'vue-router';
import { useAuthStore } from '@/modules/auth/stores/auth.store';

const auth = useAuthStore();
const router = useRouter();

function logout(): void {
  auth.logout();
  void router.push({ name: 'login' });
}
</script>

<template>
  <div class="vault-bg flex min-h-svh">
    <aside
      class="flex w-56 shrink-0 flex-col border-r border-base-300/80 bg-base-200/50 md:w-60"
      aria-label="Menu principal"
    >
      <div class="flex h-14 items-center px-5">
        <p class="font-brand text-xl leading-none tracking-tight">envm-vault</p>
      </div>

      <nav class="flex flex-1 flex-col gap-1 px-3 pt-2">
        <RouterLink
          :to="{ name: 'home' }"
          class="flex items-center gap-3 rounded-md px-3 py-2.5 text-sm text-base-content/75 transition-colors hover:bg-base-300/60 hover:text-base-content"
          active-class="!bg-base-300 !text-primary"
        >
          <HomeIcon class="size-5 shrink-0" aria-hidden="true" />
          Início
        </RouterLink>
      </nav>
    </aside>

    <div class="flex min-w-0 flex-1 flex-col">
      <header
        class="flex h-14 shrink-0 items-center justify-end gap-4 border-b border-base-300/80 bg-base-200/40 px-6 md:px-8"
      >
        <p v-if="auth.user" class="truncate text-sm text-base-content/70">
          {{ auth.user.email }}
        </p>
        <button class="btn btn-ghost btn-sm gap-2" type="button" @click="logout">
          <ArrowRightStartOnRectangleIcon class="size-4" aria-hidden="true" />
          Sair
        </button>
      </header>

      <main class="min-h-0 flex-1 overflow-auto">
        <RouterView />
      </main>
    </div>
  </div>
</template>
