import type { HealthResponse } from '@/modules/health/types';
import { httpClient } from '@/shared/http';

export interface HealthService {
  check(): Promise<HealthResponse>;
}

export const healthService: HealthService = {
  check(): Promise<HealthResponse> {
    return httpClient.request<HealthResponse>('/v1/health');
  },
};
