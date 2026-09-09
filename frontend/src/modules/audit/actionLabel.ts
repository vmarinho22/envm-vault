/** Known audit action codes from the API. */
export const AUDIT_ACTIONS = [
  'login.requested',
  'login.succeeded',
  'login.failed',
  'user.bootstrap_super_admin',
  'user.registered',
  'user.blocked',
  'user.unblocked',
  'personal_token.created',
  'personal_token.rotated',
  'personal_token.revoked',
  'personal_token.revoked_by_admin',
  'team.created',
  'team.updated',
  'team.deleted',
  'tag.created',
  'tag.updated',
  'tag.deleted',
  'environment.created',
  'environment.updated',
  'environment.deleted',
  'project.created',
  'project.updated',
  'project.deleted',
  'secret.created',
  'secret.updated',
  'secret.deleted',
  'secret.attached',
  'secret.detached',
  'secret_value.set',
  'secret_value.cleared',
  'secret.downloaded',
] as const;

export type AuditAction = (typeof AUDIT_ACTIONS)[number];

const ACTION_I18N_KEYS: Record<string, string> = {
  'login.requested': 'audit.actions.loginRequested',
  'login.succeeded': 'audit.actions.loginSucceeded',
  'login.failed': 'audit.actions.loginFailed',
  'user.bootstrap_super_admin': 'audit.actions.userBootstrapSuperAdmin',
  'user.registered': 'audit.actions.userRegistered',
  'user.blocked': 'audit.actions.userBlocked',
  'user.unblocked': 'audit.actions.userUnblocked',
  'personal_token.created': 'audit.actions.personalTokenCreated',
  'personal_token.rotated': 'audit.actions.personalTokenRotated',
  'personal_token.revoked': 'audit.actions.personalTokenRevoked',
  'personal_token.revoked_by_admin': 'audit.actions.personalTokenRevokedByAdmin',
  'team.created': 'audit.actions.teamCreated',
  'team.updated': 'audit.actions.teamUpdated',
  'team.deleted': 'audit.actions.teamDeleted',
  'tag.created': 'audit.actions.tagCreated',
  'tag.updated': 'audit.actions.tagUpdated',
  'tag.deleted': 'audit.actions.tagDeleted',
  'environment.created': 'audit.actions.environmentCreated',
  'environment.updated': 'audit.actions.environmentUpdated',
  'environment.deleted': 'audit.actions.environmentDeleted',
  'project.created': 'audit.actions.projectCreated',
  'project.updated': 'audit.actions.projectUpdated',
  'project.deleted': 'audit.actions.projectDeleted',
  'secret.created': 'audit.actions.secretCreated',
  'secret.updated': 'audit.actions.secretUpdated',
  'secret.deleted': 'audit.actions.secretDeleted',
  'secret.attached': 'audit.actions.secretAttached',
  'secret.detached': 'audit.actions.secretDetached',
  'secret_value.set': 'audit.actions.secretValueSet',
  'secret_value.cleared': 'audit.actions.secretValueCleared',
  'secret.downloaded': 'audit.actions.secretDownloaded',
};

export function actionLabel(action: string, t: (key: string) => string): string {
  const key = ACTION_I18N_KEYS[action];
  if (!key) {
    return action;
  }
  return t(key);
}
