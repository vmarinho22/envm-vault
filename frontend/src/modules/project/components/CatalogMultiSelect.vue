<script setup lang="ts">
import { PlusIcon, XMarkIcon } from '@heroicons/vue/24/outline';
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

export interface CatalogOption {
  id: string;
  label: string;
}

const props = withDefaults(
  defineProps<{
    modelValue: string[];
    options: CatalogOption[];
    label: string;
    placeholder: string;
    canAdd?: boolean;
    addAriaLabel?: string;
  }>(),
  {
    canAdd: false,
    addAriaLabel: '',
  },
);

const emit = defineEmits<{
  'update:modelValue': [value: string[]];
  add: [];
}>();

const { t } = useI18n();
const rootEl = ref<HTMLElement | null>(null);
const triggerEl = ref<HTMLElement | null>(null);
const panelEl = ref<HTMLElement | null>(null);
const open = ref(false);
const query = ref('');
const teleportTo = ref<HTMLElement | string>('body');
const panelStyle = ref<Record<string, string>>({});

const selectedSet = computed(() => new Set(props.modelValue));

const filtered = computed(() => {
  const q = query.value.trim().toLowerCase();
  if (!q) return props.options;
  return props.options.filter((o) => o.label.toLowerCase().includes(q));
});

const selected = computed(() => props.options.filter((o) => selectedSet.value.has(o.id)));

function toggle(id: string): void {
  const next = new Set(props.modelValue);
  if (next.has(id)) next.delete(id);
  else next.add(id);
  emit('update:modelValue', [...next]);
}

function remove(id: string): void {
  emit(
    'update:modelValue',
    props.modelValue.filter((v) => v !== id),
  );
}

function updatePanelPosition(): void {
  const el = triggerEl.value;
  if (!el) return;
  const rect = el.getBoundingClientRect();
  const gap = 4;
  const preferredMax = 224;
  const spaceBelow = window.innerHeight - rect.bottom - gap - 8;
  const spaceAbove = rect.top - gap - 8;
  const openUp = spaceBelow < 140 && spaceAbove > spaceBelow;
  const available = openUp ? spaceAbove : spaceBelow;
  const maxHeight = Math.max(120, Math.min(preferredMax, available));

  panelStyle.value = {
    position: 'fixed',
    left: `${rect.left}px`,
    width: `${rect.width}px`,
    maxHeight: `${maxHeight}px`,
    zIndex: '100',
    ...(openUp
      ? { bottom: `${window.innerHeight - rect.top + gap}px`, top: 'auto' }
      : { top: `${rect.bottom + gap}px`, bottom: 'auto' }),
  };
}

async function setOpen(next: boolean): Promise<void> {
  if (!next) {
    open.value = false;
    query.value = '';
    return;
  }
  const dialog = triggerEl.value?.closest('dialog') ?? null;
  teleportTo.value = dialog ?? 'body';
  open.value = true;
  await nextTick();
  updatePanelPosition();
}

function onDocumentPointerDown(event: PointerEvent): void {
  const target = event.target as Node | null;
  if (!target) return;
  if (rootEl.value?.contains(target)) return;
  if (panelEl.value?.contains(target)) return;
  void setOpen(false);
}

function onViewportChange(): void {
  if (open.value) updatePanelPosition();
}

watch(open, (isOpen) => {
  if (isOpen) {
    window.addEventListener('resize', onViewportChange);
    document.addEventListener('scroll', onViewportChange, true);
  } else {
    window.removeEventListener('resize', onViewportChange);
    document.removeEventListener('scroll', onViewportChange, true);
  }
});

onMounted(() => {
  document.addEventListener('pointerdown', onDocumentPointerDown);
});

onBeforeUnmount(() => {
  document.removeEventListener('pointerdown', onDocumentPointerDown);
  window.removeEventListener('resize', onViewportChange);
  document.removeEventListener('scroll', onViewportChange, true);
});
</script>

<template>
  <div ref="rootEl" class="form-control w-full">
    <div class="flex items-center justify-between gap-2">
      <span class="label-text my-1">{{ label }}</span>
      <button
        v-if="canAdd"
        class="btn btn-ghost btn-square btn-xs"
        type="button"
        :aria-label="addAriaLabel || label"
        :title="addAriaLabel || label"
        @pointerdown.stop
        @click.prevent.stop="emit('add')"
      >
        <PlusIcon class="size-4" aria-hidden="true" />
      </button>
    </div>

    <button
      ref="triggerEl"
      class="input input-bordered flex min-h-10 w-full flex-wrap items-center gap-1 py-1 text-left focus:outline-none focus:ring-2 focus:ring-base-content/20"
      type="button"
      :aria-expanded="open"
      :aria-label="label"
      @click="setOpen(!open)"
    >
      <template v-if="selected.length === 0">
        <span class="text-base-content/50">{{ placeholder }}</span>
      </template>
      <span
        v-for="item in selected"
        :key="item.id"
        class="badge badge-sm gap-1 pr-0.5"
      >
        {{ item.label }}
        <button
          class="btn btn-ghost btn-square btn-xs size-5 min-h-0 cursor-pointer p-0"
          type="button"
          :aria-label="t('projects.search.removeOption', { label: item.label })"
          @click.stop="remove(item.id)"
        >
          <XMarkIcon class="size-3.5" aria-hidden="true" />
        </button>
      </span>
    </button>

    <Teleport :to="teleportTo" :disabled="!open">
      <div
        v-if="open"
        ref="panelEl"
        class="overflow-auto rounded-box border border-base-300 bg-base-100 p-2 shadow-lg"
        role="listbox"
        :style="panelStyle"
      >
        <input
          v-model="query"
          class="input input-bordered input-sm mb-2 w-full"
          type="search"
          :placeholder="t('projects.search.filterOptions')"
          @click.stop
        />
        <label
          v-for="opt in filtered"
          :key="opt.id"
          class="flex cursor-pointer items-center gap-2 rounded px-2 py-1.5 hover:bg-base-200"
        >
          <input
            class="checkbox checkbox-sm"
            type="checkbox"
            :checked="selectedSet.has(opt.id)"
            @change="toggle(opt.id)"
          />
          <span class="text-sm">{{ opt.label }}</span>
        </label>
        <p v-if="filtered.length === 0" class="px-2 py-1 text-sm text-base-content/60">
          {{ placeholder }}
        </p>
      </div>
    </Teleport>
  </div>
</template>
