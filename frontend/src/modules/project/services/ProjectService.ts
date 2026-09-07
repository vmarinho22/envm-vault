import type {
  CreateProjectInput,
  ListProjectsParams,
  Project,
  ProjectListResponse,
  UpdateProjectInput,
} from '@/modules/project/types';
import { httpClient } from '@/shared/http';

export interface ProjectService {
  list(params?: ListProjectsParams): Promise<ProjectListResponse>;
  getBySlug(slug: string): Promise<Project>;
  create(input: CreateProjectInput): Promise<Project>;
  update(id: string, input: UpdateProjectInput): Promise<Project>;
  remove(id: string): Promise<void>;
}

function buildListPath(params?: ListProjectsParams): string {
  const query = new URLSearchParams();
  const title = params?.title?.trim();
  if (title) query.set('title', title);
  if (params?.team_ids?.length) query.set('team_ids', params.team_ids.join(','));
  if (params?.tag_ids?.length) query.set('tag_ids', params.tag_ids.join(','));
  if (params?.page != null) query.set('page', String(params.page));
  if (params?.page_size != null) query.set('page_size', String(params.page_size));
  const qs = query.toString();
  return qs ? `/v1/projects?${qs}` : '/v1/projects';
}

export const projectService: ProjectService = {
  list(params?: ListProjectsParams): Promise<ProjectListResponse> {
    return httpClient.request<ProjectListResponse>(buildListPath(params));
  },

  getBySlug(slug: string): Promise<Project> {
    return httpClient.request<Project>(`/v1/projects/by-slug/${encodeURIComponent(slug)}`);
  },

  create(input: CreateProjectInput): Promise<Project> {
    return httpClient.request<Project>('/v1/projects', {
      method: 'POST',
      body: input,
    });
  },

  update(id: string, input: UpdateProjectInput): Promise<Project> {
    return httpClient.request<Project>(`/v1/projects/${id}`, {
      method: 'PATCH',
      body: input,
    });
  },

  remove(id: string): Promise<void> {
    return httpClient.request<void>(`/v1/projects/${id}`, {
      method: 'DELETE',
    });
  },
};
