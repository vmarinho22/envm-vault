import type {
  CreateSecretInput,
  DownloadItem,
  DownloadResponse,
  RevealSecretResponse,
  Secret,
  SecretListResponse,
} from '@/modules/secret/types';
import { httpClient } from '@/shared/http';

export interface SecretService {
  list(projectId: string): Promise<SecretListResponse>;
  create(projectId: string, input: CreateSecretInput): Promise<Secret>;
  updateName(projectId: string, secretId: string, name: string): Promise<Secret>;
  remove(projectId: string, secretId: string): Promise<void>;
  setValue(
    projectId: string,
    secretId: string,
    environmentId: string,
    value: string,
  ): Promise<Secret>;
  revealValue(
    projectId: string,
    secretId: string,
    environmentId: string,
  ): Promise<RevealSecretResponse>;
  clearValue(projectId: string, secretId: string, environmentId: string): Promise<void>;
  exportEnvironment(projectId: string, environmentId: string): Promise<DownloadResponse>;
}

export function formatDotEnv(items: DownloadItem[]): string {
  return items
    .map((item) => `${item.name}=${quoteDotEnv(item.value)}`)
    .join('\n')
    .concat(items.length > 0 ? '\n' : '');
}

function quoteDotEnv(value: string): string {
  if (value === '') return '""';
  if (!/[\s"'#]/.test(value) && !value.includes('\n') && !value.includes('\r')) {
    return value;
  }
  return `"${value.replace(/\\/g, '\\\\').replace(/"/g, '\\"').replace(/\n/g, '\\n').replace(/\r/g, '\\r')}"`;
}

export const secretService: SecretService = {
  list(projectId: string): Promise<SecretListResponse> {
    return httpClient.request<SecretListResponse>(`/v1/projects/${projectId}/secrets`);
  },

  create(projectId: string, input: CreateSecretInput): Promise<Secret> {
    return httpClient.request<Secret>(`/v1/projects/${projectId}/secrets`, {
      method: 'POST',
      body: input,
    });
  },

  updateName(projectId: string, secretId: string, name: string): Promise<Secret> {
    return httpClient.request<Secret>(`/v1/projects/${projectId}/secrets/${secretId}`, {
      method: 'PATCH',
      body: { name },
    });
  },

  remove(projectId: string, secretId: string): Promise<void> {
    return httpClient.request<void>(`/v1/projects/${projectId}/secrets/${secretId}`, {
      method: 'DELETE',
    });
  },

  setValue(
    projectId: string,
    secretId: string,
    environmentId: string,
    value: string,
  ): Promise<Secret> {
    return httpClient.request<Secret>(
      `/v1/projects/${projectId}/secrets/${secretId}/values/${environmentId}`,
      { method: 'PUT', body: { value } },
    );
  },

  revealValue(
    projectId: string,
    secretId: string,
    environmentId: string,
  ): Promise<RevealSecretResponse> {
    return httpClient.request<RevealSecretResponse>(
      `/v1/projects/${projectId}/secrets/${secretId}/values/${environmentId}`,
    );
  },

  clearValue(projectId: string, secretId: string, environmentId: string): Promise<void> {
    return httpClient.request<void>(
      `/v1/projects/${projectId}/secrets/${secretId}/values/${environmentId}`,
      { method: 'DELETE' },
    );
  },

  exportEnvironment(projectId: string, environmentId: string): Promise<DownloadResponse> {
    return httpClient.request<DownloadResponse>(
      `/v1/projects/${projectId}/environments/${environmentId}/values`,
    );
  },
};
