<script setup lang="ts">
import {
  ArrowRightStartOnRectangleIcon,
  ClipboardDocumentListIcon,
  CubeTransparentIcon,
  HomeIcon,
  KeyIcon,
  TagIcon,
  UserGroupIcon,
  UsersIcon,
} from '@heroicons/vue/24/outline';
import { useI18n } from 'vue-i18n';
import { RouterLink, RouterView, useRouter } from 'vue-router';
import LocaleSwitcher from '@/app/components/LocaleSwitcher.vue';
import { useAuthStore } from '@/modules/auth/stores/auth.store';
import ConfirmDialog from '@/shared/confirm/ConfirmDialog.vue';

const { t } = useI18n();
const auth = useAuthStore();
const router = useRouter();

const navLinkClass =
  'flex items-center gap-3 rounded-md px-3 py-2.5 text-sm text-base-content/75 transition-colors hover:bg-base-300/60 hover:text-base-content';
const navActiveClass = '!bg-base-300 !text-primary';

function logout(): void {
  auth.logout();
  void router.push({ name: 'login' });
}
</script>

<template>
  <div class="vault-bg flex min-h-svh">
    <aside
      class="flex w-56 shrink-0 flex-col border-r border-base-300/80 bg-base-200/50 md:w-60"
      :aria-label="t('nav.menuAriaLabel')"
    >
      <div class="flex h-14 items-center px-5">
        <p class="font-brand text-xl leading-none tracking-tight">envm-vault</p>
      </div>

      <nav class="flex flex-1 flex-col gap-1 px-3 pt-2">
        <RouterLink
          :to="{ name: 'home' }"
          :class="navLinkClass"
          :exact-active-class="navActiveClass"
        >
          <HomeIcon class="size-5 shrink-0" aria-hidden="true" />
          {{ t('nav.home') }}
        </RouterLink>
        <RouterLink
          :to="{ name: 'api-key' }"
          :class="navLinkClass"
          :exact-active-class="navActiveClass"
        >
          <KeyIcon class="size-5 shrink-0" aria-hidden="true" />
          {{ t('nav.apiKey') }}
        </RouterLink>
        <RouterLink
          :to="{ name: 'teams' }"
          :class="navLinkClass"
          :exact-active-class="navActiveClass"
        >
          <UserGroupIcon class="size-5 shrink-0" aria-hidden="true" />
          {{ t('nav.teams') }}
        </RouterLink>
        <RouterLink
          :to="{ name: 'tags' }"
          :class="navLinkClass"
          :exact-active-class="navActiveClass"
        >
          <TagIcon class="size-5 shrink-0" aria-hidden="true" />
          {{ t('nav.tags') }}
        </RouterLink>
        <RouterLink
          :to="{ name: 'environments' }"
          :class="navLinkClass"
          :exact-active-class="navActiveClass"
        >
          <CubeTransparentIcon class="size-5 shrink-0" aria-hidden="true" />
          {{ t('nav.environments') }}
        </RouterLink>
        <RouterLink
          v-if="auth.isAdmin"
          :to="{ name: 'users' }"
          :class="navLinkClass"
          :exact-active-class="navActiveClass"
        >
          <UsersIcon class="size-5 shrink-0" aria-hidden="true" />
          {{ t('nav.users') }}
        </RouterLink>
        <RouterLink
          v-if="auth.isAdmin"
          :to="{ name: 'audit-logs' }"
          :class="navLinkClass"
          :exact-active-class="navActiveClass"
        >
          <ClipboardDocumentListIcon class="size-5 shrink-0" aria-hidden="true" />
          {{ t('nav.audit') }}
        </RouterLink>
      </nav>
    </aside>

    <div class="flex min-w-0 flex-1 flex-col">
      <header
        class="flex h-14 shrink-0 items-center justify-end gap-3 border-b border-base-300/80 bg-base-200/40 px-6 md:gap-4 md:px-8"
      >
        <p v-if="auth.user" class="truncate text-sm text-base-content/70">
          {{ auth.user.email }}
        </p>
        <LocaleSwitcher />
        <button
          class="btn btn-ghost btn-sm gap-2"
          type="button"
          data-testid="logout"
          @click="logout"
        >
          <ArrowRightStartOnRectangleIcon class="size-4" aria-hidden="true" />
          {{ t('common.logout') }}
        </button>
      </header>

      <main class="min-h-0 flex-1 overflow-auto">
        <RouterView />
      </main>
    </div>

    <ConfirmDialog />
  </div>
</template>
