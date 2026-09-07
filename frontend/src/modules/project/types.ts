export interface ProjectTeam {
  id: string;
  name: string;
  slug: string;
}

export interface ProjectEnvironment {
  id: string;
  title: string;
}

export interface ProjectTag {
  id: string;
  title: string;
  color: string;
}

export interface Project {
  id: string;
  title: string;
  slug: string;
  teams: ProjectTeam[];
  environments: ProjectEnvironment[];
  tags: ProjectTag[];
  created_at: string;
  updated_at: string;
}

export interface ProjectListResponse {
  items: Project[];
  page: number;
  page_size: number;
  total: number;
}

export interface CreateProjectInput {
  title: string;
  team_ids: string[];
  environment_ids: string[];
  tag_ids: string[];
}

export interface UpdateProjectInput {
  title?: string;
  team_ids?: string[];
  environment_ids?: string[];
  tag_ids?: string[];
}

export interface ListProjectsParams {
  title?: string;
  team_ids?: string[];
  tag_ids?: string[];
  page?: number;
  page_size?: number;
}

export const PROJECT_PAGE_SIZE = 12;
