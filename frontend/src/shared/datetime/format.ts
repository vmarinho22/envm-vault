import dayjs from 'dayjs';
import localizedFormat from 'dayjs/plugin/localizedFormat';
import utc from 'dayjs/plugin/utc';
import 'dayjs/locale/en';
import 'dayjs/locale/pt-br';
import { i18n, isAppLocale, type AppLocale } from '@/app/translations';

dayjs.extend(utc);
dayjs.extend(localizedFormat);

const DAYJS_LOCALE: Record<AppLocale, string> = {
  'pt-BR': 'pt-br',
  'en-US': 'en',
};

/** Display formats per app locale (24h clock). */
const DATETIME_FORMAT: Record<AppLocale, string> = {
  'pt-BR': 'DD/MM/YYYY HH:mm',
  'en-US': 'MM/DD/YYYY HH:mm',
};

function resolveAppLocale(locale?: string): AppLocale {
  if (locale && isAppLocale(locale)) {
    return locale;
  }
  const current = i18n.global.locale.value;
  return isAppLocale(current) ? current : 'en-US';
}

/** Format an ISO datetime for display in the active (or given) app locale. */
export function formatDateTime(iso: string, locale?: string): string {
  const appLocale = resolveAppLocale(locale);
  const parsed = dayjs(iso).locale(DAYJS_LOCALE[appLocale]);
  if (!parsed.isValid()) {
    return iso;
  }
  return parsed.format(DATETIME_FORMAT[appLocale]);
}

/**
 * Format a value as YYYY-MM-DD for API date filters.
 * Always uses local calendar day at 00:00 (time is ignored).
 */
export function toDateParam(value: string | Date | dayjs.Dayjs | null | undefined): string {
  if (value == null || value === '') {
    return '';
  }
  const parsed = dayjs(value).startOf('day');
  if (!parsed.isValid()) {
    return '';
  }
  return parsed.format('YYYY-MM-DD');
}
