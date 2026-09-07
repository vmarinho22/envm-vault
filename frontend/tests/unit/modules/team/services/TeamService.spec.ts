import { afterEach, describe, expect, it, vi } from 'vitest';
import { teamService } from '@/modules/team/services/TeamService';

function jsonResponse(status: number, body?: unknown): Response {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: async () => body,
  } as Response;
}

describe('TeamService', () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('list gets /v1/teams without query when name is empty', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(200, []));
    vi.stubGlobal('fetch', fetchMock);

    await teamService.list();
    expect(fetchMock).toHaveBeenCalledWith('/v1/teams', expect.any(Object));
  });

  it('list gets /v1/teams?name= when searching', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(200, []));
    vi.stubGlobal('fetch', fetchMock);

    await teamService.list({ name: 'plat' });
    expect(fetchMock).toHaveBeenCalledWith('/v1/teams?name=plat', expect.any(Object));
  });

  it('create posts to /v1/teams', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      jsonResponse(201, {
        id: '1',
        name: 'Platform',
        slug: 'platform',
        accountable: { id: 'u1', email: 'a@x.com' },
        created_at: '2026-01-01T00:00:00Z',
        updated_at: '2026-01-01T00:00:00Z',
      }),
    );
    vi.stubGlobal('fetch', fetchMock);

    await teamService.create({ name: 'Platform', accountable_id: 'u1' });
    expect(fetchMock).toHaveBeenCalledWith(
      '/v1/teams',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ name: 'Platform', accountable_id: 'u1' }),
      }),
    );
  });

  it('update patches /v1/teams/:id', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      jsonResponse(200, {
        id: '1',
        name: 'Platform',
        slug: 'platform',
        accountable: { id: 'u1', email: 'a@x.com' },
        created_at: '2026-01-01T00:00:00Z',
        updated_at: '2026-01-01T00:00:00Z',
      }),
    );
    vi.stubGlobal('fetch', fetchMock);

    await teamService.update('1', { name: 'Platform' });
    expect(fetchMock).toHaveBeenCalledWith(
      '/v1/teams/1',
      expect.objectContaining({ method: 'PATCH' }),
    );
  });

  it('remove deletes /v1/teams/:id', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(204));
    vi.stubGlobal('fetch', fetchMock);

    await teamService.remove('1');
    expect(fetchMock).toHaveBeenCalledWith(
      '/v1/teams/1',
      expect.objectContaining({ method: 'DELETE' }),
    );
  });
});
