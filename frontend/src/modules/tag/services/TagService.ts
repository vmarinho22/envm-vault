import type { CreateTagInput, ListTagsParams, Tag, UpdateTagInput } from '@/modules/tag/types';
import { httpClient } from '@/shared/http';

export interface TagService {
  list(params?: ListTagsParams): Promise<Tag[]>;
  create(input: CreateTagInput): Promise<Tag>;
  update(id: string, input: UpdateTagInput): Promise<Tag>;
  remove(id: string): Promise<void>;
}

function buildListPath(params?: ListTagsParams): string {
  const query = new URLSearchParams();
  const title = params?.title?.trim();
  if (title) {
    query.set('title', title);
  }
  const qs = query.toString();
  return qs ? `/v1/tags?${qs}` : '/v1/tags';
}

export const tagService: TagService = {
  list(params?: ListTagsParams): Promise<Tag[]> {
    return httpClient.request<Tag[]>(buildListPath(params));
  },

  create(input: CreateTagInput): Promise<Tag> {
    return httpClient.request<Tag>('/v1/tags', {
      method: 'POST',
      body: input,
    });
  },

  update(id: string, input: UpdateTagInput): Promise<Tag> {
    return httpClient.request<Tag>(`/v1/tags/${id}`, {
      method: 'PATCH',
      body: input,
    });
  },

  remove(id: string): Promise<void> {
    return httpClient.request<void>(`/v1/tags/${id}`, {
      method: 'DELETE',
    });
  },
};
