import { describe, expect, it } from 'vitest';
import { fromName } from '@/shared/slug/fromName';

describe('fromName', () => {
  it.each([
    ['Platform', 'platform'],
    ['Equipe São Paulo', 'equipe-sao-paulo'],
    ['Foo & Bar!!!', 'foo-bar'],
    ['  My -- Team  ', 'my-team'],
    ['Team 42', 'team-42'],
    ['   ', ''],
    ['!!!', ''],
    ['Ação', 'acao'],
    ['Über Team', 'uber-team'],
  ])('fromName(%j) → %j', (input, expected) => {
    expect(fromName(input)).toBe(expected);
  });
});
