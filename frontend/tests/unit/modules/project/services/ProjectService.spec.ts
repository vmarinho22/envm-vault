import { afterEach, describe, expect, it, vi } from 'vitest';
import { projectService } from '@/modules/project/services/ProjectService';

function jsonResponse(status: number, body?: unknown): Response {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: async () => body,
  } as Response;
}

const sampleProject = {
  id: 'p1',
  title: 'Payment',
  slug: 'payment',
  teams: [{ id: 't1', name: 'Platform', slug: 'platform' }],
  environments: [{ id: 'e1', title: 'dev' }],
  tags: [{ id: 'g1', title: 'backend', color: '#000000' }],
  created_at: '2026-01-01T00:00:00Z',
  updated_at: '2026-01-01T00:00:00Z',
};

describe('ProjectService', () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('list gets /v1/projects without query when empty', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      jsonResponse(200, { items: [], page: 1, page_size: 12, total: 0 }),
    );
    vi.stubGlobal('fetch', fetchMock);

    await projectService.list();
    expect(fetchMock).toHaveBeenCalledWith('/v1/projects', expect.any(Object));
  });

  it('list builds title, team_ids, tag_ids and pagination query', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      jsonResponse(200, { items: [], page: 2, page_size: 12, total: 0 }),
    );
    vi.stubGlobal('fetch', fetchMock);

    await projectService.list({
      title: 'pay',
      team_ids: ['t1', 't2'],
      tag_ids: ['g1'],
      page: 2,
      page_size: 12,
    });
    expect(fetchMock).toHaveBeenCalledWith(
      '/v1/projects?title=pay&team_ids=t1%2Ct2&tag_ids=g1&page=2&page_size=12',
      expect.any(Object),
    );
  });

  it('getBySlug gets /v1/projects/by-slug/:slug', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(200, sampleProject));
    vi.stubGlobal('fetch', fetchMock);

    await projectService.getBySlug('payment');
    expect(fetchMock).toHaveBeenCalledWith(
      '/v1/projects/by-slug/payment',
      expect.any(Object),
    );
  });

  it('create posts to /v1/projects', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(201, sampleProject));
    vi.stubGlobal('fetch', fetchMock);

    const body = {
      title: 'Payment',
      team_ids: ['t1'],
      environment_ids: ['e1'],
      tag_ids: ['g1'],
    };
    await projectService.create(body);
    expect(fetchMock).toHaveBeenCalledWith(
      '/v1/projects',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify(body),
      }),
    );
  });

  it('update patches /v1/projects/:id', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(200, sampleProject));
    vi.stubGlobal('fetch', fetchMock);

    await projectService.update('p1', { title: 'Payment' });
    expect(fetchMock).toHaveBeenCalledWith(
      '/v1/projects/p1',
      expect.objectContaining({ method: 'PATCH' }),
    );
  });

  it('remove deletes /v1/projects/:id', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(204));
    vi.stubGlobal('fetch', fetchMock);

    await projectService.remove('p1');
    expect(fetchMock).toHaveBeenCalledWith(
      '/v1/projects/p1',
      expect.objectContaining({ method: 'DELETE' }),
    );
  });
});
