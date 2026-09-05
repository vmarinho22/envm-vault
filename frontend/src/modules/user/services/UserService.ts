import type { User, UserRole, UserWithToken } from '@/modules/user/types';
import { httpClient } from '@/shared/http';

export type RegisterableRole = Exclude<UserRole, 'super_admin'>;

export interface UserService {
  register(email: string, role: RegisterableRole): Promise<User>;
  list(): Promise<UserWithToken[]>;
}

export const userService: UserService = {
  register(email: string, role: RegisterableRole): Promise<User> {
    return httpClient.request<User>('/v1/auth/register', {
      method: 'POST',
      body: { email, role },
    });
  },

  list(): Promise<UserWithToken[]> {
    return httpClient.request<UserWithToken[]>('/v1/users');
  },
};
