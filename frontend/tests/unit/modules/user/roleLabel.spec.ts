import { createI18n } from 'vue-i18n';
import { describe, expect, it } from 'vitest';
import { allLocaleMessages, type AppLocale } from '@/app/translations';
import { roleLabel } from '@/modules/user/roleLabel';

function createTestI18n(locale: AppLocale) {
  return createI18n({
    legacy: false,
    locale,
    fallbackLocale: 'en-US',
    messages: allLocaleMessages(),
  });
}

describe('roleLabel', () => {
  it('maps known roles to Portuguese labels', () => {
    const { t } = createTestI18n('pt-BR').global;
    expect(roleLabel('super_admin', t)).toBe('Super admin');
    expect(roleLabel('admin', t)).toBe('Admin');
    expect(roleLabel('user', t)).toBe('Usuário');
  });

  it('maps known roles to English labels', () => {
    const { t } = createTestI18n('en-US').global;
    expect(roleLabel('super_admin', t)).toBe('Super admin');
    expect(roleLabel('admin', t)).toBe('Admin');
    expect(roleLabel('user', t)).toBe('User');
  });
});
