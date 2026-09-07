import type { UserRole } from '@/modules/user/types';

/** Mirrors API service.CanManageTarget. */
export function canManageTarget(actorRole: UserRole, targetRole: UserRole): boolean {
  if (actorRole === 'super_admin') return true;
  if (actorRole === 'admin') return targetRole === 'admin' || targetRole === 'user';
  return false;
}
