import { describe, expect, it } from 'vitest';
import { actionLabel } from '@/modules/audit/actionLabel';

const t = ((key: string) => key) as (key: string) => string;

describe('actionLabel', () => {
  it('maps known actions to i18n keys', () => {
    expect(actionLabel('personal_token.rotated', t)).toBe('audit.actions.personalTokenRotated');
    expect(actionLabel('login.succeeded', t)).toBe('audit.actions.loginSucceeded');
    expect(actionLabel('team.created', t)).toBe('audit.actions.teamCreated');
    expect(actionLabel('tag.created', t)).toBe('audit.actions.tagCreated');
    expect(actionLabel('environment.created', t)).toBe('audit.actions.environmentCreated');
    expect(actionLabel('project.created', t)).toBe('audit.actions.projectCreated');
    expect(actionLabel('secret.created', t)).toBe('audit.actions.secretCreated');
    expect(actionLabel('secret_value.set', t)).toBe('audit.actions.secretValueSet');
    expect(actionLabel('user.blocked', t)).toBe('audit.actions.userBlocked');
  });

  it('falls back to raw action when unknown', () => {
    expect(actionLabel('unknown.event', t)).toBe('unknown.event');
  });
});
