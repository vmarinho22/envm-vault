<script setup lang="ts">
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';

export interface UserOption {
  id: string;
  email: string;
}

const props = defineProps<{
  modelValue: string[];
  options: UserOption[];
  label: string;
}>();

const emit = defineEmits<{
  'update:modelValue': [value: string[]];
}>();

const { t } = useI18n();
const open = ref(false);
const query = ref('');

const selectedSet = computed(() => new Set(props.modelValue));

const filtered = computed(() => {
  const q = query.value.trim().toLowerCase();
  if (!q) return props.options;
  return props.options.filter((u) => u.email.toLowerCase().includes(q));
});

const selectedUsers = computed(() =>
  props.options.filter((u) => selectedSet.value.has(u.id)),
);

function toggle(id: string): void {
  const next = new Set(props.modelValue);
  if (next.has(id)) {
    next.delete(id);
  } else {
    next.add(id);
  }
  emit('update:modelValue', [...next]);
}

function remove(id: string): void {
  emit(
    'update:modelValue',
    props.modelValue.filter((v) => v !== id),
  );
}

function clearAll(): void {
  emit('update:modelValue', []);
}
</script>

<template>
  <div class="form-control flex w-full max-w-xs flex-col gap-[4px]">
    <span class="label-text">{{ label }}</span>
    <div class="dropdown w-full" :class="{ 'dropdown-open': open }">
      <div
        tabindex="0"
        role="button"
        class="input input-bordered input-sm flex min-h-8 w-full flex-wrap items-center gap-1 py-1"
        @click="open = !open"
      >
        <template v-if="selectedUsers.length === 0">
          <span class="text-base-content/50">{{ t('audit.filters.all') }}</span>
        </template>
        <span
          v-for="u in selectedUsers"
          :key="u.id"
          class="badge badge-sm gap-1"
          @click.stop="remove(u.id)"
        >
          {{ u.email }}
          <span aria-hidden="true">×</span>
        </span>
      </div>
      <div
        tabindex="0"
        class="dropdown-content z-20 mt-1 w-full rounded-box border border-base-300 bg-base-100 p-2 shadow"
      >
        <input
          v-model="query"
          class="input input-bordered input-sm mb-2 w-full"
          type="search"
          :placeholder="t('audit.filters.searchUsers')"
          @click.stop
        />
        <ul class="max-h-48 overflow-y-auto">
          <li v-if="filtered.length === 0" class="px-2 py-1 text-sm text-base-content/60">
            {{ t('audit.filters.noUsers') }}
          </li>
          <li v-for="u in filtered" :key="u.id">
            <label class="flex cursor-pointer items-center gap-2 rounded-md px-2 py-1.5 hover:bg-base-200">
              <input
                class="checkbox checkbox-sm"
                type="checkbox"
                :checked="selectedSet.has(u.id)"
                @change="toggle(u.id)"
              />
              <span class="truncate text-sm">{{ u.email }}</span>
            </label>
          </li>
        </ul>
        <button
          v-if="modelValue.length > 0"
          class="btn btn-ghost btn-xs mt-2"
          type="button"
          @click="clearAll"
        >
          {{ t('audit.filters.all') }}
        </button>
      </div>
    </div>
  </div>
</template>
