import { describe, expect, it } from 'vitest';
import { actionLabel } from '@/modules/audit/actionLabel';

const t = ((key: string) => key) as (key: string) => string;

describe('actionLabel', () => {
  it('maps known actions to i18n keys', () => {
    expect(actionLabel('personal_token.rotated', t)).toBe('audit.actions.personalTokenRotated');
    expect(actionLabel('login.succeeded', t)).toBe('audit.actions.loginSucceeded');
  });

  it('falls back to raw action when unknown', () => {
    expect(actionLabel('unknown.event', t)).toBe('unknown.event');
  });
});
