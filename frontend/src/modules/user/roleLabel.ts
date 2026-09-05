import type { UserRole } from '@/modules/user/types';

type TranslateFn = (key: string) => string;

export function roleLabel(role: UserRole, t: TranslateFn): string {
  return t(`users.roles.${role}`);
}
