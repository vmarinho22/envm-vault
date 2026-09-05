export interface TokenStatus {
  active: boolean;
  prefix?: string;
  created_at?: string;
}

export interface TokenCreateResponse {
  token: string;
  prefix: string;
  created_at: string;
}
