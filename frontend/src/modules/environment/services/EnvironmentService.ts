import type {
  CreateEnvironmentInput,
  Environment,
  ListEnvironmentsParams,
  UpdateEnvironmentInput,
} from '@/modules/environment/types';
import { httpClient } from '@/shared/http';

export interface EnvironmentService {
  list(params?: ListEnvironmentsParams): Promise<Environment[]>;
  create(input: CreateEnvironmentInput): Promise<Environment>;
  update(id: string, input: UpdateEnvironmentInput): Promise<Environment>;
  remove(id: string): Promise<void>;
}

function buildListPath(params?: ListEnvironmentsParams): string {
  const query = new URLSearchParams();
  const title = params?.title?.trim();
  if (title) {
    query.set('title', title);
  }
  const qs = query.toString();
  return qs ? `/v1/environments?${qs}` : '/v1/environments';
}

export const environmentService: EnvironmentService = {
  list(params?: ListEnvironmentsParams): Promise<Environment[]> {
    return httpClient.request<Environment[]>(buildListPath(params));
  },

  create(input: CreateEnvironmentInput): Promise<Environment> {
    return httpClient.request<Environment>('/v1/environments', {
      method: 'POST',
      body: input,
    });
  },

  update(id: string, input: UpdateEnvironmentInput): Promise<Environment> {
    return httpClient.request<Environment>(`/v1/environments/${id}`, {
      method: 'PATCH',
      body: input,
    });
  },

  remove(id: string): Promise<void> {
    return httpClient.request<void>(`/v1/environments/${id}`, {
      method: 'DELETE',
    });
  },
};
