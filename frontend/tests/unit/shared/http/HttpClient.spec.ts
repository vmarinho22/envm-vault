import { afterEach, describe, expect, it, vi } from 'vitest';
import { ApiError } from '@/shared/http/ApiError';
import { httpClient } from '@/shared/http';

function jsonResponse(status: number, body: unknown): Response {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: async () => body,
  } as Response;
}

describe('HttpClient', () => {
  afterEach(() => {
    vi.unstubAllGlobals();
    httpClient.configure(() => null);
  });

  it('returns parsed JSON on success', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(200, { message: 'ok' }));
    vi.stubGlobal('fetch', fetchMock);

    const result = await httpClient.request<{ message: string }>('/v1/health');

    expect(result).toEqual({ message: 'ok' });
    expect(fetchMock).toHaveBeenCalledWith(
      '/v1/health',
      expect.objectContaining({
        headers: expect.any(Headers),
      }),
    );
  });

  it('throws ApiError with API message on 4xx', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue(jsonResponse(400, { error: 'invalid email' })),
    );

    const err = await httpClient
      .request('/v1/auth/login', {
        method: 'POST',
        body: { email: 'x' },
      })
      .catch((e: unknown) => e);

    expect(err).toBeInstanceOf(ApiError);
    expect(err).toMatchObject({ status: 400, message: 'invalid email' });
  });

  it('injects Bearer token from configure', async () => {
    const fetchMock = vi.fn().mockResolvedValue(jsonResponse(200, { id: '1' }));
    vi.stubGlobal('fetch', fetchMock);
    httpClient.configure(() => 'test-token');

    await httpClient.request('/v1/auth/register', { method: 'POST', body: { email: 'a@b.com' } });

    const headers = fetchMock.mock.calls[0]?.[1]?.headers as Headers;
    expect(headers.get('Authorization')).toBe('Bearer test-token');
  });
});
