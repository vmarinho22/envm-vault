import { afterEach, describe, expect, it, vi } from 'vitest';
import { tagService } from '@/modules/tag/services/TagService';

function jsonResponse(status: number, body?: unknown): Response {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: async () => body,
  } as Response;
}

describe('TagService', () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('list gets /v1/tags without query when title is empty', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(200, []));
    vi.stubGlobal('fetch', fetchMock);

    await tagService.list();
    expect(fetchMock).toHaveBeenCalledWith('/v1/tags', expect.any(Object));
  });

  it('list gets /v1/tags?title= when searching', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(200, []));
    vi.stubGlobal('fetch', fetchMock);

    await tagService.list({ title: 'back' });
    expect(fetchMock).toHaveBeenCalledWith('/v1/tags?title=back', expect.any(Object));
  });

  it('create posts to /v1/tags', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      jsonResponse(201, {
        id: '1',
        title: 'Backend',
        color: '#6366F1',
        created_at: '2026-01-01T00:00:00Z',
        updated_at: '2026-01-01T00:00:00Z',
      }),
    );
    vi.stubGlobal('fetch', fetchMock);

    await tagService.create({ title: 'Backend', color: '#6366F1' });
    expect(fetchMock).toHaveBeenCalledWith(
      '/v1/tags',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ title: 'Backend', color: '#6366F1' }),
      }),
    );
  });

  it('update patches /v1/tags/:id', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      jsonResponse(200, {
        id: '1',
        title: 'Backend',
        color: '#6366F1',
        created_at: '2026-01-01T00:00:00Z',
        updated_at: '2026-01-01T00:00:00Z',
      }),
    );
    vi.stubGlobal('fetch', fetchMock);

    await tagService.update('1', { title: 'Backend' });
    expect(fetchMock).toHaveBeenCalledWith(
      '/v1/tags/1',
      expect.objectContaining({ method: 'PATCH' }),
    );
  });

  it('remove deletes /v1/tags/:id', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(204));
    vi.stubGlobal('fetch', fetchMock);

    await tagService.remove('1');
    expect(fetchMock).toHaveBeenCalledWith(
      '/v1/tags/1',
      expect.objectContaining({ method: 'DELETE' }),
    );
  });
});
