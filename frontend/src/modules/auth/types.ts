import type { User } from '@/modules/user/types';

export interface LoginResponse {
  message: string;
}

export interface VerifyResponse {
  token: string;
  user: User;
}
