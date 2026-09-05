<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const CODE_LENGTH = 8;

const model = defineModel<string>({ default: '' });

const props = defineProps<{
  disabled?: boolean;
}>();

const { t } = useI18n();

const digits = ref<string[]>(Array.from({ length: CODE_LENGTH }, () => ''));
const inputRefs = ref<(HTMLInputElement | null)[]>([]);

const codeValue = computed(() => digits.value.join(''));

watch(
  model,
  (value) => {
    const normalized = value.replace(/\D/g, '').slice(0, CODE_LENGTH);
    const next = Array.from({ length: CODE_LENGTH }, (_, index) => normalized[index] ?? '');
    if (next.join('') === digits.value.join('')) {
      return;
    }
    digits.value = next;
  },
  { immediate: true },
);

watch(codeValue, (value) => {
  if (value !== model.value) {
    model.value = value;
  }
});

function setInputRef(element: unknown, index: number): void {
  inputRefs.value[index] = element instanceof HTMLInputElement ? element : null;
}

function focusAt(index: number): void {
  void nextTick(() => {
    inputRefs.value[index]?.focus();
    inputRefs.value[index]?.select();
  });
}

function applyDigits(raw: string, startIndex = 0): void {
  const cleaned = raw.replace(/\D/g, '');
  if (!cleaned) {
    return;
  }

  const next = [...digits.value];
  let cursor = startIndex;
  for (const char of cleaned) {
    if (cursor >= CODE_LENGTH) {
      break;
    }
    next[cursor] = char;
    cursor += 1;
  }
  digits.value = next;
  focusAt(Math.min(cursor, CODE_LENGTH - 1));
}

function onInput(index: number, event: Event): void {
  const target = event.target as HTMLInputElement;
  const cleaned = target.value.replace(/\D/g, '');

  if (!cleaned) {
    digits.value[index] = '';
    return;
  }

  if (cleaned.length > 1) {
    applyDigits(cleaned, index);
    return;
  }

  digits.value[index] = cleaned;
  if (index < CODE_LENGTH - 1) {
    focusAt(index + 1);
  }
}

function onKeydown(index: number, event: KeyboardEvent): void {
  if (event.key === 'Backspace') {
    if (digits.value[index]) {
      digits.value[index] = '';
      return;
    }
    if (index > 0) {
      event.preventDefault();
      digits.value[index - 1] = '';
      focusAt(index - 1);
    }
    return;
  }

  if (event.key === 'ArrowLeft' && index > 0) {
    event.preventDefault();
    focusAt(index - 1);
    return;
  }

  if (event.key === 'ArrowRight' && index < CODE_LENGTH - 1) {
    event.preventDefault();
    focusAt(index + 1);
  }
}

function onPaste(index: number, event: ClipboardEvent): void {
  event.preventDefault();
  const pasted = event.clipboardData?.getData('text') ?? '';
  applyDigits(pasted, index);
}
</script>

<template>
  <div
    class="flex flex-wrap gap-2"
    role="group"
    :aria-label="t('code.groupAriaLabel')"
  >
    <input
      v-for="(_, index) in digits"
      :key="index"
      :ref="(el) => setInputRef(el, index)"
      class="input input-bordered w-10 px-0 text-center text-base tracking-normal"
      type="text"
      inputmode="numeric"
      maxlength="1"
      pattern="[0-9]"
      :name="index === 0 ? 'code' : undefined"
      :autocomplete="index === 0 ? 'one-time-code' : 'off'"
      :aria-label="t('code.digitAriaLabel', { n: index + 1 })"
      :value="digits[index]"
      :disabled="props.disabled"
      @input="onInput(index, $event)"
      @keydown="onKeydown(index, $event)"
      @paste="onPaste(index, $event)"
      @focus="($event.target as HTMLInputElement).select()"
    />
  </div>
</template>
