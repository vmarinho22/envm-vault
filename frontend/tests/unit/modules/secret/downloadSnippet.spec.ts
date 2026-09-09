import { describe, expect, it } from 'vitest';
import { buildDownloadSnippet, buildDownloadUrl } from '@/modules/secret/downloadSnippet';

describe('downloadSnippet', () => {
  it('builds download url against the API origin, not the Vite UI', () => {
    expect(buildDownloadUrl('envm-api', 'production')).toBe(
      'http://localhost:8080/v1/projects/envm-api/env?environment=production&format=env',
    );
  });

  it('builds curl for linux/macos and PowerShell for windows', () => {
    const linux = buildDownloadSnippet('linux', 'envm-api', 'dev');
    expect(linux).toContain('curl -fsSL');
    expect(linux).toContain('Authorization: Bearer ${ENVM_TOKEN}');
    expect(linux).toContain(
      'http://localhost:8080/v1/projects/envm-api/env?environment=dev&format=env',
    );
    expect(linux).toContain('-o .env');

    const macos = buildDownloadSnippet('macos', 'envm-api', 'dev');
    expect(macos).toContain('curl -fsSL');

    const win = buildDownloadSnippet('windows', 'envm-api', 'staging');
    expect(win).toContain('Invoke-WebRequest');
    expect(win).toContain('$env:ENVM_TOKEN');
    expect(win).toContain('-OutFile .env');
    expect(win).toContain(
      'http://localhost:8080/v1/projects/envm-api/env?environment=staging&format=env',
    );
  });
});
