<script setup lang="ts">
import { ArrowPathIcon } from '@heroicons/vue/24/outline';
import { computed, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import CatalogMultiSelect from '@/modules/project/components/CatalogMultiSelect.vue';
import type { CatalogOption } from '@/modules/project/components/CatalogMultiSelect.vue';
import type { CatalogQuickKind } from '@/modules/project/composables/useProjects';
import { DEFAULT_TAG_COLOR, randomTagColor } from '@/modules/tag/types';
import type { UserWithToken } from '@/modules/user/types';
import { fromName } from '@/shared/slug/fromName';

const props = defineProps<{
  open: boolean;
  mode: 'create' | 'edit';
  title: string;
  teamIds: string[];
  environmentIds: string[];
  tagIds: string[];
  teamOptions: CatalogOption[];
  environmentOptions: CatalogOption[];
  tagOptions: CatalogOption[];
  saving: boolean;
  canAddTeam?: boolean;
  canAddEnvironment?: boolean;
  canAddTag?: boolean;
  catalogQuick: CatalogQuickKind | null;
  quickUsers: UserWithToken[];
  quickSaving: boolean;
  quickError: string | null;
  teamQuickName: string;
  teamQuickAccountableId: string;
  environmentQuickTitle: string;
  tagQuickTitle: string;
  tagQuickColor: string;
}>();

const emit = defineEmits<{
  close: [];
  submit: [];
  'update:title': [value: string];
  'update:teamIds': [value: string[]];
  'update:environmentIds': [value: string[]];
  'update:tagIds': [value: string[]];
  addTeam: [];
  addEnvironment: [];
  addTag: [];
  closeCatalogQuick: [];
  submitQuickTeam: [];
  submitQuickEnvironment: [];
  submitQuickTag: [];
  'update:teamQuickName': [value: string];
  'update:teamQuickAccountableId': [value: string];
  'update:environmentQuickTitle': [value: string];
  'update:tagQuickTitle': [value: string];
  'update:tagQuickColor': [value: string];
}>();

const { t } = useI18n();

const slugPreview = computed(() => fromName(props.title));
const teamSlugPreview = computed(() => fromName(props.teamQuickName));

watch(
  () => props.open,
  (isOpen) => {
    const dialog = document.getElementById('project-form-modal') as HTMLDialogElement | null;
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

function normalizeHex(raw: string): string {
  const value = raw.trim();
  if (/^#[0-9A-Fa-f]{6}$/.test(value)) {
    return value.toUpperCase();
  }
  return value;
}
</script>

<template>
  <dialog id="project-form-modal" class="modal" @close="onClose">
    <div class="modal-box flex max-h-[85vh] max-w-lg flex-col overflow-hidden">
      <template v-if="catalogQuick === 'team'">
        <h3 class="shrink-0 text-lg font-medium">{{ t('teams.form.createTitle') }}</h3>
        <form class="mt-4 flex min-h-0 flex-1 flex-col" @submit.prevent="emit('submitQuickTeam')">
          <div class="min-h-0 flex-1 space-y-4 overflow-y-auto pr-1">
            <p v-if="quickError" class="text-sm text-error" role="alert">{{ quickError }}</p>
            <label class="form-control w-full">
              <span class="label-text my-1">{{ t('teams.form.name') }}</span>
              <input
                class="input input-bordered w-full"
                type="text"
                name="quick-team-name"
                autocomplete="off"
                required
                maxlength="80"
                :value="teamQuickName"
                :aria-label="t('teams.form.name')"
                @input="emit('update:teamQuickName', ($event.target as HTMLInputElement).value)"
              />
            </label>
            <p class="text-sm text-base-content/60">
              {{ t('teams.form.slugPreview') }}:
              <span class="font-mono">{{ teamSlugPreview || t('common.emDash') }}</span>
            </p>
            <label class="form-control w-full">
              <span class="label-text my-1">{{ t('teams.form.accountable') }}</span>
              <select
                class="select select-bordered w-full"
                name="quick-team-accountable"
                required
                :value="teamQuickAccountableId"
                :aria-label="t('teams.form.accountable')"
                @change="
                  emit('update:teamQuickAccountableId', ($event.target as HTMLSelectElement).value)
                "
              >
                <option disabled value="">{{ t('teams.form.selectUser') }}</option>
                <option v-for="user in quickUsers" :key="user.id" :value="user.id">
                  {{ user.email }}
                </option>
              </select>
            </label>
          </div>
          <div class="modal-action shrink-0">
            <button class="btn" type="button" @click="emit('closeCatalogQuick')">
              {{ t('projects.form.back') }}
            </button>
            <button class="btn btn-primary" type="submit" :disabled="quickSaving">
              <span v-if="quickSaving" class="loading loading-spinner loading-xs" />
              {{ t('teams.form.create') }}
            </button>
          </div>
        </form>
      </template>

      <template v-else-if="catalogQuick === 'environment'">
        <h3 class="shrink-0 text-lg font-medium">{{ t('environments.form.createTitle') }}</h3>
        <form
          class="mt-4 flex min-h-0 flex-1 flex-col"
          @submit.prevent="emit('submitQuickEnvironment')"
        >
          <div class="min-h-0 flex-1 space-y-4 overflow-y-auto pr-1">
            <p v-if="quickError" class="text-sm text-error" role="alert">{{ quickError }}</p>
            <label class="form-control w-full">
              <span class="label-text my-1">{{ t('environments.form.title') }}</span>
              <input
                class="input input-bordered w-full"
                type="text"
                name="quick-environment-title"
                autocomplete="off"
                required
                maxlength="80"
                :value="environmentQuickTitle"
                :aria-label="t('environments.form.title')"
                @input="
                  emit('update:environmentQuickTitle', ($event.target as HTMLInputElement).value)
                "
              />
            </label>
          </div>
          <div class="modal-action shrink-0">
            <button class="btn" type="button" @click="emit('closeCatalogQuick')">
              {{ t('projects.form.back') }}
            </button>
            <button class="btn btn-primary" type="submit" :disabled="quickSaving">
              <span v-if="quickSaving" class="loading loading-spinner loading-xs" />
              {{ t('environments.form.create') }}
            </button>
          </div>
        </form>
      </template>

      <template v-else-if="catalogQuick === 'tag'">
        <h3 class="shrink-0 text-lg font-medium">{{ t('tags.form.createTitle') }}</h3>
        <form class="mt-4 flex min-h-0 flex-1 flex-col" @submit.prevent="emit('submitQuickTag')">
          <div class="min-h-0 flex-1 space-y-4 overflow-y-auto pr-1">
            <p v-if="quickError" class="text-sm text-error" role="alert">{{ quickError }}</p>
            <label class="form-control w-full">
              <span class="label-text my-1">{{ t('tags.form.title') }}</span>
              <input
                class="input input-bordered w-full"
                type="text"
                name="quick-tag-title"
                autocomplete="off"
                required
                maxlength="80"
                :value="tagQuickTitle"
                :aria-label="t('tags.form.title')"
                @input="emit('update:tagQuickTitle', ($event.target as HTMLInputElement).value)"
              />
            </label>
            <div class="form-control w-full">
              <span class="label-text my-1">{{ t('tags.form.color') }}</span>
              <div class="mt-1 flex items-center gap-2">
                <input
                  class="h-10 w-14 cursor-pointer rounded border border-base-300 bg-transparent p-1"
                  type="color"
                  name="quick-tag-color-picker"
                  :value="/^#[0-9A-Fa-f]{6}$/.test(tagQuickColor) ? tagQuickColor : DEFAULT_TAG_COLOR"
                  :aria-label="t('tags.form.colorPicker')"
                  @input="
                    emit(
                      'update:tagQuickColor',
                      ($event.target as HTMLInputElement).value.toUpperCase(),
                    )
                  "
                />
                <input
                  class="input input-bordered min-w-0 flex-1 font-mono uppercase"
                  type="text"
                  name="quick-tag-color-hex"
                  autocomplete="off"
                  required
                  maxlength="7"
                  pattern="#[0-9A-Fa-f]{6}"
                  :value="tagQuickColor"
                  :aria-label="t('tags.form.colorHex')"
                  @input="
                    emit(
                      'update:tagQuickColor',
                      normalizeHex(($event.target as HTMLInputElement).value),
                    )
                  "
                />
                <button
                  class="btn btn-ghost btn-square"
                  type="button"
                  :aria-label="t('tags.form.randomizeColor')"
                  :title="t('tags.form.randomizeColor')"
                  @click="emit('update:tagQuickColor', randomTagColor())"
                >
                  <ArrowPathIcon class="size-5" aria-hidden="true" />
                </button>
              </div>
            </div>
          </div>
          <div class="modal-action shrink-0">
            <button class="btn" type="button" @click="emit('closeCatalogQuick')">
              {{ t('projects.form.back') }}
            </button>
            <button class="btn btn-primary" type="submit" :disabled="quickSaving">
              <span v-if="quickSaving" class="loading loading-spinner loading-xs" />
              {{ t('tags.form.create') }}
            </button>
          </div>
        </form>
      </template>

      <template v-else>
        <h3 class="shrink-0 text-lg font-medium">
          {{ mode === 'create' ? t('projects.form.createTitle') : t('projects.form.editTitle') }}
        </h3>

        <form class="mt-4 flex min-h-0 flex-1 flex-col" @submit.prevent="emit('submit')">
          <div class="min-h-0 flex-1 space-y-4 overflow-y-auto pr-1">
            <label class="form-control w-full">
              <span class="label-text my-1">{{ t('projects.form.title') }}</span>
              <input
                class="input input-bordered w-full"
                type="text"
                name="project-title"
                autocomplete="off"
                required
                maxlength="80"
                :value="title"
                :aria-label="t('projects.form.title')"
                @input="emit('update:title', ($event.target as HTMLInputElement).value)"
              />
            </label>

            <p class="text-sm text-base-content/60">
              {{ t('projects.form.slugPreview') }}:
              <span class="font-mono">{{ slugPreview || t('common.emDash') }}</span>
            </p>

            <CatalogMultiSelect
              :model-value="teamIds"
              :options="teamOptions"
              :label="t('projects.form.teams')"
              :placeholder="t('projects.form.selectTeams')"
              :can-add="canAddTeam"
              :add-aria-label="t('projects.form.addTeam')"
              @update:model-value="emit('update:teamIds', $event)"
              @add="emit('addTeam')"
            />

            <CatalogMultiSelect
              :model-value="environmentIds"
              :options="environmentOptions"
              :label="t('projects.form.environments')"
              :placeholder="t('projects.form.selectEnvironments')"
              :can-add="canAddEnvironment"
              :add-aria-label="t('projects.form.addEnvironment')"
              @update:model-value="emit('update:environmentIds', $event)"
              @add="emit('addEnvironment')"
            />

            <CatalogMultiSelect
              :model-value="tagIds"
              :options="tagOptions"
              :label="t('projects.form.tags')"
              :placeholder="t('projects.form.selectTags')"
              :can-add="canAddTag"
              :add-aria-label="t('projects.form.addTag')"
              @update:model-value="emit('update:tagIds', $event)"
              @add="emit('addTag')"
            />
          </div>

          <div class="modal-action shrink-0">
            <button class="btn" type="button" @click="onClose">
              {{ t('projects.form.cancel') }}
            </button>
            <button class="btn btn-primary" type="submit" :disabled="saving">
              <span v-if="saving" class="loading loading-spinner loading-xs" />
              {{ mode === 'create' ? t('projects.form.create') : t('projects.form.save') }}
            </button>
          </div>
        </form>
      </template>
    </div>
    <form method="dialog" class="modal-backdrop">
      <button type="submit">{{ t('projects.form.cancel') }}</button>
    </form>
  </dialog>
</template>
