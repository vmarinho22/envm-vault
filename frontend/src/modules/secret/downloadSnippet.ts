import { envs } from '@/shared/env/envs';

/** Absolute API origin for shell snippets (Vite proxy origin when VITE_API_BASE_URL is empty). */
export function publicApiBaseUrl(): string {
  const configured = envs.api.baseUrl.replace(/\/$/, '');
  if (configured) return configured;
  if (typeof window !== 'undefined' && window.location?.origin) {
    return window.location.origin;
  }
  return 'http://localhost:8080';
}

export type DownloadOs = 'linux' | 'macos' | 'windows';

export function buildDownloadUrl(projectSlug: string, envTitle: string): string {
  const base = publicApiBaseUrl();
  const slug = encodeURIComponent(projectSlug);
  const params = new URLSearchParams({
    environment: envTitle,
    format: 'env',
  });
  return `${base}/v1/projects/${slug}/env?${params.toString()}`;
}

export function buildDownloadSnippet(
  os: DownloadOs,
  projectSlug: string,
  envTitle: string,
): string {
  const url = buildDownloadUrl(projectSlug, envTitle);

  if (os === 'windows') {
    return [
      `# Set your personal API key once in the session:`,
      `#   $env:ENVM_TOKEN = "ev_…"`,
      `Invoke-WebRequest \``,
      `  -Headers @{ Authorization = "Bearer $env:ENVM_TOKEN" } \``,
      `  -Uri "${url}" \``,
      `  -OutFile .env`,
    ].join('\n');
  }

  const shellHint =
    os === 'macos'
      ? '# Set your personal API key once in the session:\n#   export ENVM_TOKEN="ev_…"'
      : '# Set your personal API key once in the session:\n#   export ENVM_TOKEN="ev_…"';

  return [
    shellHint,
    `curl -fsSL \\`,
    `  -H "Authorization: Bearer \${ENVM_TOKEN}" \\`,
    `  "${url}" \\`,
    `  -o .env`,
  ].join('\n');
}
