export type UserRole = 'super_admin' | 'admin' | 'user';

export interface User {
  id: string;
  email: string;
  role: UserRole;
}
