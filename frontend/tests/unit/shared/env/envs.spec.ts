import { describe, expect, it } from 'vitest';
import { assertEnvsValid, createEnvs } from '@/shared/env/envs';

describe('envs', () => {
  it('accepts empty api.baseUrl for Vite proxy', () => {
    expect(() => createEnvs({ api: { baseUrl: '' } })).not.toThrow();
  });

  it('accepts absolute http(s) api.baseUrl', () => {
    expect(() => createEnvs({ api: { baseUrl: 'https://api.example.com' } })).not.toThrow();
  });

  it('rejects relative or invalid api.baseUrl', () => {
    expect(() => assertEnvsValid({ api: { baseUrl: '/v1' } })).toThrow(/VITE_API_BASE_URL/);
  });
});
