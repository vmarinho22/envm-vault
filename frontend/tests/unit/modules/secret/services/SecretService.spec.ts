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

  it('setValue puts plaintext body', async () => {
    requestMock.mockResolvedValue({ id: 's1', name: 'X', values: {}, created_at: '', updated_at: '' });
    await secretService.setValue('p1', 's1', 'e1', 'secret');
    expect(requestMock).toHaveBeenCalledWith('/v1/projects/p1/secrets/s1/values/e1', {
      method: 'PUT',
      body: { value: 'secret' },
    });
  });
});
