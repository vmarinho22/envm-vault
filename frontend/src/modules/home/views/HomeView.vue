<script setup lang="ts">
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
</script>

<template>
  <section class="px-8 py-10 md:px-12">
    <h1 class="text-2xl font-medium tracking-tight md:text-3xl">{{ t('home.title') }}</h1>
    <p class="mt-3 max-w-lg text-base text-base-content/70">{{ t('home.welcome') }}</p>
  </section>
</template>
