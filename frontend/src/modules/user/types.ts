export type UserRole = 'super_admin' | 'admin' | 'user';

export interface User {
  id: string;
  email: string;
  role: UserRole;
}

export interface UserTokenStatus {
  prefix: string;
  created_at: string;
  revoked_at?: string;
}

export interface UserWithToken extends User {
  token: UserTokenStatus | null;
}
