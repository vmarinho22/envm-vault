<script setup lang="ts">
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import {
  isAppLocale,
  persistLocale,
  SUPPORTED_LOCALES,
  type AppLocale,
} from '@/app/translations';

const { locale, t } = useI18n();

const LOCALE_LABELS: Record<AppLocale, string> = {
  'pt-BR': 'Português',
  'en-US': 'English',
};

const selected = computed({
  get: () => locale.value as AppLocale,
  set: (value: string) => {
    if (!isAppLocale(value)) return;
    locale.value = value;
    persistLocale(value);
  },
});
</script>

<template>
  <label class="sr-only" for="locale-switcher">{{ t('common.language') }}</label>
  <select
    id="locale-switcher"
    v-model="selected"
    class="select select-bordered select-sm w-auto min-w-28"
    data-testid="locale-switcher"
    :aria-label="t('common.language')"
  >
    <option v-for="code in SUPPORTED_LOCALES" :key="code" :value="code">
      {{ LOCALE_LABELS[code] }}
    </option>
  </select>
</template>
