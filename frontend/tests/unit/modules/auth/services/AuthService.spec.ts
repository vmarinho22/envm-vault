import { afterEach, describe, expect, it, vi } from 'vitest';
import { authService } from '@/modules/auth/services/AuthService';

function jsonResponse(body: unknown): Response {
  return {
    ok: true,
    status: 200,
    json: async () => body,
  } as Response;
}

describe('AuthService', () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('requestLogin posts to /v1/auth/login', async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValue(jsonResponse({ message: 'verification code sent' }));
    vi.stubGlobal('fetch', fetchMock);

    await authService.requestLogin('admin@example.com');

    expect(fetchMock).toHaveBeenCalledWith(
      '/v1/auth/login',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ email: 'admin@example.com' }),
      }),
    );
  });

  it('verifyLogin posts to /v1/auth/verify', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      jsonResponse({
        token: 'jwt',
        user: { id: '1', email: 'admin@example.com', role: 'super_admin' },
      }),
    );
    vi.stubGlobal('fetch', fetchMock);

    await authService.verifyLogin('admin@example.com', '12345678');

    expect(fetchMock).toHaveBeenCalledWith(
      '/v1/auth/verify',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ email: 'admin@example.com', code: '12345678' }),
      }),
    );
  });
});
