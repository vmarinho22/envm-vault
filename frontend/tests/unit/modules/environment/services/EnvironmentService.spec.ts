import { afterEach, describe, expect, it, vi } from 'vitest';
import { environmentService } from '@/modules/environment/services/EnvironmentService';

function jsonResponse(status: number, body?: unknown): Response {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: async () => body,
  } as Response;
}

describe('EnvironmentService', () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('list gets /v1/environments without query when title is empty', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(200, []));
    vi.stubGlobal('fetch', fetchMock);

    await environmentService.list();
    expect(fetchMock).toHaveBeenCalledWith('/v1/environments', expect.any(Object));
  });

  it('list gets /v1/environments?title= when searching', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(200, []));
    vi.stubGlobal('fetch', fetchMock);

    await environmentService.list({ title: 'stag' });
    expect(fetchMock).toHaveBeenCalledWith('/v1/environments?title=stag', expect.any(Object));
  });

  it('create posts to /v1/environments', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      jsonResponse(201, {
        id: '1',
        title: 'qa',
        created_at: '2026-01-01T00:00:00Z',
        updated_at: '2026-01-01T00:00:00Z',
      }),
    );
    vi.stubGlobal('fetch', fetchMock);

    await environmentService.create({ title: 'qa' });
    expect(fetchMock).toHaveBeenCalledWith(
      '/v1/environments',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ title: 'qa' }),
      }),
    );
  });

  it('update patches /v1/environments/:id', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      jsonResponse(200, {
        id: '1',
        title: 'qa',
        created_at: '2026-01-01T00:00:00Z',
        updated_at: '2026-01-01T00:00:00Z',
      }),
    );
    vi.stubGlobal('fetch', fetchMock);

    await environmentService.update('1', { title: 'qa' });
    expect(fetchMock).toHaveBeenCalledWith(
      '/v1/environments/1',
      expect.objectContaining({ method: 'PATCH' }),
    );
  });

  it('remove deletes /v1/environments/:id', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(204));
    vi.stubGlobal('fetch', fetchMock);

    await environmentService.remove('1');
    expect(fetchMock).toHaveBeenCalledWith(
      '/v1/environments/1',
      expect.objectContaining({ method: 'DELETE' }),
    );
  });
});
