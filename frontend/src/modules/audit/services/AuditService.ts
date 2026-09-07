import type { AuditLogListParams, AuditLogListResponse } from '@/modules/audit/types';
import { httpClient } from '@/shared/http';

export interface AuditService {
  list(params?: AuditLogListParams): Promise<AuditLogListResponse>;
}

function buildQuery(params: AuditLogListParams = {}): string {
  const search = new URLSearchParams();
  if (params.page !== null) search.set('page', String(params.page));
  if (params.page_size != null) search.set('page_size', String(params.page_size));
  if (params.sort) search.set('sort', params.sort);
  for (const id of params.actor_user_ids ?? []) {
    search.append('actor_user_id', id);
  }
  for (const id of params.target_user_ids ?? []) {
    search.append('target_user_id', id);
  }
  if (params.occurred_from) search.set('occurred_from', params.occurred_from);
  if (params.occurred_to) search.set('occurred_to', params.occurred_to);
  const q = search.toString();
  return q ? `?${q}` : '';
}

export const auditService: AuditService = {
  list(params?: AuditLogListParams): Promise<AuditLogListResponse> {
    return httpClient.request<AuditLogListResponse>(`/v1/audit-logs${buildQuery(params)}`);
  },
};
