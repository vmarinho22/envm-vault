import { envs } from '@/shared/env/envs';
import { ApiError } from './ApiError';

export type TokenGetter = () => string | null;

export type RequestOptions = Omit<RequestInit, 'body'> & {
  body?: unknown;
};

export interface HttpClient {
  configure(getter: TokenGetter): void;
  request<T>(path: string, options?: RequestOptions): Promise<T>;
}

const DEFAULT_ERROR_MESSAGE = 'request failed';

let getAccessToken: TokenGetter = () => null;

function apiBaseUrl(): string {
  return envs.api.baseUrl;
}

function buildRequestHeaders(options: RequestOptions): Headers {
  const headers = new Headers(options.headers);

  if (options.body !== undefined && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }

  const accessToken = getAccessToken();
  if (accessToken && !headers.has('Authorization')) {
    headers.set('Authorization', `Bearer ${accessToken}`);
  }

  return headers;
}

function serializeBody(body: unknown): string | undefined {
  if (body === undefined) {
    return undefined;
  }

  return JSON.stringify(body);
}

function readErrorMessage(payload: unknown): string | null {
  if (typeof payload !== 'object' || payload === null || !('error' in payload)) {
    return null;
  }

  const { error } = payload;
  if (typeof error !== 'string' || error.length === 0) {
    return null;
  }

  return error;
}

async function throwApiError(response: Response): Promise<never> {
  let message = DEFAULT_ERROR_MESSAGE;

  try {
    const payload: unknown = await response.json();
    message = readErrorMessage(payload) ?? message;
  } catch {
    // Keep the default message when the error body is not JSON.
  }

  throw new ApiError(response.status, message);
}

async function parseSuccessBody<T>(response: Response): Promise<T> {
  if (response.status === 204) {
    return undefined as T;
  }

  return (await response.json()) as T;
}

export const httpClient: HttpClient = {
  configure(getter: TokenGetter): void {
    getAccessToken = getter;
  },

  async request<T>(path: string, options: RequestOptions = {}): Promise<T> {
    const { body, ...fetchOptions } = options;

    const response = await fetch(`${apiBaseUrl()}${path}`, {
      ...fetchOptions,
      headers: buildRequestHeaders(options),
      body: serializeBody(body),
    });

    if (!response.ok) {
      return throwApiError(response);
    }

    return parseSuccessBody<T>(response);
  },
};
