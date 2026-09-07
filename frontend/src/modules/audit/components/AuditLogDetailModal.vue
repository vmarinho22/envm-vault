<script setup lang="ts">
import { watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { actionLabel } from '@/modules/audit/actionLabel';
import { channelLabel } from '@/modules/audit/channelLabel';
import type { AuditLog } from '@/modules/audit/types';
import { formatDateTime } from '@/shared/datetime/format';

const props = defineProps<{
  open: boolean;
  log: AuditLog | null;
}>();

const emit = defineEmits<{
  close: [];
}>();

const { t } = useI18n();

watch(
  () => props.open,
  (isOpen) => {
    const dialog = document.getElementById('audit-log-detail-modal') as HTMLDialogElement | null;
    if (!dialog) return;
    if (isOpen) {
      if (!dialog.open) dialog.showModal();
    } else if (dialog.open) {
      dialog.close();
    }
  },
);

function onClose(): void {
  emit('close');
}

function dash(value: string | null | undefined): string {
  if (value == null || value === '') return '—';
  return value;
}
</script>

<template>
  <dialog id="audit-log-detail-modal" class="modal" @close="onClose">
    <div class="modal-box flex max-h-[85vh] max-w-2xl flex-col">
      <h3 class="shrink-0 text-lg font-medium">{{ t('audit.detail.title') }}</h3>
      <dl v-if="log" class="mt-4 min-h-0 flex-1 space-y-3 overflow-y-auto text-sm">
        <div>
          <dt class="text-base-content/60">{{ t('audit.detail.id') }}</dt>
          <dd class="font-mono break-all" :title="log.id">{{ log.id }}</dd>
        </div>
        <div>
          <dt class="text-base-content/60">{{ t('audit.columns.action') }}</dt>
          <dd>{{ actionLabel(log.action, t) }}</dd>
          <dd class="mt-0.5 font-mono text-xs text-base-content/50">
            {{ t('audit.detail.actionCode') }}: {{ log.action }}
          </dd>
        </div>
        <div>
          <dt class="text-base-content/60">{{ t('audit.detail.actorUserId') }}</dt>
          <dd class="font-mono break-all" :title="dash(log.actor_user_id)">
            {{ dash(log.actor_user_id) }}
          </dd>
        </div>
        <div>
          <dt class="text-base-content/60">{{ t('audit.columns.actor') }}</dt>
          <dd>{{ dash(log.actor_email) }}</dd>
        </div>
        <div>
          <dt class="text-base-content/60">{{ t('audit.detail.actorRole') }}</dt>
          <dd>{{ dash(log.actor_role) }}</dd>
        </div>
        <div>
          <dt class="text-base-content/60">{{ t('audit.detail.targetUserId') }}</dt>
          <dd class="font-mono break-all" :title="dash(log.target_user_id)">
            {{ dash(log.target_user_id) }}
          </dd>
        </div>
        <div>
          <dt class="text-base-content/60">{{ t('audit.columns.target') }}</dt>
          <dd>{{ dash(log.target_email) }}</dd>
        </div>
        <div>
          <dt class="text-base-content/60">{{ t('audit.detail.projectId') }}</dt>
          <dd class="font-mono break-all" :title="dash(log.project_id)">
            {{ dash(log.project_id) }}
          </dd>
        </div>
        <div>
          <dt class="text-base-content/60">{{ t('audit.columns.project') }}</dt>
          <dd>{{ dash(log.project_title) }}</dd>
        </div>
        <div>
          <dt class="text-base-content/60">{{ t('audit.detail.projectSlug') }}</dt>
          <dd class="font-mono">{{ dash(log.project_slug) }}</dd>
        </div>
        <div>
          <dt class="text-base-content/60">{{ t('audit.detail.secretId') }}</dt>
          <dd class="font-mono break-all" :title="dash(log.secret_id)">
            {{ dash(log.secret_id) }}
          </dd>
        </div>
        <div>
          <dt class="text-base-content/60">{{ t('audit.columns.secret') }}</dt>
          <dd class="font-mono">{{ dash(log.secret_name) }}</dd>
        </div>
        <div>
          <dt class="text-base-content/60">{{ t('audit.columns.channel') }}</dt>
          <dd>
            {{ channelLabel(log.auth_channel, t) }}
            <span class="text-base-content/50">({{ log.auth_channel }})</span>
          </dd>
        </div>
        <div>
          <dt class="text-base-content/60">{{ t('audit.detail.ip') }}</dt>
          <dd class="font-mono">{{ dash(log.ip) }}</dd>
        </div>
        <div>
          <dt class="text-base-content/60">{{ t('audit.detail.userAgent') }}</dt>
          <dd class="break-all" :title="dash(log.user_agent)">{{ dash(log.user_agent) }}</dd>
        </div>
        <div>
          <dt class="text-base-content/60">{{ t('audit.detail.origin') }}</dt>
          <dd class="break-all" :title="dash(log.origin)">{{ dash(log.origin) }}</dd>
        </div>
        <div>
          <dt class="text-base-content/60">{{ t('audit.columns.occurredAt') }}</dt>
          <dd>{{ formatDateTime(log.occurred_at) }}</dd>
        </div>
        <div>
          <dt class="text-base-content/60">{{ t('audit.detail.createdAt') }}</dt>
          <dd>{{ formatDateTime(log.created_at) }}</dd>
        </div>
      </dl>
      <div class="modal-action shrink-0">
        <form method="dialog">
          <button class="btn" type="submit">{{ t('audit.detail.close') }}</button>
        </form>
      </div>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button type="submit">{{ t('audit.detail.close') }}</button>
    </form>
  </dialog>
</template>
