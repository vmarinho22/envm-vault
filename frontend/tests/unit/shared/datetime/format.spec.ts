import { afterEach, describe, expect, it } from 'vitest';
import { i18n } from '@/app/translations';
import { formatDateTime, toDateParam } from '@/shared/datetime/format';

describe('formatDateTime', () => {
  afterEach(() => {
    i18n.global.locale.value = 'en-US';
  });

  it('formats ISO timestamps in pt-BR style', () => {
    expect(formatDateTime('2026-09-06T15:43:00.000Z', 'pt-BR')).toMatch(
      /^\d{2}\/\d{2}\/2026 \d{2}:\d{2}$/,
    );
    expect(formatDateTime('2026-09-06T15:43:00.000Z', 'pt-BR')).toMatch(/^06\/09\/2026/);
  });

  it('formats ISO timestamps in en-US style', () => {
    expect(formatDateTime('2026-09-06T15:43:00.000Z', 'en-US')).toMatch(
      /^\d{2}\/\d{2}\/2026 \d{2}:\d{2}$/,
    );
    expect(formatDateTime('2026-09-06T15:43:00.000Z', 'en-US')).toMatch(/^09\/06\/2026/);
  });

  it('uses the active i18n locale when none is passed', () => {
    i18n.global.locale.value = 'pt-BR';
    expect(formatDateTime('2026-09-06T15:43:00.000Z')).toMatch(/^06\/09\/2026/);

    i18n.global.locale.value = 'en-US';
    expect(formatDateTime('2026-09-06T15:43:00.000Z')).toMatch(/^09\/06\/2026/);
  });

  it('returns original string when invalid', () => {
    expect(formatDateTime('not-a-date')).toBe('not-a-date');
  });
});

describe('toDateParam', () => {
  it('formats as YYYY-MM-DD', () => {
    expect(toDateParam('2026-09-06T12:00:00.000Z')).toBe('2026-09-06');
  });

  it('returns empty for nullish', () => {
    expect(toDateParam(null)).toBe('');
    expect(toDateParam(undefined)).toBe('');
  });
});
