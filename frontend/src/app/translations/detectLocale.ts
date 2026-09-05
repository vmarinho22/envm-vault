export const SUPPORTED_LOCALES = ['pt-BR', 'en-US'] as const;

export type AppLocale = (typeof SUPPORTED_LOCALES)[number];

export const DEFAULT_LOCALE: AppLocale = 'en-US';

function normalizeTag(tag: string): string {
  return tag.trim().replace(/_/g, '-');
}

function matchLocale(tag: string): AppLocale | null {
  const normalized = normalizeTag(tag);
  if (!normalized) return null;

  const lower = normalized.toLowerCase();
  for (const locale of SUPPORTED_LOCALES) {
    if (lower === locale.toLowerCase()) {
      return locale;
    }
  }

  const language = lower.split('-')[0];
  if (language === 'pt') return 'pt-BR';
  if (language === 'en') return 'en-US';
  return null;
}

/** Pick the first supported locale from a browser language preference list. */
export function detectLocale(languages: readonly string[]): AppLocale {
  for (const tag of languages) {
    const matched = matchLocale(tag);
    if (matched) return matched;
  }
  return DEFAULT_LOCALE;
}

export function detectBrowserLocale(): AppLocale {
  const languages =
    typeof navigator !== 'undefined'
      ? navigator.languages?.length
        ? [...navigator.languages]
        : navigator.language
          ? [navigator.language]
          : []
      : [];
  return detectLocale(languages);
}
