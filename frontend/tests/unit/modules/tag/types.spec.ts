import { describe, expect, it } from 'vitest';
import { randomTagColor } from '@/modules/tag/types';

describe('randomTagColor', () => {
  it('returns a valid uppercase #RRGGBB hex', () => {
    for (let i = 0; i < 20; i += 1) {
      expect(randomTagColor()).toMatch(/^#[0-9A-F]{6}$/);
    }
  });
});
