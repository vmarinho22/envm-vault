import type { CreateTeamInput, ListTeamsParams, Team, UpdateTeamInput } from '@/modules/team/types';
import { httpClient } from '@/shared/http';

export interface TeamService {
  list(params?: ListTeamsParams): Promise<Team[]>;
  create(input: CreateTeamInput): Promise<Team>;
  update(id: string, input: UpdateTeamInput): Promise<Team>;
  remove(id: string): Promise<void>;
}

function buildListPath(params?: ListTeamsParams): string {
  const query = new URLSearchParams();
  const name = params?.name?.trim();
  if (name) {
    query.set('name', name);
  }
  const qs = query.toString();
  return qs ? `/v1/teams?${qs}` : '/v1/teams';
}

export const teamService: TeamService = {
  list(params?: ListTeamsParams): Promise<Team[]> {
    return httpClient.request<Team[]>(buildListPath(params));
  },

  create(input: CreateTeamInput): Promise<Team> {
    return httpClient.request<Team>('/v1/teams', {
      method: 'POST',
      body: input,
    });
  },

  update(id: string, input: UpdateTeamInput): Promise<Team> {
    return httpClient.request<Team>(`/v1/teams/${id}`, {
      method: 'PATCH',
      body: input,
    });
  },

  remove(id: string): Promise<void> {
    return httpClient.request<void>(`/v1/teams/${id}`, {
      method: 'DELETE',
    });
  },
};
