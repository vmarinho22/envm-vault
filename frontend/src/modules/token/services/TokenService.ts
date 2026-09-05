import type { TokenCreateResponse, TokenStatus } from '@/modules/token/types';
import { httpClient } from '@/shared/http';

export interface TokenService {
  getMine(): Promise<TokenStatus>;
  create(): Promise<TokenCreateResponse>;
  revokeOwn(): Promise<void>;
  revokeForUser(userId: string): Promise<void>;
}

export const tokenService: TokenService = {
  getMine(): Promise<TokenStatus> {
    return httpClient.request<TokenStatus>('/v1/tokens/me');
  },

  create(): Promise<TokenCreateResponse> {
    return httpClient.request<TokenCreateResponse>('/v1/tokens', {
      method: 'POST',
    });
  },

  revokeOwn(): Promise<void> {
    return httpClient.request<void>('/v1/tokens', {
      method: 'DELETE',
    });
  },

  revokeForUser(userId: string): Promise<void> {
    return httpClient.request<void>(`/v1/users/${userId}/token`, {
      method: 'DELETE',
    });
  },
};
