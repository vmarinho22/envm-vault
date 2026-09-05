import { afterEach, describe, expect, it, vi } from 'vitest';
import { tokenService } from '@/modules/token/services/TokenService';

function jsonResponse(status: number, body?: unknown): Response {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: async () => body,
  } as Response;
}

describe('TokenService', () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('getMine requests /v1/tokens/me', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(200, { active: false }));
    vi.stubGlobal('fetch', fetchMock);

    await expect(tokenService.getMine()).resolves.toEqual({ active: false });
    expect(fetchMock).toHaveBeenCalledWith('/v1/tokens/me', expect.any(Object));
  });

  it('create posts to /v1/tokens', async () => {
    const body = { token: 'ev_abc', prefix: 'ev_abc1234', created_at: '2026-01-01T00:00:00Z' };
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(201, body));
    vi.stubGlobal('fetch', fetchMock);

    await expect(tokenService.create()).resolves.toEqual(body);
    expect(fetchMock).toHaveBeenCalledWith(
      '/v1/tokens',
      expect.objectContaining({ method: 'POST' }),
    );
  });

  it('revokeOwn deletes /v1/tokens', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(204));
    vi.stubGlobal('fetch', fetchMock);

    await tokenService.revokeOwn();
    expect(fetchMock).toHaveBeenCalledWith(
      '/v1/tokens',
      expect.objectContaining({ method: 'DELETE' }),
    );
  });

  it('revokeForUser deletes /v1/users/:id/token', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(204));
    vi.stubGlobal('fetch', fetchMock);

    await tokenService.revokeForUser('user-1');
    expect(fetchMock).toHaveBeenCalledWith(
      '/v1/users/user-1/token',
      expect.objectContaining({ method: 'DELETE' }),
    );
  });
});
