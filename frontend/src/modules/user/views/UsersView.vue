<script setup lang="ts">
import { TrashIcon } from '@heroicons/vue/24/outline';
import { onMounted } from 'vue';
import { useUsers } from '@/modules/user/composables/useUsers';
import { roleLabel } from '@/modules/user/roleLabel';
import type { UserWithToken } from '@/modules/user/types';

const { users, loading, mutatingId, error, fetchUsers, revokeToken } = useUsers();

function tokenStatus(user: UserWithToken): string {
  if (!user.token) return 'Sem chave';
  if (user.token.revoked_at) return 'Revogada';
  return user.token.prefix;
}

function canRevoke(user: UserWithToken): boolean {
  return user.token != null && user.token.revoked_at == null;
}

onMounted(() => {
  void fetchUsers();
});
</script>

<template>
  <section class="px-8 py-10 md:px-12">
    <h1 class="text-2xl font-medium tracking-tight md:text-3xl">Usuários</h1>
    <p class="mt-3 max-w-lg text-base text-base-content/70">
      Revogue a chave de API de um usuário. Ele precisará gerar uma nova para continuar.
    </p>

    <p v-if="error" class="mt-6 text-sm text-error" role="alert">{{ error }}</p>
    <p v-else-if="loading" class="mt-8 text-base-content/70">Carregando...</p>

    <div v-else class="mt-8 overflow-x-auto">
      <table class="table">
        <thead>
          <tr>
            <th>E-mail</th>
            <th>Papel</th>
            <th>Chave</th>
            <th>
              <span class="sr-only">Ações</span>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="users.length === 0">
            <td colspan="4" class="text-base-content/70">Nenhum usuário encontrado.</td>
          </tr>
          <tr v-for="user in users" :key="user.id">
            <td>{{ user.email }}</td>
            <td>{{ roleLabel(user.role) }}</td>
            <td class="font-mono text-sm">{{ tokenStatus(user) }}</td>
            <td class="text-right">
              <button
                class="btn btn-error btn-outline btn-sm gap-2"
                type="button"
                :disabled="!canRevoke(user) || mutatingId === user.id"
                @click="revokeToken(user.id)"
              >
                <span v-if="mutatingId === user.id" class="loading loading-spinner loading-xs" />
                <TrashIcon v-else class="size-4" aria-hidden="true" />
                Revogar
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </section>
</template>
