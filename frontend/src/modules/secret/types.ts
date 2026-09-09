export const SECRET_PAGE_SIZE = 12;

export interface SecretValueMeta {
  has_value: boolean;
  updated_at?: string;
}

export interface Secret {
  id: string;
  name: string;
  values: Record<string, SecretValueMeta>;
  project_count: number;
  created_at: string;
  updated_at: string;
}

export interface SecretListResponse {
  items: Secret[];
}

export interface CatalogSecret {
  id: string;
  name: string;
  project_count: number;
  created_at: string;
  updated_at: string;
}

export interface CatalogListResponse {
  items: CatalogSecret[];
  page: number;
  page_size: number;
  total: number;
}

export interface CatalogListParams {
  title?: string;
  excludeProjectId?: string;
  page?: number;
  pageSize?: number;
}

export interface SecretProjectLink {
  id: string;
  title: string;
  slug: string;
}

export interface SecretDetail {
  id: string;
  name: string;
  project_count: number;
  values: Record<string, SecretValueMeta>;
  projects: SecretProjectLink[];
  created_at: string;
  updated_at: string;
}

export interface CreateSecretInput {
  name: string;
  values?: Record<string, string>;
}

export interface RevealSecretResponse {
  name: string;
  value: string;
}

export interface DownloadItem {
  name: string;
  value: string;
}

export interface DownloadResponse {
  items: DownloadItem[];
}
