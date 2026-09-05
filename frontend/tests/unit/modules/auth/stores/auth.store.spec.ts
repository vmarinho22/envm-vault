import { createPinia, setActivePinia } from 'pinia';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { useAuthStore } from '@/modules/auth/stores/auth.store';
import { tokenService } from '@/modules/token/services/TokenService';

vi.mock('@/modules/token/services/TokenService', () => ({
  tokenService: {
    getMine: vi.fn(),
    create: vi.fn(),
    revokeOwn: vi.fn(),
    revokeForUser: vi.fn(),
  },
}));

const getMineMock = vi.mocked(tokenService.getMine);

const session = {
  token: 'jwt-token',
  user: { id: '1', email: 'admin@example.com', role: 'user' as const },
};

describe('auth store', () => {
  beforeEach(() => {
    sessionStorage.clear();
    setActivePinia(createPinia());
    getMineMock.mockReset();
  });

  it('ensureHasApiToken fetches once and caches the result', async () => {
    getMineMock.mockResolvedValue({ active: true, prefix: 'ev_abcd1234' });
    const auth = useAuthStore();

    await expect(auth.ensureHasApiToken()).resolves.toBe(true);
    await expect(auth.ensureHasApiToken()).resolves.toBe(true);

    expect(getMineMock).toHaveBeenCalledTimes(1);
    expect(auth.hasApiToken).toBe(true);
  });

  it('ensureHasApiToken treats fetch errors as no token', async () => {
    getMineMock.mockRejectedValue(new Error('network'));
    const auth = useAuthStore();

    await expect(auth.ensureHasApiToken()).resolves.toBe(false);
    expect(auth.hasApiToken).toBe(false);
  });

  it('setSession clears the cached token status', async () => {
    getMineMock.mockResolvedValue({ active: true, prefix: 'ev_abcd1234' });
    const auth = useAuthStore();

    await auth.ensureHasApiToken();
    auth.setSession(session);

    expect(auth.hasApiToken).toBeNull();

    getMineMock.mockResolvedValue({ active: false });
    await expect(auth.ensureHasApiToken()).resolves.toBe(false);
    expect(getMineMock).toHaveBeenCalledTimes(2);
  });
});
