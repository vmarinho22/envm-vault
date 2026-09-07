export interface SecretValueMeta {
  has_value: boolean;
  updated_at?: string;
}

export interface Secret {
  id: string;
  name: string;
  values: Record<string, SecretValueMeta>;
  created_at: string;
  updated_at: string;
}

export interface SecretListResponse {
  items: Secret[];
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
