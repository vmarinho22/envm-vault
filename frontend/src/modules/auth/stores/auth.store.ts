import { computed, ref } from 'vue';
import { defineStore } from 'pinia';
import { tokenService } from '@/modules/token/services/TokenService';
import type { User, UserRole } from '@/modules/user/types';

const STORAGE_KEY = 'envm-vault.auth';

const USER_ROLES: readonly UserRole[] = ['super_admin', 'admin', 'user'];

interface AuthSession {
  token: string;
  user: User;
}

function isPlainObject(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value);
}

function isUserRole(value: string): value is UserRole {
  return (USER_ROLES as readonly string[]).includes(value);
}

function parseStoredUser(value: unknown): User | undefined {
  if (!isPlainObject(value)) return;

  const { id, email, role } = value;
  if (typeof id !== 'string' || typeof email !== 'string' || typeof role !== 'string') return;
  if (!isUserRole(role)) return;

  return { id, email, role };
}

function parseStoredSession(value: unknown): AuthSession | undefined {
  if (!isPlainObject(value)) return;

  const { token, user: storedUser } = value;
  if (typeof token !== 'string') return;

  const user = parseStoredUser(storedUser);
  if (!user) return;

  return { token, user };
}

function loadSessionFromStorage(): AuthSession | undefined {
  const storedJson = sessionStorage.getItem(STORAGE_KEY);
  if (!storedJson) return;

  try {
    const sessionData: unknown = JSON.parse(storedJson);
    return parseStoredSession(sessionData);
  } catch {
    sessionStorage.removeItem(STORAGE_KEY);
  }
}

export const useAuthStore = defineStore('auth', () => {
  const restoredSession = loadSessionFromStorage();
  const token = ref<string | null>(restoredSession?.token ?? null);
  const user = ref<User | null>(restoredSession?.user ?? null);
  const hasApiToken = ref<boolean | null>(null);

  const isAuthenticated = computed(() => token.value !== null && user.value !== null);
  const isAdmin = computed(() => user.value?.role === 'admin' || user.value?.role === 'super_admin');
  const isSuperAdmin = computed(() => user.value?.role === 'super_admin');

  function setSession(session: AuthSession): void {
    token.value = session.token;
    user.value = session.user;
    hasApiToken.value = null;
    sessionStorage.setItem(STORAGE_KEY, JSON.stringify(session));
  }

  function setHasApiToken(value: boolean): void {
    hasApiToken.value = value;
  }

  async function ensureHasApiToken(): Promise<boolean> {
    if (hasApiToken.value !== null) {
      return hasApiToken.value;
    }

    try {
      const status = await tokenService.getMine();
      hasApiToken.value = status.active;
    } catch {
      hasApiToken.value = false;
    }

    return hasApiToken.value;
  }

  function logout(): void {
    token.value = null;
    user.value = null;
    hasApiToken.value = null;
    sessionStorage.removeItem(STORAGE_KEY);
  }

  return {
    token,
    user,
    hasApiToken,
    isAuthenticated,
    isAdmin,
    isSuperAdmin,
    setSession,
    setHasApiToken,
    ensureHasApiToken,
    logout,
  };
});
