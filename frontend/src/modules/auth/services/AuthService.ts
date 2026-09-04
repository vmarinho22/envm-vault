import type { LoginResponse, VerifyResponse } from '@/modules/auth/types';
import { httpClient } from '@/shared/http';

export interface AuthService {
  requestLogin(email: string): Promise<LoginResponse>;
  verifyLogin(email: string, code: string): Promise<VerifyResponse>;
}

export const authService: AuthService = {
  requestLogin(email: string): Promise<LoginResponse> {
    return httpClient.request<LoginResponse>('/v1/auth/login', {
      method: 'POST',
      body: { email },
    });
  },

  verifyLogin(email: string, code: string): Promise<VerifyResponse> {
    return httpClient.request<VerifyResponse>('/v1/auth/verify', {
      method: 'POST',
      body: { email, code },
    });
  },
};
