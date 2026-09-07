<script setup lang="ts">
import { KeyIcon, LockClosedIcon, LockOpenIcon, PlusIcon } from '@heroicons/vue/24/outline';
import { onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useAuthStore } from '@/modules/auth/stores/auth.store';
import { canManageTarget } from '@/modules/user/canManageTarget';
import UserFormModal from '@/modules/user/components/UserFormModal.vue';
import { useUsers } from '@/modules/user/composables/useUsers';
import { roleLabel } from '@/modules/user/roleLabel';
import type { UserWithToken } from '@/modules/user/types';

const { t } = useI18n();
const auth = useAuthStore();
const {
  users,
  loading,
  mutatingId,
  creating,
  error,
  formOpen,
  formEmail,
  formRole,
  fetchUsers,
  openCreate,
  closeForm,
  submitCreate,
  revokeToken,
  blockUser,
  unblockUser,
} = useUsers();

function tokenStatus(user: UserWithToken): string {
  if (!user.token) return t('users.noKey');
  if (user.token.revoked_at) return t('users.revoked');
  return user.token.prefix;
}

function isBlocked(user: UserWithToken): boolean {
  return user.blocked_at != null && user.blocked_at !== '';
}

function canRevoke(user: UserWithToken): boolean {
  return user.token != null && user.token.revoked_at == null;
}

function canToggleBlock(user: UserWithToken): boolean {
  const actor = auth.user;
  if (!actor) return false;
  if (actor.id === user.id) return false;
  return canManageTarget(actor.role, user.role);
}

onMounted(() => {
  void fetchUsers();
});
</script>

<template>
  <section class="px-8 py-10 md:px-12">
    <div class="flex flex-wrap items-start justify-between gap-4">
      <div>
        <h1 class="text-2xl font-medium tracking-tight md:text-3xl">{{ t('users.title') }}</h1>
        <p class="mt-3 max-w-lg text-base text-base-content/70">
          {{ t('users.description') }}
        </p>
      </div>
      <button class="btn btn-primary gap-2" type="button" @click="openCreate">
        <PlusIcon class="size-4" aria-hidden="true" />
        {{ t('users.create') }}
      </button>
    </div>

    <p v-if="error" class="mt-6 text-sm text-error" role="alert">{{ error }}</p>
    <p v-else-if="loading" class="mt-8 text-base-content/70">{{ t('common.loading') }}</p>

    <div v-else class="mt-8 overflow-x-auto">
      <table class="table">
        <thead>
          <tr>
            <th>{{ t('users.columns.email') }}</th>
            <th>{{ t('users.columns.role') }}</th>
            <th>{{ t('users.columns.status') }}</th>
            <th>{{ t('users.columns.key') }}</th>
            <th>
              <span class="sr-only">{{ t('users.columns.actions') }}</span>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="users.length === 0">
            <td colspan="5" class="text-base-content/70">{{ t('users.empty') }}</td>
          </tr>
          <tr v-for="user in users" :key="user.id">
            <td>{{ user.email }}</td>
            <td>{{ roleLabel(user.role, t) }}</td>
            <td>
              <span
                class="badge badge-sm"
                :class="isBlocked(user) ? 'badge-error' : 'badge-success'"
              >
                {{ isBlocked(user) ? t('users.status.blocked') : t('users.status.active') }}
              </span>
            </td>
            <td class="font-mono text-sm">{{ tokenStatus(user) }}</td>
            <td class="text-right">
              <div class="flex justify-end gap-1">
                <button
                  v-if="canToggleBlock(user) && !isBlocked(user)"
                  class="btn btn-ghost btn-square btn-sm"
                  type="button"
                  :disabled="mutatingId === user.id"
                  :aria-label="t('users.block')"
                  @click="blockUser(user)"
                >
                  <LockClosedIcon class="size-4" aria-hidden="true" />
                </button>
                <button
                  v-else-if="canToggleBlock(user) && isBlocked(user)"
                  class="btn btn-ghost btn-square btn-sm"
                  type="button"
                  :disabled="mutatingId === user.id"
                  :aria-label="t('users.unblock')"
                  @click="unblockUser(user)"
                >
                  <LockOpenIcon class="size-4" aria-hidden="true" />
                </button>
                <button
                  class="btn btn-error btn-outline btn-sm gap-2"
                  type="button"
                  :disabled="!canRevoke(user) || mutatingId === user.id"
                  @click="revokeToken(user.id)"
                >
                  <span v-if="mutatingId === user.id" class="loading loading-spinner loading-xs" />
                  <KeyIcon v-else class="size-4" aria-hidden="true" />
                  {{ t('users.revoke') }}
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <UserFormModal
      :open="formOpen"
      :email="formEmail"
      :role="formRole"
      :saving="creating"
      @close="closeForm"
      @submit="submitCreate"
      @update:email="formEmail = $event"
      @update:role="formRole = $event"
    />
  </section>
</template>
