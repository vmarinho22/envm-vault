import { createI18n } from 'vue-i18n';
import {
  detectBrowserLocale,
  SUPPORTED_LOCALES,
  type AppLocale,
} from './detectLocale';
import { allLocaleMessages, type MessageSchema } from './loadLocaleMessages';

export type { AppLocale, MessageSchema };
export { detectLocale, detectBrowserLocale, DEFAULT_LOCALE, SUPPORTED_LOCALES } from './detectLocale';
export { allLocaleMessages, loadLocaleMessages } from './loadLocaleMessages';

const LOCALE_STORAGE_KEY = 'envm-vault.locale';

export function isAppLocale(value: string): value is AppLocale {
  return (SUPPORTED_LOCALES as readonly string[]).includes(value);
}

export function readStoredLocale(): AppLocale | null {
  if (typeof localStorage === 'undefined') return null;
  const stored = localStorage.getItem(LOCALE_STORAGE_KEY);
  if (stored && isAppLocale(stored)) return stored;
  return null;
}

export function resolveInitialLocale(): AppLocale {
  return readStoredLocale() ?? detectBrowserLocale();
}

const locale = resolveInitialLocale();

export const i18n = createI18n({
  legacy: false,
  locale,
  fallbackLocale: 'en-US',
  messages: allLocaleMessages(),
});

export function applyDocumentLang(nextLocale: AppLocale): void {
  if (typeof document !== 'undefined') {
    document.documentElement.lang = nextLocale;
  }
}

/** Persist preference and sync <html lang>. Does not mutate a Composer — callers set locale via useI18n(). */
export function persistLocale(nextLocale: AppLocale): void {
  applyDocumentLang(nextLocale);
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem(LOCALE_STORAGE_KEY, nextLocale);
  }
}

export function setAppLocale(nextLocale: AppLocale): void {
  i18n.global.locale.value = nextLocale;
  persistLocale(nextLocale);
}
