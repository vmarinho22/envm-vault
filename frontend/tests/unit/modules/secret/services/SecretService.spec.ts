import { beforeEach, describe, expect, it, vi } from 'vitest';
import { formatDotEnv, secretService } from '@/modules/secret/services/SecretService';
import { httpClient } from '@/shared/http';

vi.mock('@/shared/http', () => ({
  httpClient: {
    request: vi.fn(),
  },
}));

const requestMock = vi.mocked(httpClient.request);

describe('formatDotEnv', () => {
  it('formats simple and quoted values', () => {
    expect(
      formatDotEnv([
        { name: 'A', value: 'simple' },
        { name: 'B', value: 'has space' },
        { name: 'C', value: '' },
      ]),
    ).toBe('A=simple\nB="has space"\nC=""\n');
  });
});

describe('secretService', () => {
  beforeEach(() => {
    requestMock.mockReset();
  });

  it('list hits project secrets path', async () => {
    requestMock.mockResolvedValue({ items: [] });
    await secretService.list('proj-1');
    expect(requestMock).toHaveBeenCalledWith('/v1/projects/proj-1/secrets');
  });

  it('search hits catalog with query params', async () => {
    requestMock.mockResolvedValue({ items: [] });
    await secretService.search({ title: 'DB', excludeProjectId: 'p1', limit: 10 });
    expect(requestMock).toHaveBeenCalledWith(
      '/v1/secrets?title=DB&exclude_project_id=p1&limit=10',
    );
  });

  it('attach posts secret_id', async () => {
    requestMock.mockResolvedValue({
      id: 's1',
      name: 'X',
      values: {},
      project_count: 2,
      created_at: '',
      updated_at: '',
    });
    await secretService.attach('p1', 's1');
    expect(requestMock).toHaveBeenCalledWith('/v1/projects/p1/secrets/attach', {
      method: 'POST',
      body: { secret_id: 's1' },
    });
  });

  it('setValue puts plaintext body', async () => {
    requestMock.mockResolvedValue({
      id: 's1',
      name: 'X',
      values: {},
      project_count: 1,
      created_at: '',
      updated_at: '',
    });
    await secretService.setValue('p1', 's1', 'e1', 'secret');
    expect(requestMock).toHaveBeenCalledWith('/v1/projects/p1/secrets/s1/values/e1', {
      method: 'PUT',
      body: { value: 'secret' },
    });
  });
});
