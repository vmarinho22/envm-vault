export type AuthChannel = 'login_attempt' | 'logged_in_user' | 'personal_token' | string;

export type AuditSort = 'asc' | 'desc';

export interface AuditLog {
  id: string;
  action: string;
  actor_user_id: string | null;
  actor_email: string;
  actor_role: string;
  target_user_id: string | null;
  target_email: string | null;
  auth_channel: AuthChannel;
  ip: string;
  user_agent: string;
  origin: string;
  occurred_at: string;
  created_at: string;
}

export interface AuditLogListResponse {
  items: AuditLog[];
  page: number;
  page_size: number;
  total: number;
}

export interface AuditLogListParams {
  page?: number;
  page_size?: number;
  sort?: AuditSort;
  actor_user_ids?: string[];
  target_user_ids?: string[];
  occurred_from?: string;
  occurred_to?: string;
}

/** Client-side filter criteria for the audit logs list. */
export interface AuditListFilters {
  actorUserIds: string[];
  targetUserIds: string[];
  occurredFrom: string;
  occurredTo: string;
}

/** Client-side pagination/sort for the audit logs list. */
export interface AuditListQuery {
  page: number;
  sort: AuditSort;
}
