import { describe, expect, it } from 'vitest';
import { detectLocale } from '@/app/translations/detectLocale';

describe('detectLocale', () => {
  it('matches exact supported tags', () => {
    expect(detectLocale(['pt-BR'])).toBe('pt-BR');
    expect(detectLocale(['en-US'])).toBe('en-US');
  });

  it('matches language prefixes', () => {
    expect(detectLocale(['pt-PT'])).toBe('pt-BR');
    expect(detectLocale(['en-GB'])).toBe('en-US');
    expect(detectLocale(['pt'])).toBe('pt-BR');
    expect(detectLocale(['en'])).toBe('en-US');
  });

  it('normalizes underscores', () => {
    expect(detectLocale(['pt_BR'])).toBe('pt-BR');
    expect(detectLocale(['en_US'])).toBe('en-US');
  });

  it('uses the first supported tag in the list', () => {
    expect(detectLocale(['fr-FR', 'en-GB', 'pt-BR'])).toBe('en-US');
  });

  it('falls back to en-US when nothing matches', () => {
    expect(detectLocale([])).toBe('en-US');
    expect(detectLocale(['fr-FR', 'de-DE'])).toBe('en-US');
  });
});
