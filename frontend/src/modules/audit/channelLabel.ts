import type { AuthChannel } from '@/modules/audit/types';

export function channelLabel(channel: AuthChannel, t: (key: string) => string): string {
  if (channel === 'personal_token') {
    return t('audit.channels.api');
  }
  return t('audit.channels.web');
}
