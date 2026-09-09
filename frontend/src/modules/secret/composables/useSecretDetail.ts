import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute } from 'vue-router';
import { environmentService } from '@/modules/environment/services/EnvironmentService';
import { secretService } from '@/modules/secret/services/SecretService';
import type { SecretDetail, SecretProjectLink } from '@/modules/secret/types';
import { useConfirmDialog } from '@/shared/confirm/useConfirmDialog';
import { ApiError } from '@/shared/http/ApiError';

export function useSecretDetail() {
  const { t } = useI18n();
  const route = useRoute();
  const { confirm } = useConfirmDialog();

  const secret = ref<SecretDetail | null>(null);
  const envTitles = ref<Record<string, string>>({});
  const loading = ref(false);
  const mutating = ref(false);
  const error = ref<string | null>(null);

  const secretId = computed(() => String(route.params.id ?? ''));

  const envBadges = computed(() => {
    if (!secret.value) return [];
    const titles: string[] = [];
    for (const [envId, meta] of Object.entries(secret.value.values)) {
      if (!meta.has_value) continue;
      titles.push(envTitles.value[envId] ?? envId);
    }
    return titles.sort((a, b) => a.localeCompare(b));
  });

  function toErrorMessage(err: unknown, fallback: string): string {
    if (err instanceof ApiError) return err.message;
    return fallback;
  }

  async function load(): Promise<void> {
    error.value = null;
    loading.value = true;
    try {
      const [detail, environments] = await Promise.all([
        secretService.getById(secretId.value),
        environmentService.list(),
      ]);
      secret.value = detail;
      envTitles.value = Object.fromEntries(environments.map((e) => [e.id, e.title]));
    } catch (err) {
      error.value = toErrorMessage(err, t('secrets.catalog.errors.loadDetailFailed'));
      secret.value = null;
    } finally {
      loading.value = false;
    }
  }

  async function detachProject(project: SecretProjectLink): Promise<void> {
    if (!secret.value) return;
    const ok = await confirm({
      title: t('secrets.catalog.unlink'),
      message: t('secrets.catalog.confirmUnlink', {
        secret: secret.value.name,
        project: project.title,
      }),
      confirmLabel: t('common.confirm'),
      cancelLabel: t('common.cancel'),
      variant: 'danger',
    });
    if (!ok) return;

    error.value = null;
    mutating.value = true;
    try {
      await secretService.remove(project.id, secret.value.id);
      await load();
    } catch (err) {
      error.value = toErrorMessage(err, t('secrets.catalog.errors.unlinkFailed'));
    } finally {
      mutating.value = false;
    }
  }

  return {
    secret,
    envBadges,
    loading,
    mutating,
    error,
    secretId,
    load,
    detachProject,
  };
}
