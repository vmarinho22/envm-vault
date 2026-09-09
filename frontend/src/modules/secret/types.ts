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

export interface CatalogSearchResponse {
  items: CatalogSecret[];
}

export interface SearchSecretsParams {
  title?: string;
  excludeProjectId?: string;
  limit?: number;
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
