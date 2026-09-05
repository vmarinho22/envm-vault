import { describe, expect, it } from 'vitest';
import { roleLabel } from '@/modules/user/roleLabel';

describe('roleLabel', () => {
  it('maps known roles to Portuguese labels', () => {
    expect(roleLabel('super_admin')).toBe('Super admin');
    expect(roleLabel('admin')).toBe('Admin');
    expect(roleLabel('user')).toBe('Usuário');
  });
});
