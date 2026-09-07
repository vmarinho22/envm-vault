import { afterEach, describe, expect, it, vi } from 'vitest';
import { auditService } from '@/modules/audit/services/AuditService';

function jsonResponse(status: number, body: unknown): Response {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: async () => body,
  } as Response;
}

describe('AuditService', () => {
  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('list gets /v1/audit-logs with repeated user id params', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      jsonResponse(200, { items: [], page: 1, page_size: 25, total: 0 }),
    );
    vi.stubGlobal('fetch', fetchMock);

    await auditService.list({
      page: 2,
      page_size: 25,
      sort: 'asc',
      actor_user_ids: ['actor-1', 'actor-2'],
      project_ids: ['proj-1'],
      occurred_from: '2026-01-01',
    });

    expect(fetchMock).toHaveBeenCalledWith(
      '/v1/audit-logs?page=2&page_size=25&sort=asc&actor_user_id=actor-1&actor_user_id=actor-2&project_id=proj-1&occurred_from=2026-01-01',
      expect.any(Object),
    );
  });
});
