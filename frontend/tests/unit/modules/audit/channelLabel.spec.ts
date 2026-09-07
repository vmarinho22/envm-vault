import { describe, expect, it } from 'vitest';
import { channelLabel } from '@/modules/audit/channelLabel';

const t = ((key: string) => {
  const map: Record<string, string> = {
    'audit.channels.web': 'Web',
    'audit.channels.api': 'API',
  };
  return map[key] ?? key;
}) as (key: string) => string;

describe('channelLabel', () => {
  it('maps personal_token to api', () => {
    expect(channelLabel('personal_token', t as never)).toBe('API');
  });

  it('maps login channels to web', () => {
    expect(channelLabel('login_attempt', t as never)).toBe('Web');
    expect(channelLabel('logged_in_user', t as never)).toBe('Web');
  });
});
