export interface ApiEnv {
  /** Absolute API URL, or empty string to use Vite proxy in development. */
  baseUrl: string;
}

export interface Envs {
  api: ApiEnv;
}

function readApiBaseUrl(): string {
  return (import.meta.env.VITE_API_BASE_URL ?? '').trim();
}

function isAbsoluteHttpUrl(value: string): boolean {
  try {
    const url = new URL(value);
    return url.protocol === 'http:' || url.protocol === 'https:';
  } catch {
    return false;
  }
}

function loadEnvs(): Envs {
  return {
    api: {
      baseUrl: readApiBaseUrl(),
    },
  };
}

function assertEnvs(values: Envs): void {
  if (values.api.baseUrl !== '' && !isAbsoluteHttpUrl(values.api.baseUrl)) {
    throw new Error('VITE_API_BASE_URL must be an absolute http(s) URL when set');
  }
}

/** Central public env — same idea as `import { envs } from 'env'`. */
export const envs: Envs = loadEnvs();

/** Throws if public env is malformed (empty api.baseUrl is allowed). */
export function assertEnvsValid(values: Envs = envs): void {
  assertEnvs(values);
}

/** Test helper: validate explicit values without reading import.meta.env. */
export function createEnvs(values: Envs): Envs {
  assertEnvs(values);
  return values;
}
