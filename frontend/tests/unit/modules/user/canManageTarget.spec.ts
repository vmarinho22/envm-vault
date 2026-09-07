import { describe, expect, it } from 'vitest';
import { canManageTarget } from '@/modules/user/canManageTarget';

describe('canManageTarget', () => {
  it('allows super_admin to manage anyone', () => {
    expect(canManageTarget('super_admin', 'super_admin')).toBe(true);
    expect(canManageTarget('super_admin', 'admin')).toBe(true);
    expect(canManageTarget('super_admin', 'user')).toBe(true);
  });

  it('allows admin to manage admin and user only', () => {
    expect(canManageTarget('admin', 'admin')).toBe(true);
    expect(canManageTarget('admin', 'user')).toBe(true);
    expect(canManageTarget('admin', 'super_admin')).toBe(false);
  });

  it('denies user role', () => {
    expect(canManageTarget('user', 'user')).toBe(false);
  });
});
