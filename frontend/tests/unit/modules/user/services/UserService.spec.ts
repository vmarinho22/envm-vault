import { afterEach, describe, expect, it, vi } from 'vitest';
import { userService } from '@/modules/user/services/UserService';

function jsonResponse(status: number, body: unknown): Response {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: async () => body,
  } as Response;
}

describe('UserService', () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('register posts to /v1/auth/register', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      jsonResponse(201, { id: '1', email: 'user@example.com', role: 'user' }),
    );
    vi.stubGlobal('fetch', fetchMock);

    await userService.register('user@example.com', 'user');

    expect(fetchMock).toHaveBeenCalledWith(
      '/v1/auth/register',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ email: 'user@example.com', role: 'user' }),
      }),
    );
  });

  it('list gets /v1/users', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(200, []));
    vi.stubGlobal('fetch', fetchMock);

    await expect(userService.list()).resolves.toEqual([]);
    expect(fetchMock).toHaveBeenCalledWith('/v1/users', expect.any(Object));
  });
});
