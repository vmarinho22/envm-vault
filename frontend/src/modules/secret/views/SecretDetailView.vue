<script setup lang="ts">
import { ArrowLeftIcon, LinkSlashIcon } from '@heroicons/vue/24/outline';
import { onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { RouterLink } from 'vue-router';
import { useSecretDetail } from '@/modules/secret/composables/useSecretDetail';
import { formatDateTime } from '@/shared/datetime/format';

const { t } = useI18n();
const { secret, envBadges, loading, mutating, error, load, detachProject } = useSecretDetail();

onMounted(() => {
  void load();
});
</script>

<template>
  <section class="px-8 py-10 md:px-12">
    <RouterLink
      :to="{ name: 'secrets' }"
      class="inline-flex items-center gap-2 text-sm text-base-content/70 hover:text-base-content"
    >
      <ArrowLeftIcon class="size-4" aria-hidden="true" />
      {{ t('secrets.catalog.back') }}
    </RouterLink>

    <p v-if="error" class="mt-6 text-sm text-error" role="alert">{{ error }}</p>
    <p v-else-if="loading" class="mt-8 text-base-content/70">{{ t('common.loading') }}</p>
    <p v-else-if="!secret" class="mt-8 text-base-content/70">
      {{ t('secrets.catalog.notFound') }}
    </p>

    <template v-else>
      <div class="mt-4">
        <h1 class="font-mono text-2xl font-medium tracking-tight md:text-3xl">
          {{ secret.name }}
        </h1>
        <p class="mt-2 text-sm text-base-content/60">
          {{
            t('secrets.catalog.detail.meta', {
              count: secret.project_count,
              created: formatDateTime(secret.created_at),
            })
          }}
        </p>
      </div>

      <div class="mt-8">
        <h2 class="text-lg font-medium">{{ t('secrets.catalog.detail.environments') }}</h2>
        <div class="mt-3 flex flex-wrap gap-1">
          <span
            v-for="title in envBadges"
            :key="title"
            class="badge badge-ghost badge-sm"
          >
            {{ title }}
          </span>
          <span v-if="envBadges.length === 0" class="text-sm text-base-content/60">
            {{ t('secrets.catalog.detail.noValues') }}
          </span>
        </div>
      </div>

      <div class="mt-8">
        <h2 class="text-lg font-medium">{{ t('secrets.catalog.detail.projects') }}</h2>
        <div class="mt-4 overflow-x-auto">
          <table class="table">
            <thead>
              <tr>
                <th>{{ t('secrets.catalog.detail.projectTitle') }}</th>
                <th>{{ t('secrets.catalog.detail.projectSlug') }}</th>
                <th>
                  <span class="sr-only">{{ t('secrets.catalog.columns.actions') }}</span>
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="secret.projects.length === 0">
                <td colspan="3" class="text-base-content/70">
                  {{ t('secrets.catalog.detail.noProjects') }}
                </td>
              </tr>
              <tr v-for="project in secret.projects" :key="project.id">
                <td>
                  <RouterLink
                    class="link link-hover font-medium"
                    :to="{ name: 'project-detail', params: { slug: project.slug } }"
                  >
                    {{ project.title }}
                  </RouterLink>
                </td>
                <td class="font-mono text-sm text-base-content/70">{{ project.slug }}</td>
                <td class="text-right">
                  <button
                    class="btn btn-error btn-outline btn-sm gap-1"
                    type="button"
                    :disabled="mutating"
                    :aria-label="t('secrets.catalog.unlink')"
                    @click="detachProject(project)"
                  >
                    <LinkSlashIcon class="size-4" aria-hidden="true" />
                    {{ t('secrets.catalog.unlink') }}
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>
  </section>
</template>
