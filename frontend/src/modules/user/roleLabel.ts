import type { UserRole } from '@/modules/user/types';

const ROLE_LABELS: Record<UserRole, string> = {
  super_admin: 'Super admin',
  admin: 'Admin',
  user: 'Usuário',
};

export function roleLabel(role: UserRole): string {
  return ROLE_LABELS[role];
}
